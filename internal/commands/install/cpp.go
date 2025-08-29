package install

import "devbox/internal/commands"

var CPP_INSTALLABLE_TOOLCHAIN = &commands.Toolchain{
	// C++ toolchain
	Name:        "cpp",
	Description: "C++ toolchain including gcc, clang, make, cmake, gdb, and common linters/formatters.",
	InstalledPackages: []string{
		"gcc",
		"g++",
		"clang",
		"clang++",
		"make",
		"cmake",
		"ninja",
		"clang-tidy",
		"cppcheck",
		"clang-format",
		"valgrind",
		"lcov",
		"gcovr",
		"gdb",
	},
	ExportedBinaries: []string{
		"gcc",
		"g++",
		"clang",
		"clang++",
		"make",
		"cmake",
		"ninja",
		"clang-tidy",
		"cppcheck",
		"clang-format",
		"valgrind",
		"lcov",
		"gcovr",
		"gdb",
	},
	VSCodeExtensions: []string{
		"ms-vscode.cpptools",
		"ms-vscode.cpptools-extension-pack",
		"ms-vscode.cmake-tools",
		"ms-vscode.cpptools-themes",
		"ms-vscode.makefile-tools",
	},
	VSCodeSettings: map[string]any{
		"makefile.configureOnOpen":             true,
		"cmake.deleteBuildDirOnCleanConfigure": true,
		"cmake.generator":                      "Ninja",
	},
	EnvironmentVariables: map[string]string{
		"CC":        "$(which gcc)",
		"CXX":       "$(which g++)",
		"MAKEFLAGS": "-j$(nproc)",
	},
}
