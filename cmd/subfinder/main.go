package main

import (
	"github.com/derinil/subfinder/v2/pkg/runner"
	// Attempts to increase the OS file descriptors - Fail silently
	"github.com/derinil/subfinder/v2/pkg/gologger"
	_ "github.com/projectdiscovery/fdmax/autofdmax"
)

func main() {
	// Parse the command line flags and read config files
	options := runner.ParseOptions()

	newRunner, err := runner.NewRunner(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}

	err = newRunner.RunEnumeration()
	if err != nil {
		gologger.Fatal().Msgf("Could not run enumeration: %s\n", err)
	}
}
