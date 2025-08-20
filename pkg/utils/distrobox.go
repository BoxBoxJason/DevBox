package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

const (
	// DISTROBOX_EXPORT_COMMAND is the command used to export binaries and applications from the distro
	DISTROBOX_EXPORT_COMMAND = "distrobox-export"
)

func init() {
	// Check if the distrobox-export command is available
	_, err := exec.LookPath(DISTROBOX_EXPORT_COMMAND)
	if err != nil {
		panic(fmt.Sprintf("%s command not found. Please install %s. Or make sure you are inside the distrobox", DISTROBOX_EXPORT_COMMAND, DISTROBOX_EXPORT_COMMAND))
	}
}

// IsDistroboxBinaryExported checks if a package is exported from the distrobox.
func IsDistroboxBinaryExported(packageName string) (bool, error) {
	cmd := exec.Command(DISTROBOX_EXPORT_COMMAND, "--list-binaries")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list exported binaries: %w", err)
	}

	return strings.Contains(string(output), packageName), nil
}

// ExportDistroboxBinaries exports a list of binaries from the distrobox to the host system.
// It returns an error if the export fails.
func ExportDistroboxBinaries(binaries []string) error {
	for _, binary := range binaries {
		if err := ExportDistroboxBinary(binary); err != nil {
			return fmt.Errorf("failed to export binary %s: %w", binary, err)
		}
	}
	return nil
}

// ExportDistroboxBinary exports a binary from the distrobox to the host system.
// It returns an error if the export fails.
func ExportDistroboxBinary(binaryName string) error {
	// Check if the binary is already exported
	exported, err := IsDistroboxBinaryExported(binaryName)
	if err != nil {
		return fmt.Errorf("failed to check if binary is exported: %w", err)
	} else if exported {
		return nil // Binary is already exported, no need to export again
	}

	// Determine the path of the binary in the distrobox
	binaryPath, err := exec.LookPath(binaryName)
	if err != nil {
		return err
	}

	// Construct the export command
	cmd := exec.Command(DISTROBOX_EXPORT_COMMAND, "--bin", binaryPath)
	cmd.Stdout = nil // Redirect output to nil
	cmd.Stderr = nil // Redirect error to nil
	if err := cmd.Run(); err != nil {
		return errors.New("failed to export binary: " + err.Error())
	}
	return nil
}

// IsDistroboxApplicationExported checks if an application is exported from the distrobox.
func IsDistroboxApplicationExported(appName string) (bool, error) {
	cmd := exec.Command(DISTROBOX_EXPORT_COMMAND, "--list-apps")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to list exported applications: %w", err)
	}

	return strings.Contains(string(output), appName), nil
}

// ExportDistroboxApplications exports a list of applications from the distrobox to the host system.
// It returns an error if the export fails.
func ExportDistroboxApplications(apps []string) error {
	for _, app := range apps {
		if err := ExportDistroboxApplication(app); err != nil {
			return fmt.Errorf("failed to export application %s: %w", app, err)
		}
	}
	return nil
}

// ExportDistroboxApplication exports an application from the distrobox to the host system.
// It returns an error if the export fails.
func ExportDistroboxApplication(appName string) error {
	// Check if the application is already exported
	exported, err := IsDistroboxApplicationExported(appName)
	if err != nil {
		return fmt.Errorf("failed to check if application is exported: %w", err)
	} else if exported {
		return nil // Application is already exported, no need to export again
	}

	// Construct the export command
	cmd := exec.Command(DISTROBOX_EXPORT_COMMAND, "--app", appName)
	cmd.Stdout = nil // Redirect output to nil
	cmd.Stderr = nil // Redirect error to nil
	if err := cmd.Run(); err != nil {
		return errors.New("failed to export application: " + err.Error())
	}
	return nil
}
