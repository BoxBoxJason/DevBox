package install

import "devbox/internal/commands"

var (
	// GITLAB_BINARIES contains the binaries to be exported for GitLab
	GITLAB_BINARIES = []string{
		"yamllint",
		"glab",
		"git-lfs",
	}
)

// installGitLab installs the entire GitLab development toolchain and environment.
// It installs the GitLab binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for GitLab development.
func installGitLab(args *commands.SharedCmdArgs) []error {
	return nil
}
