package utils

import (
	"strings"
)

func StrPtr(s string) *string {
	return &s
}

func TrimSpacesAndQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 1 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1 : len(s)-1]
	}
	return s
}
