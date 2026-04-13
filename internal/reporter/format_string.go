package reporter

import "fmt"

// String returns the string representation of a Format.
func (f Format) String() string {
	return string(f)
}

// ParseFormat converts a raw string into a Format value.
// It returns an error if the string does not match a known format.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatText, FormatJSON:
		return Format(s), nil
	default:
		return "", fmt.Errorf("unknown report format %q: must be \"text\" or \"json\"", s)
	}
}

// Formats returns all supported Format values.
func Formats() []Format {
	return []Format{FormatText, FormatJSON}
}
