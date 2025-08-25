package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_ReadFileLines(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"basic lines", args{"one\ntwo\nthree\n"}, []string{"one", "two", "three"}},
		{"trims spaces and ignores empty lines", args{"  a  \n\n  b\n   \n c  \n"}, []string{"a", "b", "c"}},
		{"only empty and spaces", args{"   \n\n\t\n"}, []string{}},
		{"single line no newline", args{" single line "}, []string{"single line"}},
		{"lines with internal quotes and spaces", args{"  'hi'  \n\"hello\"\n mid \" quote \n"}, []string{"'hi'", "\"hello\"", "mid \" quote"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			fp := filepath.Join(dir, "testfile.txt")
			if err := os.WriteFile(fp, []byte(tt.args.content), 0644); err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}
			got, err := ReadFileLines(fp)
			if err != nil {
				t.Fatalf("ReadFileLines() returned error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ReadFileLines() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("non-existent file returns error", func(t *testing.T) {
		t.Parallel()
		non := filepath.Join(t.TempDir(), "does_not_exist.txt")
		if _, err := ReadFileLines(non); err == nil {
			t.Fatalf("expected error for non-existent file, got nil")
		}
	})
}

func Test_FileEndsWithNewline(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty file", args{""}, true},
		{"ends with newline", args{"a\nb\n"}, true},
		{"no newline at end", args{"a\nb"}, false},
		{"single newline only", args{"\n"}, true},
		{"single char no newline", args{"x"}, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			fp := filepath.Join(dir, "f.txt")
			if err := os.WriteFile(fp, []byte(tt.args.content), 0644); err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}
			got, err := FileEndsWithNewline(fp)
			if err != nil {
				t.Fatalf("FileEndsWithNewline() returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("FileEndsWithNewline() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("non-existent file returns error", func(t *testing.T) {
		t.Parallel()
		non := filepath.Join(t.TempDir(), "nofile")
		if _, err := FileEndsWithNewline(non); err == nil {
			t.Fatalf("expected error for non-existent file, got nil")
		}
	})
}

func Test_CreateFileIfNotExists(t *testing.T) {
	fileContent := []byte("initial content")
	t.Run("creates file if not exists", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		fp := filepath.Join(dir, "newfile.txt")
		if err := CreateFileIfNotExists(fp, fileContent); err != nil {
			t.Fatalf("CreateFileIfNotExists() returned error: %v", err)
		}
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			t.Fatalf("file was not created")
		}
	})

	t.Run("does not overwrite existing file", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		fp := filepath.Join(dir, "existingfile.txt")
		content := []byte("existing content")
		if err := os.WriteFile(fp, content, 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}
		if err := CreateFileIfNotExists(fp, fileContent); err != nil {
			t.Fatalf("CreateFileIfNotExists() returned error: %v", err)
		}
		data, err := os.ReadFile(fp)
		if err != nil {
			t.Fatalf("failed to read file: %v", err)
		}
		if string(data) != string(content) {
			t.Fatalf("file content was altered, got %s, want %s", string(data), string(content))
		}
	})
}
