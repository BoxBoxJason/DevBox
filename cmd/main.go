package main

import (
	"devbox/internal/commands/install"
	"devbox/internal/commands/setup"
	"devbox/pkg/utils"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	version = "dev"

	args ParserArgs

	mainCmd = &cobra.Command{
		Use:     "devbox {setup|install|share} [flags]",
		Version: version,
		Short:   "devbox is the package manager for the distrobox ecosystem",
		Long: `devbox is the package manager for the distrobox ecosystem.
It helps you install packages on your distrobox and export them to your host system.`,
		PersistentPreRun: func(cmd *cobra.Command, preRunArgs []string) {
			// Setup the logger with the verbosity level and file output if specified
			SetupZapLogger(args.Verbose, args.LogFilePath)
			zap.L().Debug("Verbose mode enabled")
		},
	}

	setupCmd = &cobra.Command{
		Use:   "setup [--skip-ide] [--verbose] [--log-file <PATH>] [--no-export]",
		Short: "Setup the devbox by installing the minimum required packages",
		Long: `Setup the devbox by installing the minimum required packages.
This command will install the necessary packages to get started with devbox.
It installs the minimal required packages to start developing with devbox.`,
		Run: func(cmd *cobra.Command, commandArgs []string) {
			errs := setup.SetupDevbox(&args.SharedCmdArgs)
			if errs != nil {
				zap.L().Fatal("Failed to setup devbox", zap.Errors("errors", errs))
			}
		},
	}

	installCmd = &cobra.Command{
		Use:   "install [--skip-ide] [--no-export] [--file <PATH>] [toolchain...]",
		Short: "Install a language toolchain or a package",
		Long: `Install a language toolchain or a package.
Supports installing language toolchains for Bash, Go, Rust, Python, Node, Kubernetes, Container, Java, GitLab, GitHub, C, C++`,
		Run: func(cmd *cobra.Command, commandArgs []string) {
			// Aggregate arguments
			var allArgs []string
			allArgs = append(allArgs, commandArgs...)

			args.InstallCmdFilePath = strings.TrimSpace(args.InstallCmdFilePath)

			// If --file flag is provided, read the file and append its content line by line
			if args.InstallCmdFilePath != "" {
				zap.L().Debug("Reading install packages file", zap.String("file", args.InstallCmdFilePath))
				fileArgs, err := utils.ReadFileLines(args.InstallCmdFilePath)
				if err != nil {
					zap.L().Fatal("Failed to read install packages file", zap.String("file", args.InstallCmdFilePath), zap.Error(err))
				}
				allArgs = append(allArgs, fileArgs...)
			}

			if len(allArgs) == 0 {
				zap.L().Fatal("No toolchains or packages specified for installation. Use --file to specify a file or provide arguments directly.")
			}

			// Call the install function with all arguments
			err := install.InstallToolchains(&args.SharedCmdArgs, allArgs...)
			if err != nil {
				zap.L().Fatal("Failed to install toolchains", zap.Errors("errors", err))
			}
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
	installCmd.Flags().StringVar(&args.InstallCmdFilePath, "file", "", "Path to a file containing a list of languages toolchains to install, one per line")

	mainCmd.PersistentFlags().BoolVarP(&args.SkipIde, "skip-ide", "n", false, "Skip IDE installation")
	mainCmd.PersistentFlags().BoolVarP(&args.Verbose, "verbose", "v", false, "Enable verbose output")
	mainCmd.PersistentFlags().StringVarP(&args.LogFilePath, "log-file", "l", "", "Path to the log file")
	mainCmd.PersistentFlags().BoolVar(&args.NoExport, "no-export", false, "Do not export the package to the host system")

	mainCmd.AddCommand(setupCmd, installCmd, sharePackageCmd)
	if err := mainCmd.Execute(); err != nil {
		zap.L().Fatal("devbox runtime error", zap.Error(err))
	}
}
