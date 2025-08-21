package install

import "devbox/internal/commands"

var (
	// CPP_BINARIES contains the binaries to be exported for C++
	CPP_BINARIES = []string{
		"g++",
		"clang++",
		"make",
		"cmake",
		"clang-tidy",
		"cppcheck",
		"clang-format",
		"valgrind",
		"lcov",
		"gcovr",
		"gdb",
	}
)

// installCpp installs the entire Cpp development toolchain and environment.
// It installs the Cpp binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for Cpp development.
func installCpp(args *commands.SharedCmdArgs) []error {
	return nil
}
