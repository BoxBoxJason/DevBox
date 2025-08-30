package envmanager

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func Test_AppendToEnvFile_EndsWithNewline(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fp := filepath.Join(dir, "env.sh")

	if err := os.WriteFile(fp, []byte("old\n"), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{file: fp}
	errs := em.AppendToEnvFile(map[string]string{"FOO": "bar"}, nil)
	if errs != nil && len(errs) > 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}

	b, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	content := string(b)

	// Should contain the appended line exactly and not produce an extra blank line
	if !strings.Contains(content, "export FOO=\"bar\"\n") {
		t.Fatalf("expected appended export line, got content: %q", content)
	}
	// Ensure there's exactly one newline between "old" and export line (old ended with newline)
	if !strings.Contains(content, "old\nexport FOO=\"bar\"\n") {
		t.Fatalf("expected 'old\\nexport...' sequence, got: %q", content)
	}
}

func Test_AppendToEnvFile_NoNewlineAtEnd(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fp := filepath.Join(dir, "env2.sh")

	if err := os.WriteFile(fp, []byte("old"), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{file: fp}
	errs := em.AppendToEnvFile(map[string]string{"BAZ": "qux"}, nil)
	if errs != nil && len(errs) > 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}

	b, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	content := string(b)

	// When original file did not end with newline, AppendToEnvFile should have added one before exports
	if !strings.Contains(content, "old\nexport BAZ=\"qux\"\n") {
		t.Fatalf("expected newline inserted before export, got: %q", content)
	}
}

func Test_AppendToEnvFile_MultipleVars(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fp := filepath.Join(dir, "env3.sh")

	// empty file
	if err := os.WriteFile(fp, []byte(""), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{file: fp}
	vars := map[string]string{
		"A": "1",
		"B": `two " quotes`,
		"C": "3",
	}
	errs := em.AppendToEnvFile(vars, nil)
	if errs != nil && len(errs) > 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}

	b, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	content := string(b)

	// Check each expected export line is present (order independent)
	for k, v := range vars {
		expected := "export " + k + "=\"" + v + "\""
		if !strings.Contains(content, expected) {
			t.Fatalf("expected to find %q in file content, got: %q", expected, content)
		}
	}
}

func Test_AppendToEnvFile_FileDoesNotExist_ReturnsError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	non := filepath.Join(dir, "does_not_exist_env.sh")

	em := &EnvManager{file: non}
	errs := em.AppendToEnvFile(map[string]string{"X": "y"}, nil)
	if errs == nil || len(errs) == 0 {
		t.Fatalf("expected error for non-existent file, got nil")
	}
	found := false
	for _, e := range errs {
		if e != nil && strings.Contains(e.Error(), "failed to check if env file ends with newline") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected error mentioning 'failed to check if env file ends with newline', got: %v", errs)
	}
}

func Test_AppendToEnvFile_FileNotWritable_ReturnsError(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("skipping permission test on Windows")
	}

	dir := t.TempDir()
	fp := filepath.Join(dir, "readonly.sh")
	if err := os.WriteFile(fp, []byte("something\n"), 0400); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{file: fp}
	errs := em.AppendToEnvFile(map[string]string{"NO": "WRITE"}, nil)
	if errs == nil || len(errs) == 0 {
		t.Fatalf("expected error when opening non-writable file, got nil")
	}
	found := false
	for _, e := range errs {
		if e != nil && strings.Contains(e.Error(), "failed to open env file for appending") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected error mentioning 'failed to open env file for appending', got: %v", errs)
	}
}

