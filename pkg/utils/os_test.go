package utils

import (
	"os"
	"testing"
)

func TestGetenv_DefaultBehavior(t *testing.T) {
	key := "DETECT_ENV_TEST"
	// Ensure it's unset
	_ = os.Unsetenv(key)
	val := Getenv(key, "default-val")
	if val != "default-val" {
		t.Fatalf("expected default value, got %q", val)
	}

	// Set and verify
	if err := os.Setenv(key, "set-val"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer func() { _ = os.Unsetenv(key) }()

	val = Getenv(key, "default-val")
	if val != "set-val" {
		t.Fatalf("expected set value, got %q", val)
	}
}
