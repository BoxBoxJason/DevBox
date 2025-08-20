package install

import "fmt"

func InstallToolchains(toolchains ...string) error {
	for _, toolchain := range toolchains {
		switch toolchain {
		case "bash":
			if err := installBash(); err != nil {
				return fmt.Errorf("failed to install Bash toolchain: %w", err)
			}
		case "c":
			if err := installC(); err != nil {
				return fmt.Errorf("failed to install C toolchain: %w", err)
			}
		case "container":
			if err := installContainer(); err != nil {
				return fmt.Errorf("failed to install Container toolchain: %w", err)
			}
		case "cpp":
			if err := installCpp(); err != nil {
				return fmt.Errorf("failed to install C++ toolchain: %w", err)
			}
		case "github":
			if err := installGitHub(); err != nil {
				return fmt.Errorf("failed to install GitHub toolchain: %w", err)
			}
		case "gitlab":
			if err := installGitLab(); err != nil {
				return fmt.Errorf("failed to install GitLab toolchain: %w", err)
			}
		case "golang":
			if err := installGolang(); err != nil {
				return fmt.Errorf("failed to install Go toolchain: %w", err)
			}
		case "java":
			if err := installJava(); err != nil {
				return fmt.Errorf("failed to install Java toolchain: %w", err)
			}
		case "kubernetes":
			if err := installKubernetes(); err != nil {
				return fmt.Errorf("failed to install Kubernetes toolchain: %w", err)
			}
		case "node":
			if err := installNode(); err != nil {
				return fmt.Errorf("failed to install Node toolchain: %w", err)
			}
		case "python":
			if err := installPython(); err != nil {
				return fmt.Errorf("failed to install Python toolchain: %w", err)
			}
		case "rust":
			if err := installRust(); err != nil {
				return fmt.Errorf("failed to install Rust toolchain: %w", err)
			}

		default:
			return fmt.Errorf("unknown toolchain: %s", toolchain)
		}
	}
	return nil

}
