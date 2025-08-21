package utils

import (
	"fmt"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

var (
	SystemPackageManager *PackageManager = nil

	PACKAGE_MANAGERS = []PackageManager{
		{
			Name:             "apt",
			InstallCmd:       "install",
			NoInteractiveArg: strPtr("-y"),
			MultiInstall:     true,
			SudoRequired:     true,
		},
		{
			Name:             "dnf",
			InstallCmd:       "install",
			NoInteractiveArg: strPtr("-y"),
			MultiInstall:     true,
			SudoRequired:     true,
		},
		{
			Name:             "microdnf",
			InstallCmd:       "install",
			NoInteractiveArg: strPtr("-y"),
			MultiInstall:     true,
			SudoRequired:     true,
		},
		{
			Name:             "yum",
			InstallCmd:       "install",
			NoInteractiveArg: strPtr("-y"),
			MultiInstall:     true,
			SudoRequired:     true,
		},
		{
			Name:             "apk",
			InstallCmd:       "add",
			NoInteractiveArg: nil,
			MultiInstall:     true,
			SudoRequired:     true,
		},
		{
			Name:             "brew",
			InstallCmd:       "install",
			NoInteractiveArg: nil,
			MultiInstall:     true,
			SudoRequired:     false,
		},
		{
			Name:             "pacman",
			InstallCmd:       "-S",
			NoInteractiveArg: strPtr("--noconfirm"),
			MultiInstall:     true,
			SudoRequired:     true,
		},
		{
			Name:             "zypper",
			InstallCmd:       "install",
			NoInteractiveArg: strPtr("--non-interactive"),
			MultiInstall:     true,
			SudoRequired:     true,
		},
		{
			Name:             "port",
			InstallCmd:       "install",
			NoInteractiveArg: nil,
			MultiInstall:     true,
			SudoRequired:     true,
		},
		{
			Name:             "nix-env",
			InstallCmd:       "-i",
			NoInteractiveArg: nil,
			MultiInstall:     true,
			SudoRequired:     false,
		},
		{
			Name:             "flatpak",
			InstallCmd:       "install",
			NoInteractiveArg: strPtr("-y"),
			MultiInstall:     false,
			SudoRequired:     false,
		},
		{
			Name:             "snap",
			InstallCmd:       "install",
			NoInteractiveArg: nil,
			MultiInstall:     false,
			SudoRequired:     false,
		},
	}
)

func init() {
	// Detect the package manager at initialization
	pm, err := DetectPackageManager()
	if err != nil {
		zap.L().Fatal("Failed to detect package manager", zap.Error(err))
	} else {
		SystemPackageManager = pm
	}
}

// DetectPackageManager attempts to detect the package manager used by the system.
// It returns the name of the package manager or an error if detection fails.
// The package manager must be able to support the syntax "package install <package_name> -y".
func DetectPackageManager() (*PackageManager, error) {
	zap.L().Debug("Detecting package manager")
	for _, pm := range PACKAGE_MANAGERS {
		if _, err := exec.LookPath(pm.Name); err == nil {
			return &pm, nil
		}
	}
	return nil, fmt.Errorf("no supported package manager found")
}

// Getenv retrieves the value of the environment variable named by key.
// If the variable is not set, it returns the provided default value.
func Getenv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultValue
	}
	return value
}
