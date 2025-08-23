package vscode

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func strptr(s string) *string { return &s }

type setupResult struct {
	settingsFile string
	cleanup      func()
}

func createFile(t *testing.T, dir, name, content string, perm os.FileMode) setupResult {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	called := false
	cleanup := func() {
		if called {
			return
		}
		_ = os.Chmod(path, 0o600) // try to ensure it's deletable
		_ = os.Remove(path)
		called = true
	}
	return setupResult{settingsFile: path, cleanup: cleanup}
}

func Test_UpdateSettings_Cases_Factorized(t *testing.T) {
	tests := []struct {
		name              string
		prepare           func(t *testing.T) setupResult
		newSettings       map[string]any
		wantErrContains   string
		skipOnWindowsPerm bool
	}{
		{
			name: "non-existent file",
			prepare: func(t *testing.T) setupResult {
				dir := t.TempDir()
				file := filepath.Join(dir, "no-such.json")
				// do NOT create file
				return setupResult{settingsFile: file, cleanup: func() {}}
			},
			newSettings:     map[string]any{"testDevBoxKey": "value"},
			wantErrContains: "failed to read settings file",
		},
		{
			name: "path is directory",
			prepare: func(t *testing.T) setupResult {
				dir := t.TempDir()
				// use the directory itself as the "file"
				return setupResult{settingsFile: dir, cleanup: func() {}}
			},
			newSettings:     map[string]any{"testDevBoxKey": "value"},
			wantErrContains: "failed to read settings file",
		},
		{
			name: "no read permission",
			prepare: func(t *testing.T) setupResult {
				dir := t.TempDir()
				r := createFile(t, dir, "settings.json", `{"existing":"v"}`, 0o600)
				// remove all permissions
				if err := os.Chmod(r.settingsFile, 0o000); err != nil {
					t.Fatalf("chmod failed: %v", err)
				}
				cleanup := func() {
					_ = os.Chmod(r.settingsFile, 0o600)
					_ = os.Remove(r.settingsFile)
				}
				return setupResult{settingsFile: r.settingsFile, cleanup: cleanup}
			},
			newSettings:       map[string]any{"testDevBoxKey": "value"},
			wantErrContains:   "failed to read settings file",
			skipOnWindowsPerm: true,
		},
		{
			name: "no write permission",
			prepare: func(t *testing.T) setupResult {
				dir := t.TempDir()
				r := createFile(t, dir, "settings.json", `{"existing":"v"}`, 0o400) // read-only
				cleanup := func() {
					_ = os.Chmod(r.settingsFile, 0o600)
					_ = os.Remove(r.settingsFile)
				}
				return setupResult{settingsFile: r.settingsFile, cleanup: cleanup}
			},
			newSettings:       map[string]any{"testDevBoxKey": "value"},
			wantErrContains:   "failed to write updated settings to file",
			skipOnWindowsPerm: true,
		},
		{
			name: "success merges testDevBox keys",
			prepare: func(t *testing.T) setupResult {
				dir := t.TempDir()
				existing := map[string]any{
					"existingKey":   "old",
					"testDevBoxOld": "oldVal",
				}
				b, _ := json.MarshalIndent(existing, "", "  ")
				r := createFile(t, dir, "settings.json", string(b), 0o600)
				return r
			},
			newSettings:     map[string]any{"testDevBoxNew": "newVal", "testDevBoxOld": "updatedVal"},
			wantErrContains: "",
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			// Skip permission-altering tests on Windows
			if tt.skipOnWindowsPerm && runtime.GOOS == "windows" {
				t.Skip("skipping permission test on windows")
			}
			// allow parallelism except for permission tests (they modify FS perms of temp files but isolated)
			if !tt.skipOnWindowsPerm {
				t.Parallel()
			}

			setup := tt.prepare(t)
			defer setup.cleanup()

			code := &VSCode{SettingsFile: strptr(setup.settingsFile)}
			err := code.UpdateSettings(tt.newSettings)

			if tt.wantErrContains != "" {
				if err == nil {
					t.Fatalf("expected error containing %q but got nil", tt.wantErrContains)
				}
				if !strings.Contains(err.Error(), tt.wantErrContains) {
					t.Fatalf("expected error containing %q but got: %v", tt.wantErrContains, err)
				}
				return
			}

			// success path: assert file contains merged settings
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			gotB, err := os.ReadFile(setup.settingsFile)
			if err != nil {
				t.Fatalf("failed to read resulting settings: %v", err)
			}
			var merged map[string]any
			if err := json.Unmarshal(gotB, &merged); err != nil {
				t.Fatalf("failed to parse resulting settings: %v", err)
			}

			// existing key should remain
			if merged["existingKey"] != "old" {
				t.Fatalf("expected existingKey to remain 'old', got: %v", merged["existingKey"])
			}
			// new keys from tt.newSettings must be present and set
			for k, v := range tt.newSettings {
				if !strings.HasPrefix(k, "testDevBox") {
					t.Fatalf("test provided a key that does not start with testDevBox: %s", k)
				}
				if merged[k] != v {
					t.Fatalf("expected %s to be %v, got: %v", k, v, merged[k])
				}
			}
		})
	}

}
