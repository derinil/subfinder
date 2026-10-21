package stringsutil

import (
	"strings"
	"unicode"
)

type NormalizeOptions struct {
	TrimSpaces    bool
	TrimCutset    string
	Lowercase     bool
	Uppercase     bool
	StripComments bool
}

var DefaultNormalizeOptions NormalizeOptions = NormalizeOptions{
	TrimSpaces: true,
}

func NormalizeWithOptions(data string, options NormalizeOptions) string {
	if options.TrimSpaces {
		data = strings.TrimSpace(data)
	}

	if options.TrimCutset != "" {
		data = strings.Trim(data, options.TrimCutset)
	}

	if options.Lowercase {
		data = strings.ToLower(data)
	}

	if options.Uppercase {
		data = strings.ToUpper(data)
	}

	if options.StripComments {
		if cut := strings.IndexAny(data, "#"); cut >= 0 {
			data = strings.TrimRightFunc(data[:cut], unicode.IsSpace)
		}
	}

	return data
}

func Normalize(data string) string {
	return NormalizeWithOptions(data, DefaultNormalizeOptions)
}
