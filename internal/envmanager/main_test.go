package envmanager

import (
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
	errs := em.AppendToEnvFile(map[string]string{"FOO": "bar"})
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
	errs := em.AppendToEnvFile(map[string]string{"BAZ": "qux"})
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
	errs := em.AppendToEnvFile(vars)
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
	errs := em.AppendToEnvFile(map[string]string{"X": "y"})
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
	errs := em.AppendToEnvFile(map[string]string{"NO": "WRITE"})
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

func Test_ReloadEnvFile_expandsEnvVars(t *testing.T) {
	t.Parallel()
	// ensure a clean environment for the keys we use
	keys := []string{"T1_TEST", "T2_TEST", "T3_TEST"}
	orig := make(map[string]string)
	for _, k := range keys {
		orig[k] = os.Getenv(k)
		_ = os.Unsetenv(k)
	}
	// restore at end
	defer func() {
		for _, k := range keys {
			if v := orig[k]; v != "" {
				_ = os.Setenv(k, v)
			} else {
				_ = os.Unsetenv(k)
			}
		}
	}()

	em := &EnvManager{
		variables: map[string]string{
			"T1_TEST": "value1",
			"T2_TEST": "$T1_TEST-suffix",
			"T3_TEST": "$NON_EXISTENT_VAR",
		},
	}

	errs := em.ReloadEnvFile()
	if errs != nil && len(errs) > 0 {
		t.Fatalf("expected no errors from ReloadEnvFile, got: %v", errs)
	}

	if got := os.Getenv("T1_TEST"); got != "value1" {
		t.Fatalf("T1_TEST expected 'value1', got %q", got)
	}
	if got := os.Getenv("T2_TEST"); got != "value1-suffix" {
		t.Fatalf("T2_TEST expected 'value1-suffix', got %q", got)
	}
	if got := os.Getenv("T3_TEST"); got != "" {
		t.Fatalf("T3_TEST expected empty string for missing var expansion, got %q", got)
	}
}

func Test_Set_noChange_and_withChange(t *testing.T) {
	t.Parallel()
	// no-change case
	{
		dir := t.TempDir()
		fp := filepath.Join(dir, "env_nochange.sh")
		initial := `export K="v"` + "\n"
		if err := os.WriteFile(fp, []byte(initial), 0600); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		em := &EnvManager{
			file:      fp,
			variables: map[string]string{"K": "v"},
		}
		before, _ := os.ReadFile(fp)
		if errs := em.Set(map[string]string{"K": "v"}); errs != nil && len(errs) > 0 {
			t.Fatalf("expected no errors for no-change Set, got: %v", errs)
		}
		after, _ := os.ReadFile(fp)
		if string(before) != string(after) {
			t.Fatalf("file changed on no-op Set; before=%q after=%q", string(before), string(after))
		}
	}

	// change case
	{
		dir := t.TempDir()
		fp := filepath.Join(dir, "env_change.sh")
		// start with empty file
		if err := os.WriteFile(fp, []byte(""), 0600); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		// ensure env var not set prior
		_ = os.Unsetenv("NEW_TEST_VAR")
		em := &EnvManager{
			file:      fp,
			variables: map[string]string{}, // start empty
		}
		if errs := em.Set(map[string]string{"NEW_TEST_VAR": "someval"}); errs != nil && len(errs) > 0 {
			t.Fatalf("expected no errors for Set, got: %v", errs)
		}
		// file should contain the export line
		b, _ := os.ReadFile(fp)
		content := string(b)
		if !strings.Contains(content, `export NEW_TEST_VAR="someval"`) {
			t.Fatalf("expected export line in file, got: %q", content)
		}
		// em.variables should have been updated
		if v, ok := em.variables["NEW_TEST_VAR"]; !ok || v != "someval" {
			t.Fatalf("expected em.variables to contain NEW_TEST_VAR=someval, got: %v", em.variables)
		}
		// environment should be set
		if got := os.Getenv("NEW_TEST_VAR"); got != "someval" {
			t.Fatalf("expected environment NEW_TEST_VAR=someval, got %q", got)
		}
		// cleanup
		_ = os.Unsetenv("NEW_TEST_VAR")
	}
}
