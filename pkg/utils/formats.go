package utils

import (
	"strings"
)

// StrPtr returns a pointer to the given string.
func StrPtr(s string) *string {
	return &s
}

// TrimSpacesAndQuotes trims leading and trailing spaces and quotes from a string.
func TrimSpacesAndQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 1 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}

// MergeStringSlices merges multiple string slices into one, removing duplicates while preserving order.
func MergeStringSlices(slices ...[]string) []string {
	merged := []string{}
	seen := make(map[string]struct{})
	for _, slice := range slices {
		for _, str := range slice {
			if _, exists := seen[str]; !exists {
				seen[str] = struct{}{}
				merged = append(merged, str)
			}
		}
	}
	return merged
}
