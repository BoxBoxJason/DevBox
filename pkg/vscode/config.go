package vscode

import (
	"devbox/pkg/packagemanager"
	"devbox/pkg/utils"
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
	SystemVSCode *VSCode = &VSCode{SettingsFile: nil}

	VSCODE_PACKAGE_MANAGER = &packagemanager.PackageManager{
		Name:         "code",
		InstallCmd:   "--install-extension",
		MultiInstall: false,
		SudoRequired: false,
	}
)

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
	var defaultPath string
	switch runtime.GOOS {
	case "windows":
		// Default Windows path
		appData := os.Getenv("APPDATA")
		defaultPath = filepath.Join(appData, "Code", "User", "settings.json")
		if appData != "" {
			paths = append(paths, defaultPath)
		}

	case "darwin":
		// Default macOS path
		defaultPath = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Code", "User", "settings.json")
		paths = append(paths, defaultPath)

	case "linux":
		defaultPath = filepath.Join(os.Getenv("HOME"), ".config", "Code", "User", "settings.json")
		// Default Linux path
		paths = append(paths, defaultPath)

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
	// If no valid path is found, issue a warning and set the default path
	zap.L().Warn("VSCode settings.json file not found in standard locations, using default path for OS", zap.String("default_path", defaultPath), zap.String("OS", runtime.GOOS))

	zap.L().Warn("Creating directories if they do not exist", zap.String("path", defaultPath))

	// Ensure the file exists
	code.SettingsFile = &defaultPath
	return utils.CreateFileIfNotExists(defaultPath, []byte("{}"))
}

// UpdateSettings updates the VS Code settings.json file with the provided settings.
func (code *VSCode) UpdateSettings(newSettings map[string]any) error {
	if code.SettingsFile == nil {
		err := code.FindVSCodeSettings()
		if err != nil {
			errMsg := fmt.Errorf("failed to locate VSCode settings file: %w", err)
			zap.L().Error(errMsg.Error())
			return errMsg
		}
	}

	// Read existing settings
	zap.L().Debug("Reading existing VSCode settings", zap.String("file", *code.SettingsFile))
	existingSettings := make(map[string]any)
	data, err := os.ReadFile(*code.SettingsFile)
	if err != nil {
		errMsg := fmt.Errorf("failed to read settings file: %w", err)
		zap.L().Error(errMsg.Error(), zap.String("file", *code.SettingsFile))
		return errMsg
	}

	if err := json.Unmarshal(data, &existingSettings); err != nil {
		errMsg := fmt.Errorf("failed to parse existing settings: %w", err)
		zap.L().Error(errMsg.Error(), zap.String("file", *code.SettingsFile))
		return errMsg
	}

	// Update settings with new values
	maps.Copy(existingSettings, newSettings)

	// Write updated settings back to the file
	updatedData, err := json.MarshalIndent(existingSettings, "", "  ")
	if err != nil {
		errMsg := fmt.Errorf("failed to serialize updated settings: %w", err)
		zap.L().Error(errMsg.Error(), zap.String("file", *code.SettingsFile))
		return errMsg
	}

	zap.L().Debug("Writing updated VSCode settings", zap.String("file", *code.SettingsFile))
	if err := os.WriteFile(*code.SettingsFile, updatedData, 0600); err != nil {
		errMsg := fmt.Errorf("failed to write updated settings to file: %w", err)
		zap.L().Error(errMsg.Error(), zap.String("file", *code.SettingsFile))
		return errMsg
	}

	return nil
}
