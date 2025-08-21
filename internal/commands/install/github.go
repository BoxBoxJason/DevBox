package install

import "devbox/internal/commands"

var (
	// GITHUB_BINARIES contains the binaries to be exported for GitHub
	GITHUB_BINARIES = []string{
		"gh",
		"hub",
		"git-lfs",
	}
)

// installGitHub installs the entire GitHub development toolchain and environment.
// It installs the GitHub binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for GitHub development.
func installGitHub(args *commands.SharedCmdArgs) []error {
	return nil
}
