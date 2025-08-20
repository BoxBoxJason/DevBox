package utils

import (
	"fmt"
	"os/exec"
)

type PackageManager struct {
	Name             string  `yaml:"name"`
	InstallCmd       string  `yaml:"install_cmd"`
	NoInteractiveArg *string `yaml:"no_interactive_arg,omitempty"`
	MultiInstall     bool    `yaml:"multi_install,omitempty"` // Indicates if the package manager supports installing multiple packages in one command
}

func (pm *PackageManager) Install(packages []string) error {
	if pm == nil {
		return fmt.Errorf("package manager is not specified or unsupported")
	}

	// Multi-install logic
	if pm.MultiInstall {
		cmd := exec.Command(pm.Name, pm.InstallCmd)
		cmd.Args = append(cmd.Args, packages...)

		if pm.NoInteractiveArg != nil {
			cmd.Args = append(cmd.Args, *pm.NoInteractiveArg)
		}

		cmd.Stdout = nil // Redirect output to nil
		cmd.Stderr = nil // Redirect error to nil

		return cmd.Run()
	}

	// Single-install logic
	for _, pkg := range packages {
		cmd := exec.Command(pm.Name, pm.InstallCmd, pkg)
		if pm.NoInteractiveArg != nil {
			cmd.Args = append(cmd.Args, *pm.NoInteractiveArg)
		}

		cmd.Stdout = nil // Redirect output to nil
		cmd.Stderr = nil // Redirect error to nil

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install package %s using %s: %w", pkg, pm.Name, err)
		}
	}

	return nil
}
