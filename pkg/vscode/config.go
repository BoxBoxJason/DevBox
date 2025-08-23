package vscode

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"go.uber.org/zap"
)

var (
	SystemVSCode *VSCode = nil
)

func init() {
	SystemVSCode = &VSCode{}
	if err := SystemVSCode.FindVSCodeSettings(); err != nil {
		zap.L().Error("Failed to locate VSCode settings.json", zap.Error(err))
	}
}

type VSCode struct {
	SettingsFile *string
}

// FindVSCodeSettings locates the path to the VS Code settings.json file across platforms.
func (code *VSCode) FindVSCodeSettings() error {
	if code.SettingsFile != nil && strings.TrimSpace(*code.SettingsFile) != "" {
		// If a specific file path is provided, use it directly
		if _, err := os.Stat(*code.SettingsFile); err == nil {
			return nil
		}
		return fmt.Errorf("provided settings.json file does not exist: %s", *code.SettingsFile)
	}

	var paths []string

	// Determine the OS and add potential paths
	zap.L().Debug("Detecting VSCode settings.json path for OS", zap.String("os", runtime.GOOS))
	switch runtime.GOOS {
	case "windows":
		// Default Windows path
		appData := os.Getenv("APPDATA")
		if appData != "" {
			paths = append(paths, filepath.Join(appData, "Code", "User", "settings.json"))
		}

	case "darwin":
		// Default macOS path
		paths = append(paths, filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Code", "User", "settings.json"))

	case "linux":
		// Default Linux path
		paths = append(paths, filepath.Join(os.Getenv("HOME"), ".config", "Code", "User", "settings.json"))

		// Flatpak-specific path
		paths = append(paths, filepath.Join(os.Getenv("HOME"), ".var", "app", "com.visualstudio.code", "config", "Code", "User", "settings.json"))
	}

	// Search through all potential paths
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			code.SettingsFile = &path
			return nil
		}
	}

	return fmt.Errorf("settings.json file not found in standard locations")
}

// UpdateSettings updates the VS Code settings.json file with the provided settings.
func (code *VSCode) UpdateSettings(newSettings map[string]any) error {
	if code.SettingsFile == nil {
		return fmt.Errorf("settings file path is not set")
	}

	// Read existing settings
	existingSettings := make(map[string]any)
	data, err := os.ReadFile(*code.SettingsFile)
	if err != nil {
		return fmt.Errorf("failed to read settings file: %w", err)
	}

	if err := json.Unmarshal(data, &existingSettings); err != nil {
		return fmt.Errorf("failed to parse existing settings: %w", err)
	}

	// Update settings with new values
	maps.Copy(existingSettings, newSettings)

	// Write updated settings back to the file
	updatedData, err := json.MarshalIndent(existingSettings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize updated settings: %w", err)
	}

	if err := os.WriteFile(*code.SettingsFile, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated settings to file: %w", err)
	}

	return nil
}
