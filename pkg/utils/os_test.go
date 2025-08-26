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
