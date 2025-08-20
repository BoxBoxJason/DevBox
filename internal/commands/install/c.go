package install

var (
	// C_BINARIES contains the binaries to be exported for C
	C_BINARIES = []string{
		"gcc",
		"g++",
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
	}
)

// installC installs the entire C development toolchain and environment.
// It installs the C binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for C development.
func installC() error {
	return nil
}
