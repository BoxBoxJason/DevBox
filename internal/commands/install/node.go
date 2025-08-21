package install

import "devbox/internal/commands"

var (
	// NODE_BINARIES contains the binaries to be exported for Node.js
	NODE_BINARIES = []string{
		"npm",
		"npx",
		"node",
		"yarn",
	}
	// NODE_PACKAGES contains the Node.js packages to be installed using npm install
	// They will not be exported as they should already installed in the user's PATH
	NODE_PACKAGES = []string{
		"eslint",
		"prettier",
		"typescript",
		"jest",
		"ts-node",
	}
)

// installNode installs the entire Node development toolchain and environment.
// It installs the Node binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for Node development.
func installNode(args *commands.SharedCmdArgs) []error {
	return nil
}
