package graph

import "fmt"

type FormatOutput string

var allowedFormats = map[string]bool{ //nolint:gochecknoglobals
	"png": true,
	"svg": true,
}

func (f *FormatOutput) Set(value string) error {
	if !allowedFormats[value] {
		return fmt.Errorf("%w: with value: %s", ErrFormatNotAllowed, value)
	}

	*f = FormatOutput(value)

	return nil
}

func (f *FormatOutput) String() string {
	return string(*f)
}
