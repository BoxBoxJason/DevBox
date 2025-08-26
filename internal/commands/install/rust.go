package install

import (
	"devbox/internal/commands"
	"devbox/pkg/utils"
)

var (
	// RUST_INSTALLABLE_TOOLCHAIN is the installable toolchain for Rust
	RUST_INSTALLABLE_TOOLCHAIN = &commands.InstallableToolchain{
		Name:        "rust",
		Description: "Rust development environment",
		InstalledPackages: []string{
			"cargo",
			"rustc",
			"rustup",
		},
		ExportedBinaries: []string{
			"cargo",
			"rustc",
			"rustup",
		},
		EnvironmentVariables: map[string]string{
			"CARGO_HOME":  "${XDG_DATA_HOME}/cargo",
			"RUSTUP_HOME": "${XDG_DATA_HOME}/rustup",
			"PATH":        "${CARGO_HOME}/bin:${PATH}",
		},
		PackageManager: &utils.PackageManager{
			Name:             "cargo",
			InstallCmd:       "install",
			NoInteractiveArg: nil,
			SudoRequired:     false,
			MultiInstall:     true,
		},
		PackageManagerPackages: []string{
			"cargo-clippy",
			"cargo-auditable",
			"cargo-edit",
			"cargo-audit",
			"cargo-watch",
			"cargo-fmt",
		},
		VSCodeExtensions: []string{
			"rust-lang.rust-analyzer",
			"tamasfe.even-better-toml",
		},
		VSCodeSettings: map[string]any{
			"evenBetterToml.formatter.alignComments":                 true,
			"evenBetterToml.formatter.arrayAutoExpand":               true,
			"evenBetterToml.formatter.arrayAutoCollapse":             true,
			"evenBetterToml.formatter.allowedBlankLines":             1,
			"evenBetterToml.formatter.alignEntries":                  true,
			"evenBetterToml.formatter.arrayTrailingComma":            true,
			"evenBetterToml.formatter.columnWidth":                   100,
			"evenBetterToml.formatter.indentEntries":                 true,
			"evenBetterToml.formatter.indentString":                  " ",
			"evenBetterToml.formatter.indentTables":                  true,
			"evenBetterToml.formatter.reorderArrays":                 true,
			"evenBetterToml.formatter.reorderInlineTables":           true,
			"evenBetterToml.formatter.trailingNewline":               true,
			"rust-analyzer.cargo.features":                           "all",
			"rust-analyzer.completion.fullFunctionSignatures.enable": true,
		},
	}
)
