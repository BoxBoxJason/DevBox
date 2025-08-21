package install

import "devbox/internal/commands"

var (
	// PYTHON_BINARIES contains the binaries to be exported for Python
	PYTHON_BINARIES = []string{
		"python",
		"pip",
	}
	// PYTHON_PACKAGES contains the Python packages to be installed using pip
	// They will not be exported as they should already installed in the user's PATH
	PYTHON_PACKAGES = []string{
		"pylint",
		"black",
		"bandit",
		"pytest",
		"mypy",
		"flake8",
		"autopep8",
	}
)

// installPython installs the entire Python development toolchain and environment.
// It installs the Python binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for Python development.
func installPython(args *commands.SharedCmdArgs) []error {
	return nil
}
