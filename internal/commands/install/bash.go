package install

import "devbox/internal/commands"

var (
	// BASH_INSTALLABLE_TOOLCHAIN is the installable toolchain for Bash
	BASH_INSTALLABLE_TOOLCHAIN = &commands.Toolchain{
		Name:        "bash",
		Description: "Bash development environment",
		InstalledPackages: []string{
			"bash",
			"shfmt",
			"shellcheck",
			"zsh",
			"fish",
		},
		ExportedBinaries: []string{
			"bash",
			"shfmt",
			"shellcheck",
			"zsh",
			"fish",
		},
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
