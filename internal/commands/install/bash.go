package install

import "devbox/internal/commands"

var (
	// BASH_INSTALLABLE_TOOLCHAIN is the installable toolchain for Bash
	BASH_INSTALLABLE_TOOLCHAIN = &commands.InstallableToolchain{
		Name:        "bash",
		Description: "Bash development environment",
		ExportedBinaries: []string{
			"bash",
			"shfmt",
			"shellcheck",
			"zsh",
			"fish",
		},
		PackageManager: nil,
		VSCodeExtensions: []string{
			"timonwong.shellcheck",
			"foxundermoon.shell-format",
		},
		VSCodeSettings: map[string]any{
			"shellcheck.run":                   "onSave",
			"shellcheck.useWorkspaceRootAsCwd": true,
		},
	}
)
