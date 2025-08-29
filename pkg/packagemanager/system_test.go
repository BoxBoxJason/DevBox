package packagemanager

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

// writeExecutable creates a file with content and makes it executable.
func writeExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o700); err != nil {
		t.Fatalf("failed to write script %s: %v", path, err)
	}
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