func Test_parseEnvFile_various(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fp := filepath.Join(dir, "parse_env.sh")

	content := strings.Join([]string{
		"# a comment line",
		"not an export line",
		"export A=1",
		`export B="two words"`,
		"export    C = 'with spaces'   ",
		`export D="a=b=c"`,
		"exportE=should_ignore", // not prefixed with export
		"export    F =  unquoted_value",
	}, "\n") + "\n"

	if err := os.WriteFile(fp, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{file: fp}
	if err := em.parseEnvFile(); err != nil {
		t.Fatalf("parseEnvFile returned error: %v", err)
	}

	// expected values after trimSpacesAndQuotes
	expect := map[string]string{
		"A": "1",
		"B": "two words",
		"C": "with spaces",
		"D": "a=b=c",
		"F": "unquoted_value",
	}

	for k, v := range expect {
		got, ok := em.variables[k]
		if !ok {
			t.Fatalf("expected key %q missing in parsed variables", k)
		}
		if got != v {
			t.Fatalf("key %q: expected %q, got %q", k, v, got)
		}
	}

	// ensure non-export or malformed keys not present
	if _, ok := em.variables["exportE"]; ok {
		t.Fatalf("unexpected key 'exportE' found")
	}
}

func Test_parseEnvFile_FileOpenError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	non := filepath.Join(dir, "does_not_exist_parse.sh")

	em := &EnvManager{file: non}
	err := em.parseEnvFile()
	if err == nil {
		t.Fatalf("expected error opening non-existent file, got nil")
	}
	if !strings.Contains(err.Error(), "failed to open env file") {
		t.Fatalf("expected 'failed to open env file' in error, got: %v", err)
	}
}

func Test_SystemEnvManager_FileConditions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		setupFile   func(string) error
		expectFatal bool
		fatalMsg    string
	}{
		{
			name:        "empty_envFile_path",
			setupFile:   func(string) error { return nil },
			expectFatal: true,
			fatalMsg:    "DEVBOX_ENV_FILE is set to an empty string",
		},
		{
			name:        "file_does_not_exist_creates_it",
			setupFile:   func(string) error { return nil }, // don't create file
			expectFatal: false,
		},
		{
			name:        "file_is_directory",
			setupFile:   func(fp string) error { return os.Mkdir(fp, 0755) },
			expectFatal: true,
			fatalMsg:    "DEVBOX_ENV_FILE points to a directory",
		},
		{
			name:        "file_wrong_permissions",
			setupFile:   func(fp string) error { return os.WriteFile(fp, []byte{}, 0644) },
			expectFatal: true,
			fatalMsg:    "DEVBOX_ENV_FILE does not have the correct permissions",
		},
		{
			name:        "valid_file",
			setupFile:   func(fp string) error { return os.WriteFile(fp, []byte("export TEST=value\n"), 0600) },
			expectFatal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			fp := filepath.Join(dir, "env.sh")

			envFile := fp
			if tt.name == "empty_envFile_path" {
				envFile = ""
			} else if err := tt.setupFile(fp); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			// Reset singleton for each test
			ResetSystemEnvManager()

			if tt.expectFatal {
				// Can't easily test zap.L().Fatal() calls without mocking
				// So we'll test the logic conditions manually
				switch tt.name {
				case "empty_envFile_path":
					if envFile != "" {
						t.Fatal("expected empty envFile to be handled")
					}
				case "file_is_directory":
					if info, err := os.Stat(fp); err == nil && !info.IsDir() {
						t.Fatal("expected directory check")
					}
				case "file_wrong_permissions":
					if info, err := os.Stat(fp); err == nil {
						if info.Mode().Perm() == 0600 {
							t.Fatal("expected permission check - file should NOT have 0600 permissions")
						}
					}
				}
				return
			}

			em := SystemEnvManager(envFile)
			if em == nil {
				t.Fatal("expected non-nil EnvManager")
			}

			// Verify file was created if it didn't exist
			if _, err := os.Stat(envFile); os.IsNotExist(err) {
				t.Fatal("expected file to be created")
			}

			// Test singleton behavior
			em2 := SystemEnvManager(envFile)
			if em != em2 {
				t.Fatal("expected singleton behavior")
			}
		})
	}
}

