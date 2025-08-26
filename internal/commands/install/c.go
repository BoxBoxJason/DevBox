package install

import "devbox/internal/commands"

var (
	// C_INSTALLABLE_TOOLCHAIN is the installable toolchain for C
	C_INSTALLABLE_TOOLCHAIN = &commands.InstallableToolchain{
		Name:        "c",
		Description: "C toolchain including gcc, clang, make, cmake, gdb, and more.",
		InstalledPackages: []string{
			"gcc",
			"clang",
			"make",
			"cmake",
			"clang-tidy",
			"cppcheck",
			"clang-format",
			"valgrind",
			"lcov",
			"gcovr",
			"gdb",
			"cmake-gui",
		},
		ExportedBinaries: []string{
			"gcc",
			"clang",
			"make",
			"cmake",
			"clang-tidy",
			"cppcheck",
			"clang-format",
			"valgrind",
			"lcov",
			"gcovr",
			"gdb",
			"ctest",
			"cpack",
			"cmake-gui",
		},
		VSCodeExtensions: []string{
			"ms-vscode.makefile-tools",
			"ms-vscode.cpptools",
			"ms-vscode.cpptools-extension-pack",
			"ms-vscode.cmake-tools",
			"ms-vscode.cpptools-themes",
		},
		VSCodeSettings: map[string]any{
			"makefile.configureOnOpen":             true,
			"cmake.deleteBuildDirOnCleanConfigure": true,
		},
	}
)
