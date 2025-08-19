package main

import "github.com/spf13/cobra"

/* This file is part of the devbox go project
 * It defines the CLI commands for the devbox
 */

var (
	version = "dev"

	mainCmd = &cobra.Command{
		Use:     "devbox",
		Version: version,
		Short:   "devbox is the package manager for the distrobox ecosystem",
		Long: `devbox is the package manager for the distrobox ecosystem.
It helps you install packages on your distrobox and export them to your host system.`,
	}

	// TODO: Add --skip-ide flag to skip IDE installation
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup the devbox by installing the minimum required packages",
		Long: `Setup the devbox by installing the minimum required packages.
This command will install the necessary packages to get started with devbox.
It installs the minimal required packages to start developing with devbox.`,
	}

	exportCmd = &cobra.Command{
		Use:   "install",
		Short: "Install a language toolchain or a package",
		Long: `Install a language toolchain or a package.
Supports installing language toolchains for Bash, Go, Rust, Python, Node, Kubernetes, Container, Java, GitLab, GitHub, C, C++`,
	}

	sharePackageCmd = &cobra.Command{
		Use:   "share",
		Short: "Share a package with the host system",
		Long: `Share a package with the host system.
This command will install the package in the distrobox if it is not already installed and export it to the host system.`,
	}
)

func main() {
	mainCmd.AddCommand(setupCmd, exportCmd, sharePackageCmd)
	if err := mainCmd.Execute(); err != nil {
		panic(err)
	}
}