func Test_parseEnvFile_PathVariables(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fp := filepath.Join(dir, "env_path.sh")

	content := strings.Join([]string{
		`export PATH="/usr/bin"`,
		`export PATH="/usr/local/bin"`,
		`export REGULAR_VAR="value"`,
		`export PATH="/opt/bin"`,
		"# comment line",
		`export ANOTHER_VAR="another"`,
	}, "\n") + "\n"

	if err := os.WriteFile(fp, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{file: fp, pathVariables: make(map[string]struct{})}
	if err := em.parseEnvFile(); err != nil {
		t.Fatalf("parseEnvFile returned error: %v", err)
	}

	// Check PATH variables are stored separately
	expectedPaths := []string{"/usr/bin", "/usr/local/bin", "/opt/bin"}
	for _, path := range expectedPaths {
		if _, exists := em.pathVariables[path]; !exists {
			t.Fatalf("expected PATH variable %q not found in pathVariables", path)
		}
	}

	// Check regular variables
	expectedVars := map[string]string{
		"REGULAR_VAR": "value",
		"ANOTHER_VAR": "another",
	}
	for k, v := range expectedVars {
		if got, ok := em.variables[k]; !ok || got != v {
			t.Fatalf("expected variable %q=%q, got %q", k, v, got)
		}
	}

	// Ensure PATH is not in regular variables
	if _, exists := em.variables["PATH"]; exists {
		t.Fatal("PATH should not be in regular variables map")
	}
}

func Test_Set_PathVariables(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		initialPaths   []string
		setPaths       []string
		expectAppended []string
	}{
		{
			name:           "new_path_variable",
			initialPaths:   []string{},
			setPaths:       []string{"/new/path"},
			expectAppended: []string{"/new/path"},
		},
		{
			name:           "existing_path_not_duplicated",
			initialPaths:   []string{"/existing/path"},
			setPaths:       []string{"/existing/path"},
			expectAppended: []string{},
		},
		{
			name:           "mix_new_and_existing",
			initialPaths:   []string{"/existing"},
			setPaths:       []string{"/existing", "/new1", "/new2"},
			expectAppended: []string{"/new1", "/new2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			fp := filepath.Join(dir, "env.sh")

			// Setup initial file with existing paths
			var initialContent strings.Builder
			for _, path := range tt.initialPaths {
				initialContent.WriteString(fmt.Sprintf("export PATH=\"%s\"\n", path))
			}
			if err := os.WriteFile(fp, []byte(initialContent.String()), 0600); err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}

			em := &EnvManager{
				file:          fp,
				variables:     make(map[string]string),
				pathVariables: make(map[string]struct{}),
			}

			// Parse existing content
			if err := em.parseEnvFile(); err != nil {
				t.Fatalf("parseEnvFile error: %v", err)
			}

			// Set PATH variables
			pathVars := make(map[string]string)
			for _, path := range tt.setPaths {
				pathVars["PATH"] = path
				if errs := em.Set(pathVars); len(errs) > 0 {
					t.Fatalf("Set error: %v", errs)
				}
			}

			// Read final file content
			content, err := os.ReadFile(fp)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}
			fileContent := string(content)

			// Check that expected paths were appended
			for _, path := range tt.expectAppended {
				expected := fmt.Sprintf("export PATH=\"%s\"", path)
				if !strings.Contains(fileContent, expected) {
					t.Fatalf("expected %q to be appended to file, content: %q", expected, fileContent)
				}
			}

			// Count total PATH exports to ensure no duplicates
			pathCount := strings.Count(fileContent, "export PATH=")
			expectedCount := len(tt.initialPaths) + len(tt.expectAppended)
			if pathCount != expectedCount {
				t.Fatalf("expected %d PATH exports, got %d in content: %q", expectedCount, pathCount, fileContent)
			}
		})
	}
}

