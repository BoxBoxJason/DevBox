package packagemanager

import "devbox/pkg/utils"

var (
	GOLANG_PACKAGE_MANAGER = &PackageManager{
		Name:         "go",
		InstallCmd:   "install",
		MultiInstall: false,
		SudoRequired: false,
	}

	KREW_PACKAGE_MANAGER = &PackageManager{
		Name:         "krew",
		InstallCmd:   "install",
		SudoRequired: false,
		MultiInstall: true,
	}

	PYTHON_PACKAGE_MANAGER = &PackageManager{
		Name:         "pip",
		InstallCmd:   "install",
		MultiInstall: true,
		SudoRequired: false,
	}

	NODE_PACKAGE_MANAGER = &PackageManager{
		Name:             "npm",
		InstallCmd:       "install",
		MultiInstall:     true,
		NoInteractiveArg: utils.StrPtr("--user"),
	}

	CARGO_PACKAGE_MANAGER = &PackageManager{
		Name:             "cargo",
		InstallCmd:       "install",
		NoInteractiveArg: nil,
		SudoRequired:     false,
		MultiInstall:     true,
	}
)
