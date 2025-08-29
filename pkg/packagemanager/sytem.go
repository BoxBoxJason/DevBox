package packagemanager

import (
	"devbox/pkg/utils"
	"fmt"
	"os/exec"

	"go.uber.org/zap"
)

var (
	SystemPackageManager *PackageManager = nil

	SYSTEM_PACKAGE_MANAGERS = []*PackageManager{
		APT_PACKAGE_MANAGER,
		DNF_PACKAGE_MANAGER,
		MICRODNF_PACKAGE_MANAGER,
		YUM_PACKAGE_MANAGER,
		APK_PACKAGE_MANAGER,
		BREW_PACKAGE_MANAGER,
		PACMAN_PACKAGE_MANAGER,
		ZYPPER_PACKAGE_MANAGER,
		PORT_PACKAGE_MANAGER,
		NIX_ENV_PACKAGE_MANAGER,
		FLATPAK_PACKAGE_MANAGER,
		SNAP_PACKAGE_MANAGER,
	}

	APT_PACKAGE_MANAGER = &PackageManager{
		Name:             "apt",
		InstallCmd:       "install",
		NoInteractiveArg: utils.StrPtr("-y"),
		MultiInstall:     true,
		SudoRequired:     true,
	}

	DNF_PACKAGE_MANAGER = &PackageManager{
		Name:             "dnf",
		InstallCmd:       "install",
		NoInteractiveArg: utils.StrPtr("-y"),
		MultiInstall:     true,
		SudoRequired:     true,
	}

	MICRODNF_PACKAGE_MANAGER = &PackageManager{
		Name:             "microdnf",
		InstallCmd:       "install",
		NoInteractiveArg: utils.StrPtr("-y"),
		MultiInstall:     true,
		SudoRequired:     true,
	}

	YUM_PACKAGE_MANAGER = &PackageManager{
		Name:             "yum",
		InstallCmd:       "install",
		NoInteractiveArg: utils.StrPtr("-y"),
		MultiInstall:     true,
		SudoRequired:     true,
	}

	APK_PACKAGE_MANAGER = &PackageManager{
		Name:             "apk",
		InstallCmd:       "add",
		NoInteractiveArg: nil,
		MultiInstall:     true,
		SudoRequired:     true,
	}

	BREW_PACKAGE_MANAGER = &PackageManager{
		Name:             "brew",
		InstallCmd:       "install",
		NoInteractiveArg: nil,
		MultiInstall:     true,
		SudoRequired:     false,
	}

	PACMAN_PACKAGE_MANAGER = &PackageManager{
		Name:             "pacman",
		InstallCmd:       "-S",
		NoInteractiveArg: utils.StrPtr("--noconfirm"),
		MultiInstall:     true,
		SudoRequired:     true,
	}

	ZYPPER_PACKAGE_MANAGER = &PackageManager{
		Name:             "zypper",
		InstallCmd:       "install",
		NoInteractiveArg: utils.StrPtr("--non-interactive"),
		MultiInstall:     true,
		SudoRequired:     true,
	}

	PORT_PACKAGE_MANAGER = &PackageManager{
		Name:             "port",
		InstallCmd:       "install",
		NoInteractiveArg: nil,
		MultiInstall:     true,
		SudoRequired:     true,
	}

	NIX_ENV_PACKAGE_MANAGER = &PackageManager{
		Name:             "nix-env",
		InstallCmd:       "-i",
		NoInteractiveArg: nil,
		MultiInstall:     true,
		SudoRequired:     false,
	}

	FLATPAK_PACKAGE_MANAGER = &PackageManager{
		Name:             "flatpak",
		InstallCmd:       "install",
		NoInteractiveArg: utils.StrPtr("-y"),
		MultiInstall:     false,
		SudoRequired:     false,
	}

	SNAP_PACKAGE_MANAGER = &PackageManager{
		Name:             "snap",
		InstallCmd:       "install",
		NoInteractiveArg: nil,
		MultiInstall:     false,
		SudoRequired:     false,
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
	for _, pm := range SYSTEM_PACKAGE_MANAGERS {
		if _, err := exec.LookPath(pm.Name); err == nil {
			return pm, nil
		}
	}
	return nil, fmt.Errorf("no supported package manager found")
}
