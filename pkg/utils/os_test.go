package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func withModifiedPATH(t *testing.T, newPath string, fn func()) {
	t.Helper()
	orig := os.Getenv("PATH")
	if err := os.Setenv("PATH", newPath); err != nil {
		t.Fatalf("failed to set PATH: %v", err)
	}
	defer func() {
		_ = os.Setenv("PATH", orig)
	}()
	fn()
}

func TestDetectPackageManager_SuccessAndFailure(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip PATH/executable tests on Windows")
	}

	dir := t.TempDir()
	aptPath := filepath.Join(dir, "apt")
	writeExecutable(t, aptPath, "#!/bin/sh\nexit 0\n")

	withModifiedPATH(t, dir+string(os.PathListSeparator)+os.Getenv("PATH"), func() {
		pm, err := DetectPackageManager()
		if err != nil {
			t.Fatalf("expected DetectPackageManager to succeed, got error: %v", err)
		}
		if pm == nil || pm.Name != "apt" {
			t.Fatalf("expected package manager apt, got: %+v", pm)
		}
	})

	// Failure: PATH with no known package managers
	emptyDir := t.TempDir()
	withModifiedPATH(t, emptyDir, func() {
		pm, err := DetectPackageManager()
		if err == nil {
			t.Fatalf("expected error when no package manager present, got pm: %+v", pm)
		}
	})
}

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

func TestLoadEnv_SuccessAndError(t *testing.T) {
	// Success case: expansion and setting
	_ = os.Unsetenv("LOAD_BASE")
	_ = os.Unsetenv("LOAD_EXPANDED")

	vars := map[string]string{
		"LOAD_BASE":     "base",
		"LOAD_EXPANDED": "${LOAD_BASE}_x",
	}
	errs := LoadEnv(vars)
	if len(errs) != 0 {
		t.Fatalf("expected no errors from LoadEnv, got: %+v", errs)
	}
	if got := os.Getenv("LOAD_BASE"); got != "base" {
		t.Fatalf("expected LOAD_BASE=base, got %q", got)
	}
	if got := os.Getenv("LOAD_EXPANDED"); got != "base_x" {
		t.Fatalf("expected LOAD_EXPANDED=base_x, got %q", got)
	}

	// Error case: invalid key should produce at least one error
	badVars := map[string]string{
		"": "some", // empty key should cause os.Setenv to return an error on most systems
	}
	errs = LoadEnv(badVars)
	foundErr := false
	for _, e := range errs {
		if e != nil {
			foundErr = true
			break
		}
	}
	if !foundErr {
		t.Fatalf("expected at least one error from LoadEnv with invalid key, got nil slice: %+v", errs)
	}
}
