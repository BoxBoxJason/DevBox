package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"go.uber.org/zap"
)

type setupResult struct {
	dir     string
	cleanup func()
	marker  string
}

// writeExecutable creates a file with content and makes it executable.
func writeExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("failed to write script %s: %v", path, err)
	}
}

// makeDistroboxScript writes a distrobox-export script with the provided body.
func makeDistroboxScript(t *testing.T, dir, body string) string {
	t.Helper()
	script := filepath.Join(dir, "distrobox-export")
	writeExecutable(t, script, body)
	return script
}

// makeFakeBinary writes a fake binary with the given name (executable script).
func makeFakeBinary(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	writeExecutable(t, path, body)
	return path
}

func TestMain(m *testing.M) {
	// avoid zap global logger side effects
	zap.ReplaceGlobals(zap.NewNop())
	os.Exit(m.Run())
}

func Test_IsDistroboxBinaryExported_and_ExportBinary_and_Apps(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping distrobox script tests on Windows")
	}

	origPath := os.Getenv("PATH")

	tests := []struct {
		name           string
		prepare        func(t *testing.T, dir string) setupResult
		testFunc       func(t *testing.T, s setupResult)
		resetAvailable bool // whether to recompute distroboxExportAvailable after prepare
	}{
		{
			name: "IsDistroboxBinaryExported - available and present",
			prepare: func(t *testing.T, dir string) setupResult {
				// distrobox-export lists "mybin"
				body := `#!/bin/sh
if [ "$1" = "--list-binaries" ]; then
	echo "mybin"
	exit 0
fi
exit 0
`
				makeDistroboxScript(t, dir, body)
				return setupResult{dir: dir, cleanup: func() {}}
			},
			resetAvailable: true,
			testFunc: func(t *testing.T, s setupResult) {
				distroboxExportAvailable = IsDistroboxExportAvailable()
				ok, err := IsDistroboxBinaryExported("mybin")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !ok {
					t.Fatalf("expected mybin to be reported as exported")
				}
			},
		},
		{
			name: "IsDistroboxBinaryExported - available and not present",
			prepare: func(t *testing.T, dir string) setupResult {
				body := `#!/bin/sh
if [ "$1" = "--list-binaries" ]; then
	echo "other"
	exit 0
fi
exit 0
`
				makeDistroboxScript(t, dir, body)
				return setupResult{dir: dir, cleanup: func() {}}
			},
			resetAvailable: true,
			testFunc: func(t *testing.T, s setupResult) {
				distroboxExportAvailable = IsDistroboxExportAvailable()
				ok, err := IsDistroboxBinaryExported("mybin")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if ok {
					t.Fatalf("expected mybin to NOT be reported as exported")
				}
			},
		},
		{
			name: "IsDistroboxBinaryExported - not available",
			prepare: func(t *testing.T, dir string) setupResult {
				// do not create distrobox-export
				return setupResult{dir: dir, cleanup: func() {}}
			},
			resetAvailable: false,
			testFunc: func(t *testing.T, s setupResult) {
				// force not available
				distroboxExportAvailable = false
				_, err := IsDistroboxBinaryExported("mybin")
				if err == nil || !strings.Contains(err.Error(), "distrobox-export command is not available") {
					t.Fatalf("expected not available error, got: %v", err)
				}
			},
		},
		{
			name: "ExportDistroboxBinary - success calls exporter",
			prepare: func(t *testing.T, dir string) setupResult {
				// exporter returns empty list so binary is not exported,
				// and when --bin is invoked create a marker file.
				marker := filepath.Join(dir, "exported.marker")
				body := `#!/bin/sh
if [ "$1" = "--list-binaries" ]; then
	# none exported
	exit 0
fi
if [ "$1" = "--bin" ]; then
	# create marker to indicate export was invoked
	echo "exporting $2" > "` + marker + `"
	exit 0
fi
exit 0
`
				makeDistroboxScript(t, dir, body)
				// create a fake binary that will be found by LookPath
				makeFakeBinary(t, dir, "mybin", "#!/bin/sh\nexit 0\n")
				return setupResult{dir: dir, marker: marker, cleanup: func() {}}
			},
			resetAvailable: true,
			testFunc: func(t *testing.T, s setupResult) {
				distroboxExportAvailable = IsDistroboxExportAvailable()
				err := ExportDistroboxBinary("mybin")
				if err != nil {
					t.Fatalf("expected export to succeed, got: %v", err)
				}
				// marker should exist
				if _, err := os.Stat(s.marker); err != nil {
					t.Fatalf("expected marker file to be created: %v", err)
				}
			},
		},
		{
			name: "ExportDistroboxBinary - already exported does not call exporter",
			prepare: func(t *testing.T, dir string) setupResult {
				// exporter lists mybin so Export should no-op
				marker := filepath.Join(dir, "exported.marker")
				body := `#!/bin/sh
if [ "$1" = "--list-binaries" ]; then
	echo "mybin"
	exit 0
fi
if [ "$1" = "--bin" ]; then
	# would create marker if invoked
	echo "exporting $2" > "` + marker + `"
	exit 0
fi
exit 0
`
				makeDistroboxScript(t, dir, body)
				makeFakeBinary(t, dir, "mybin", "#!/bin/sh\nexit 0\n")
				return setupResult{dir: dir, marker: marker, cleanup: func() {}}
			},
			resetAvailable: true,
			testFunc: func(t *testing.T, s setupResult) {
				distroboxExportAvailable = IsDistroboxExportAvailable()
				err := ExportDistroboxBinary("mybin")
				if err != nil {
					t.Fatalf("expected export to no-op successfully, got: %v", err)
				}
				if _, err := os.Stat(s.marker); !os.IsNotExist(err) {
					t.Fatalf("expected marker NOT to be created when already exported")
				}
			},
		},
		{
			name: "ExportDistroboxBinary - exporter fails",
			prepare: func(t *testing.T, dir string) setupResult {
				body := `#!/bin/sh
if [ "$1" = "--list-binaries" ]; then
	exit 0
fi
if [ "$1" = "--bin" ]; then
	echo "boom" >&2
	exit 2
fi
exit 0
`
				makeDistroboxScript(t, dir, body)
				makeFakeBinary(t, dir, "failbin", "#!/bin/sh\nexit 0\n")
				return setupResult{dir: dir, cleanup: func() {}}
			},
			resetAvailable: true,
			testFunc: func(t *testing.T, s setupResult) {
				distroboxExportAvailable = IsDistroboxExportAvailable()
				err := ExportDistroboxBinary("failbin")
				if err == nil || !strings.Contains(err.Error(), "failed to export binary") {
					t.Fatalf("expected export binary error, got: %v", err)
				}
			},
		},
		{
			name: "IsDistroboxApplicationExported and ExportDistroboxApplication - success",
			prepare: func(t *testing.T, dir string) setupResult {
				marker := filepath.Join(dir, "app_exported.marker")
				body := `#!/bin/sh
if [ "$1" = "--list-apps" ]; then
	echo "myapp"
	exit 0
fi
if [ "$1" = "--app" ]; then
	echo "app exporting $2" > "` + marker + `"
	exit 0
fi
exit 0
`
				makeDistroboxScript(t, dir, body)
				return setupResult{dir: dir, marker: marker, cleanup: func() {}}
			},
			resetAvailable: true,
			testFunc: func(t *testing.T, s setupResult) {
				distroboxExportAvailable = IsDistroboxExportAvailable()
				ok, err := IsDistroboxApplicationExported("myapp")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !ok {
					t.Fatalf("expected myapp to be reported as exported")
				}
				// Export should no-op (already exported)
				if err := ExportDistroboxApplication("myapp"); err != nil {
					t.Fatalf("expected no error when already exported, got: %v", err)
				}
				// marker should be absent because exporter not invoked
				if _, err := os.Stat(s.marker); !os.IsNotExist(err) {
					t.Fatalf("expected app marker not to be created when already exported")
				}
			},
		},
		{
			name: "ExportDistroboxApplication - exporter invoked and fails",
			prepare: func(t *testing.T, dir string) setupResult {
				body := `#!/bin/sh
if [ "$1" = "--list-apps" ]; then
	exit 0
fi
if [ "$1" = "--app" ]; then
	echo "bad" >&2
	exit 3
fi
exit 0
`
				makeDistroboxScript(t, dir, body)
				return setupResult{dir: dir, cleanup: func() {}}
			},
			resetAvailable: true,
			testFunc: func(t *testing.T, s setupResult) {
				distroboxExportAvailable = IsDistroboxExportAvailable()
				err := ExportDistroboxApplication("failapp")
				if err == nil || !strings.Contains(err.Error(), "failed to export application") {
					t.Fatalf("expected failure exporting application, got: %v", err)
				}
			},
		},
		{
			name: "ExportDistroboxBinaries - distrobox not available returns error slice",
			prepare: func(t *testing.T, dir string) setupResult {
				// ensure distrobox not available
				return setupResult{dir: dir, cleanup: func() {}}
			},
			resetAvailable: false,
			testFunc: func(t *testing.T, s setupResult) {
				distroboxExportAvailable = false
				errs := ExportDistroboxBinaries([]string{"a", "b"})
				if len(errs) == 0 {
					t.Fatalf("expected error slice when distrobox not available")
				}
				found := false
				for _, e := range errs {
					if e != nil && strings.Contains(e.Error(), "distrobox-export command is not available") {
						found = true
					}
				}
				if !found {
					t.Fatalf("expected to find distrobox not available error in slice")
				}
			},
		},
		{
			name: "ExportDistroboxApplications - distrobox not available returns error slice",
			prepare: func(t *testing.T, dir string) setupResult {
				return setupResult{dir: dir, cleanup: func() {}}
			},
			resetAvailable: false,
			testFunc: func(t *testing.T, s setupResult) {
				distroboxExportAvailable = false
				errs := ExportDistroboxApplications([]string{"a"})
				if len(errs) == 0 {
					t.Fatalf("expected error slice when distrobox not available")
				}
				if errs[0] == nil || !strings.Contains(errs[0].Error(), "distrobox-export command is not available") {
					t.Fatalf("expected distrobox not available error, got: %v", errs[0])
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// avoid parallelism because tests modify PATH and package variable
			dir := t.TempDir()
			setup := tt.prepare(t, dir)
			defer func() {
				// restore PATH and reset distroboxExportAvailable to its initial detection
				_ = os.Setenv("PATH", origPath)
				distroboxExportAvailable = IsDistroboxExportAvailable()
				if setup.cleanup != nil {
					setup.cleanup()
				}
			}()

			// if prepare created distrobox-export, ensure PATH includes dir
			// but if not, still set PATH to not include it: keep origPath or replace
			// Prepend dir so our scripts are found first when present
			_ = os.Setenv("PATH", dir+string(os.PathListSeparator)+origPath)

			if tt.resetAvailable {
				distroboxExportAvailable = IsDistroboxExportAvailable()
			}

			tt.testFunc(t, setup)
		})
	}
}
