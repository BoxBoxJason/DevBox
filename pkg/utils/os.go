package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

var (
	PACKAGE_MANAGERS = [][]string{{"apt", "install", "-y"}, {"dnf", "install", "-y"}, {"microdnf", "install", "-y"}, {"yum", "install", "-y"}, {"apk", "add"}, {"brew", "install"}, {"pacman", "-S", "--noconfirm"}, {"zypper", "install", "--non-interactive"}, {"port", "install"}, {"nix-env", "-i"}, {"flatpak", "install", "-y"}, {"snap", "install"}}
)

// DetectPackageManager attempts to detect the package manager used by the system.
// It returns the name of the package manager or an error if detection fails.
// The package manager must be able to support the syntax "package install <package_name> -y".
func DetectPackageManager() ([]string, error) {
	for _, pm := range PACKAGE_MANAGERS {
		if _, err := exec.LookPath(pm[0]); err == nil {
			return pm, nil
		}
	}
	return []string{}, fmt.Errorf("no supported package manager found")
}

// InstallPackage installs a package using the specified package manager.
// It returns an error if the installation fails.
func InstallPackage(packageManager *[]string, packageName string) error {
	if packageManager == nil || len(*packageManager) == 0 {
		return fmt.Errorf("package manager is not specified or unsupported")
	}

	// Check if the package already exists
	if _, err := exec.LookPath(packageName); err == nil {
		return nil // Package already exists, no need to install
	}
	cmd := exec.Command((*packageManager)[0], append((*packageManager)[1:], packageName)...)
	cmd.Stdout = nil // Redirect output to nil
	cmd.Stderr = nil // Redirect error to nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install package %s using %s: %w", packageName, packageManager, err)
	}
	return nil
}

// IsDistroboxPackageExported checks if a package is exported from the distrobox.
func IsDistroboxPackageExported(packageName string) (bool, error) {
	cmd := exec.Command("distrobox-export", "--list-binaries")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list exported binaries: %w", err)
	}

	return strings.Contains(string(output), packageName), nil
}
