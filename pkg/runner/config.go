package runner

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/derinil/subfinder/v2/pkg/passive"
	"github.com/projectdiscovery/gologger"
)

// createProviderConfigYAML marshals the input map to the given location on the disk
func createProviderConfigYAML(configFilePath string) error {
	configFile, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err := configFile.Close(); err != nil {
			gologger.Error().Msgf("Error closing config file: %s", err)
		}
	}()

	sourcesRequiringApiKeysMap := make(map[string][]string)
	for _, source := range passive.AllSources {
		if source.NeedsKey() {
			sourceName := strings.ToLower(source.Name())
			sourcesRequiringApiKeysMap[sourceName] = []string{}
		}
	}

	return yaml.NewEncoder(configFile).Encode(sourcesRequiringApiKeysMap)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func ReadFile(filename string) (chan string, error) {
	if !FileExists(filename) {
		return nil, errors.New("file doesn't exist")
	}
	out := make(chan string)
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		if err != nil {
			return
		}
		defer func() {
			_ = f.Close()
		}()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			out <- scanner.Text()
		}
	}()

	return out, nil
}

func substituteEnvVars(line string) string {
	for _, word := range strings.Fields(line) {
		word = strings.Trim(word, `"`)
		if key, ok := strings.CutPrefix(word, "$"); ok {
			substituteEnv := os.Getenv(key)
			if substituteEnv != "" {
				line = strings.Replace(line, word, substituteEnv, 1)
			}
		}
	}
	return line
}

func SubstituteConfigFromEnvVars(filepath string) (io.Reader, error) {
	var config strings.Builder

	lines, err := ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	for line := range lines {
		config.WriteString(substituteEnvVars(string(line)))
		config.WriteString("\n")
	}

	return strings.NewReader(config.String()), nil
}

// UnmarshalFrom writes the marshaled yaml config to disk
func UnmarshalFrom(file string) error {
	reader, err := SubstituteConfigFromEnvVars(file)
	if err != nil {
		return err
	}

	sourceApiKeysMap := map[string][]string{}
	err = yaml.NewDecoder(reader).Decode(sourceApiKeysMap)
	for _, source := range passive.AllSources {
		sourceName := strings.ToLower(source.Name())
		apiKeys := sourceApiKeysMap[sourceName]
		if source.NeedsKey() && apiKeys != nil && len(apiKeys) > 0 {
			gologger.Debug().Msgf("API key(s) found for %s.", sourceName)
			source.AddApiKeys(apiKeys)
		}
	}
	return err
}
