package runner

import "github.com/projectdiscovery/subfinder/v2/pkg/stringsutil"

func loadFromFile(file string) ([]string, error) {
	chanItems, err := ReadFile(file)
	if err != nil {
		return nil, err
	}
	var items []string
	for item := range chanItems {
		item = preprocessDomain(item)
		if item == "" {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func preprocessDomain(s string) string {
	return stringsutil.NormalizeWithOptions(s,
		stringsutil.NormalizeOptions{
			StripComments: true,
			TrimCutset:    "\n\t\"'` ",
			Lowercase:     true,
		},
	)
}
