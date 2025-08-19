package utils

import (
	"errors"
	"os/exec"
)

// ExportDistroboxBinary exports a binary from the distrobox to the host system.
// It returns an error if the export fails.
func ExportDistroboxBinary(binaryName string) error {
	// Determine the path of the binary in the distrobox
	distroboxPath, err := exec.LookPath(binaryName)
	if err != nil {
		return err
	}

	// Check if the distrobox-export command is available
	_, err = exec.LookPath("distrobox-export")
	if err != nil {
		return errors.New("distrobox-export command not found")
	}

	// Construct the export command
	cmd := exec.Command("distrobox-export", "--bin", binaryName)
	cmd.Args = append(cmd.Args, distroboxPath)
	cmd.Stdout = nil // Redirect output to nil
	cmd.Stderr = nil // Redirect error to nil
	if err := cmd.Run(); err != nil {
		return errors.New("failed to export binary: " + err.Error())
	}
	return nil
}

// ExportDistroboxApplication exports an application from the distrobox to the host system.
// It returns an error if the export fails.
func ExportDistroboxApplication(appName string) error {
	// Check if the distrobox-export command is available
	_, err := exec.LookPath("distrobox-export")
	if err != nil {
		return errors.New("distrobox-export command not found")
	}

	// Construct the export command
	cmd := exec.Command("distrobox-export", "--app", appName)
	cmd.Stdout = nil // Redirect output to nil
	cmd.Stderr = nil // Redirect error to nil
	if err := cmd.Run(); err != nil {
		return errors.New("failed to export application: " + err.Error())
	}
	return nil
}
