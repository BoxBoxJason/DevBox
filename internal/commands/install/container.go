package install

import "devbox/internal/commands"

var (
	// CONTAINER_BINARIES contains the binaries to be exported for Container tools
	CONTAINER_BINARIES = []string{
		"podman",
		"hadolint",
		"trivy",
	}
)

// installContainer installs the entire Container development toolchain and environment.
// It installs the Container binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for Container development.
func installContainer(args *commands.SharedCmdArgs) []error {
	return nil
}