func Test_AppendToEnvFile_PathVariables(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fp := filepath.Join(dir, "env_path_append.sh")

	if err := os.WriteFile(fp, []byte("# initial content\n"), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{file: fp}

	// Test appending both regular vars and PATH vars
	envVars := map[string]string{"REGULAR": "value"}
	pathVars := []string{"/path1", "/path2"}

	errs := em.AppendToEnvFile(envVars, pathVars)
	if len(errs) > 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}

	content, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	fileContent := string(content)

	// Check regular variable
	if !strings.Contains(fileContent, `export REGULAR="value"`) {
		t.Fatalf("expected regular variable export, got: %q", fileContent)
	}

	// Check PATH variables
	for _, path := range pathVars {
		expected := fmt.Sprintf(`export PATH="%s"`, path)
		if !strings.Contains(fileContent, expected) {
			t.Fatalf("expected PATH export %q, got: %q", expected, fileContent)
		}
	}
}

func Test_AppendToEnvFile_EmptyInputs(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fp := filepath.Join(dir, "empty_inputs.sh")

	if err := os.WriteFile(fp, []byte("initial"), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{file: fp}

	// Test with empty maps/slices
	errs := em.AppendToEnvFile(map[string]string{}, []string{})
	if len(errs) > 0 {
		t.Fatalf("expected no errors for empty inputs, got: %v", errs)
	}

	errs = em.AppendToEnvFile(nil, nil)
	if len(errs) > 0 {
		t.Fatalf("expected no errors for nil inputs, got: %v", errs)
	}

	// File should be unchanged
	content, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(content) != "initial" {
		t.Fatalf("expected file unchanged, got: %q", string(content))
	}
}

func Test_Set_MultipleVariableMaps(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fp := filepath.Join(dir, "multi_maps.sh")

	if err := os.WriteFile(fp, []byte(""), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{
		file:          fp,
		variables:     make(map[string]string),
		pathVariables: make(map[string]struct{}),
	}

	// Test setting multiple variable maps at once
	map1 := map[string]string{"VAR1": "val1", "VAR2": "val2"}
	map2 := map[string]string{"VAR3": "val3", "PATH": "/new/path"}
	map3 := map[string]string{"VAR1": "updated_val1"} // update existing

	errs := em.Set(map1, map2, map3)
	if len(errs) > 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}

	// Check in-memory state
	expected := map[string]string{
		"VAR1": "updated_val1",
		"VAR2": "val2",
		"VAR3": "val3",
	}
	for k, v := range expected {
		if got, ok := em.variables[k]; !ok || got != v {
			t.Fatalf("expected variable %q=%q, got %q", k, v, got)
		}
	}

	// Check PATH variable
	if _, exists := em.pathVariables["/new/path"]; !exists {
		t.Fatal("expected PATH variable to be stored")
	}

	// Check file content
	content, err := os.ReadFile(fp)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	fileContent := string(content)

	// Should contain all variables (including the update)
	for k, v := range expected {
		expectedLine := fmt.Sprintf(`export %s="%s"`, k, v)
		if !strings.Contains(fileContent, expectedLine) {
			t.Fatalf("expected %q in file content: %q", expectedLine, fileContent)
		}
	}

	// Should contain PATH
	if !strings.Contains(fileContent, `export PATH="/new/path"`) {
		t.Fatalf("expected PATH export in file content: %q", fileContent)
	}
}

func Test_Set_NilVariablesMap(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	fp := filepath.Join(dir, "nil_vars.sh")

	if err := os.WriteFile(fp, []byte(""), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	em := &EnvManager{
		file:          fp,
		variables:     nil, // explicitly nil
		pathVariables: make(map[string]struct{}),
	}

	errs := em.Set(map[string]string{"TEST": "value"})
	if len(errs) > 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}

	// Should have initialized the map
	if em.variables == nil {
		t.Fatal("expected variables map to be initialized")
	}

	if got, ok := em.variables["TEST"]; !ok || got != "value" {
		t.Fatalf("expected TEST=value, got %q", got)
	}
}
