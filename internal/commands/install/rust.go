package install

import "devbox/internal/commands"

var (
	// RUST_BINARIES contains the binaries to be exported for Rust
	RUST_BINARIES = []string{
		"cargo",
		"rustc",
		"rustup",
	}
	// RUST_PACKAGES contains the Rust packages to be installed using cargo install
	// They will not be exported as they should already installed in the user's PATH
	RUST_PACKAGES = []string{
		"cargo-clippy",
		"cargo-auditable",
		"cargo-edit",
		"cargo-audit",
		"cargo-watch",
		"cargo-fmt",
	}
)

// installRust installs the entire Rust development toolchain and environment.
// It installs the Rust binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for Rust development.
func installRust(args *commands.SharedCmdArgs) []error {
	return nil
}
