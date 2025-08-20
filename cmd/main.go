package main

import (
	"devbox/internal/commands/install"
	"devbox/pkg/utils"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

/* This file is part of the devbox go project
 * It defines the CLI commands for the devbox
 */

var (
	version = "dev"

	installCmdFilePath string

	setupSkipIde bool

	mainCmd = &cobra.Command{
		Use:     "devbox",
		Version: version,
		Short:   "devbox is the package manager for the distrobox ecosystem",
		Long: `devbox is the package manager for the distrobox ecosystem.
It helps you install packages on your distrobox and export them to your host system.`,
	}

	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Setup the devbox by installing the minimum required packages",
		Long: `Setup the devbox by installing the minimum required packages.
This command will install the necessary packages to get started with devbox.
It installs the minimal required packages to start developing with devbox.`,
	}

	installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install a language toolchain or a package",
		Long: `Install a language toolchain or a package.
Supports installing language toolchains for Bash, Go, Rust, Python, Node, Kubernetes, Container, Java, GitLab, GitHub, C, C++`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Aggregate arguments
			var allArgs []string
			allArgs = append(allArgs, args...)

			installCmdFilePath = strings.TrimSpace(installCmdFilePath)

			// If --file flag is provided, read the file and append its content line by line
			if installCmdFilePath != "" {
				fileArgs, err := utils.ReadFileLines(installCmdFilePath)
				if err != nil {
					return fmt.Errorf("failed to process file: %w", err)
				}
				allArgs = append(allArgs, fileArgs...)
			}

			// Call the install function with all arguments
			return install.InstallToolchains(allArgs...)
		},
	}

	sharePackageCmd = &cobra.Command{
		Use:   "share",
		Short: "Share a package with the host system",
		Long: `Share a package with the host system.
This command will install the package in the distrobox if it is not already installed and export it to the host system.`,
	}
)

func main() {
	installCmd.Flags().StringVar(&installCmdFilePath, "file", "", "Path to a file containing a list of languages toolchains to install, one per line")

	setupCmd.Flags().BoolVar(&setupSkipIde, "skip-ide", false, "Skip IDE installation")

	mainCmd.AddCommand(setupCmd, installCmd, sharePackageCmd)
	if err := mainCmd.Execute(); err != nil {
		panic(err)
	}
}
