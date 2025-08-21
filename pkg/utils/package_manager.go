package utils

import (
	"bytes"
	"fmt"
	"os/exec"

	"go.uber.org/zap"
)

var (
	VSCODE_PACKAGE_MANAGER = &PackageManager{
		Name:             "code",
		InstallCmd:       "--install-extension",
		NoInteractiveArg: strPtr("--force"),
		MultiInstall:     false,
		SudoRequired:     false,
	}
)

type PackageManager struct {
	Name             string  `yaml:"name"`
	InstallCmd       string  `yaml:"install_cmd"`
	NoInteractiveArg *string `yaml:"no_interactive_arg,omitempty"`
	MultiInstall     bool    `yaml:"multi_install,omitempty"`
	SudoRequired     bool    `yaml:"sudo_required,omitempty"`
}

func (pm *PackageManager) Install(packages []string) []error {
	if pm == nil {
		return []error{fmt.Errorf("package manager is not specified or unsupported")}
	}

	// Multi-install logic
	if pm.MultiInstall {
		var cmd *exec.Cmd
		if pm.SudoRequired {
			cmd = exec.Command("sudo", pm.Name, pm.InstallCmd)
		} else {
			cmd = exec.Command(pm.Name, pm.InstallCmd)
		}
		zap.L().Info("Installing packages", zap.Strings("packages", packages), zap.String("package_manager", pm.Name))
		cmd.Args = append(cmd.Args, packages...)

		if pm.NoInteractiveArg != nil {
			cmd.Args = append(cmd.Args, *pm.NoInteractiveArg)
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		zap.L().Debug("Running command", zap.String("command", cmd.String()))
		err := cmd.Run()
		if err != nil {
			return []error{fmt.Errorf("failed to install packages using %s: %w, stderr: %s", pm.Name, err, stderr.String())}
		}
		zap.L().Info("Successfully installed packages", zap.Strings("packages", packages), zap.String("package_manager", pm.Name))
		return nil
	}

	// Single-install logic
	var errorChan = make(chan error, len(packages))
	for _, pkg := range packages {
		zap.L().Info("Installing package", zap.String("package", pkg), zap.String("package_manager", pm.Name))
		var cmd *exec.Cmd
		if pm.SudoRequired {
			cmd = exec.Command("sudo", pm.Name, pm.InstallCmd)
		} else {
			cmd = exec.Command(pm.Name, pm.InstallCmd)
		}
		cmd.Args = append(cmd.Args, pkg)

		if pm.NoInteractiveArg != nil {
			cmd.Args = append(cmd.Args, *pm.NoInteractiveArg)
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		zap.L().Debug("Running command", zap.String("command", cmd.String()))
		if err := cmd.Run(); err != nil {
			errorChan <- fmt.Errorf("failed to install package %s using %s: %w, stderr: %s", pkg, pm.Name, err, stderr.String())
		} else {
			zap.L().Info("Successfully installed package", zap.String("package", pkg), zap.String("package_manager", pm.Name))
		}
	}

	close(errorChan)
	return MergeErrors(errorChan)
}
