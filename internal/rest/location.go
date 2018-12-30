package rest

import "strings"

// Location returns a location joining strings
func Location(parts ...string) string {
	return strings.Join(parts, "/")
}
