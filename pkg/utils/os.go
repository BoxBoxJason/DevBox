package utils

import (
	"os"
)

// Getenv retrieves the value of the environment variable named by key.
// If the variable is not set, it returns the provided default value.
func Getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
