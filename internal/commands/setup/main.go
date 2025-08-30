package setup

import (
	"devbox/internal/commands"
	"devbox/internal/commands/install"
	"devbox/internal/envmanager"
	"devbox/pkg/packagemanager"
	"devbox/pkg/utils"
	"devbox/pkg/vscode"
	"sync"

	"go.uber.org/zap"
)

var (
	// DEFAULT_IDE contains the default IDE to be installed
	DEFAULT_IDE = "code"

	// DEFAULT_DEV_BINARIES contains the default development binaries to be exported
	DEFAULT_DEV_BINARIES = []string{
		"git",
		"tree",
		"curl",
		"wget",
		"vim",
		"jq",
		"yq",
		"bat",
	}

	// DEFAULT_DEV_APPS contains the default development applications to be exported
	DEFAULT_DEV_APPS = []string{}

	// DEFAULT_VSCODE_EXTENSIONS contains the default VSCode extensions to be installed
	DEFAULT_VSCODE_EXTENSIONS = []string{
		"davidanson.vscode-markdownlint",
		"bierner.markdown-mermaid",
		"fill-labs.dependi",
	}

	// DEFAULT_VSCODE_SETTINGS contains the default VSCode settings to be applied
	DEFAULT_VSCODE_SETTINGS = map[string]any{
		"dependi.rust.informPatchUpdates":                                   true,
		"dependi.npm.informPatchUpdates":                                    true,
		"dependi.npm.indexServerURL":                                        "https://registry.npmjs.org",
		"dependi.rust.indexServerURL":                                       "https://index.crates.io",
		"dependi.go.indexServerURL":                                         "https://proxy.golang.org",
		"dependi.go.informPatchUpdates":                                     true,
		"dependi.python.indexServerURL":                                     "https://pypi.org/pypi",
		"dependi.python.informPatchUpdates":                                 true,
		"dependi.php.indexServerURL":                                        "https://repo.packagist.org",
		"dependi.php.informPatchUpdates":                                    true,
		"dependi.dart.indexServerURL":                                       "https://pub.dev",
		"dependi.dart.informPatchUpdates":                                   true,
		"dependi.vulnerability.ghsa.enabled":                                true,
		"dependi.vulnerability.osvQueryURL.batch":                           "https://api.osv.dev/v1/querybatch",
		"dependi.vulnerability.osvQueryURL.single":                          "https://api.osv.dev/v1/query",
		"diffEditor.ignoreTrimWhitespace":                                   true,
		"diffEditor.experimental.showMoves":                                 true,
		"diffEditor.experimental.useTrueInlineView":                         true,
		"editor.acceptSuggestionOnEnter":                                    "smart",
		"editor.autoIndentOnPaste":                                          true,
		"editor.formatOnPaste":                                              true,
		"editor.formatOnSave":                                               true,
		"editor.bracketPairColorization.independentColorPoolPerBracketType": true,
		"editor.guides.bracketPairs":                                        "active",
		"editor.selectionHighlightMultiline":                                true,
		"editor.tabSize":                                                    4,
		"editor.trimWhitespaceOnDelete":                                     true,
		"explorer.incrementalNaming":                                        "smart",
		"files.autoGuessEncoding":                                           true,
		"files.autoSaveWhenNoErrors":                                        true,
		"files.insertFinalNewline":                                          true,
		"files.readonlyFromPermissions":                                     true,
		"files.trimFinalNewlines":                                           true,
		"files.trimTrailingWhitespace":                                      true,
		"git.autofetch":                                                     true,
		"git.confirmSync":                                                   false,
		"git.enableSmartCommit":                                             true,
		"telemetry.telemetryLevel":                                          "off",
		"telemetry.feedback.enabled":                                        false,
		"testing.coverageToolbarEnabled":                                    true,
		"workbench.commandPalette.experimental.suggestCommands":             true,
	}

	// DEFAULT_ENVIRONMENT contains the default environment variables for development
	DEFAULT_ENVIRONMENT = map[string]string{
		"XDG_CONFIG_HOME": "${XDG_CONFIG_HOME:-$HOME/.config}",
		"XDG_DATA_HOME":   "${XDG_DATA_HOME:-$HOME/.local/share}",
		"XDG_CACHE_HOME":  "${XDG_CACHE_HOME:-$HOME/.cache}",
		"XDG_STATE_HOME":  "${XDG_STATE_HOME:-$HOME/.local/state}",
		"XDG_RUNTIME_DIR": "${XDG_RUNTIME_DIR:-/run/user/$(id -u)}",
		"XDG_CONFIG_DIRS": "${XDG_CONFIG_DIRS:-/etc/xdg}",
		"XDG_DATA_DIRS":   "${XDG_DATA_DIRS:-/usr/local/share:/usr/share:$XDG_DATA_HOME}",
		"XDG_STATE_DIRS":  "${XDG_STATE_DIRS:-/var/lib/xdg}",
		"XDG_CACHE_DIRS":  "${XDG_CACHE_DIRS:-/var/cache/xdg:$XDG_CACHE_HOME}",
		"LANG":            "${LANG:-en_US.UTF-8}",
		"CLICOLOR":        "${CLICOLOR:-1}",
		"EDITOR":          "${EDITOR:-vim}",
		"OS":              "$(uname | tr '[:upper:]' '[:lower:]')",
		"ARCH":            "$(uname -m | sed -e 's/x86_64/amd64/' -e 's/(arm)(64)?.*/12/' -e 's/aarch64$/arm64/')",
		"ARCHFLAGS":       "-arch ${ARCH}",
	}
)

func SetupDevbox(args *commands.SharedCmdArgs) []error {
	// Set the development environment variables
	errs := envmanager.SystemEnvManager(envmanager.DEFAULT_SYS_ENV_FILE).Set(
		DEFAULT_ENVIRONMENT,
		install.BASH_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.C_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.CONTAINER_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.CPP_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.GITHUB_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.GITLAB_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.GOLANG_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.JAVA_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.KUBERNETES_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.NODE_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.PYTHON_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
		install.RUST_INSTALLABLE_TOOLCHAIN.EnvironmentVariables,
	)
	if errs != nil {
		return errs
	}

	// Use a WaitGroup to manage parallel installations
	var wg sync.WaitGroup

	// compute channel capacity to avoid blocking sends (which would prevent wg.Done() from running)
	maxSends := 1 // initial binaries goroutine
	if !args.SkipIde {
		maxSends += 2 // vscode extension install + settings update
	}
	if !args.NoExport {
		maxSends += 2 // export binaries + export apps
	}
	errChan := make(chan []error, maxSends) // Channel to collect errors from goroutines

	if !args.SkipIde {
		DEFAULT_DEV_BINARIES = append(DEFAULT_DEV_BINARIES, DEFAULT_IDE)
		DEFAULT_DEV_APPS = append(DEFAULT_DEV_APPS, DEFAULT_IDE)
	}

	// Install generic utility software development / unix binaries
	errChan <- packagemanager.SystemPackageManager.Install(DEFAULT_DEV_BINARIES)

	// Install the VSCode extensions for Go development
	if !args.SkipIde {
		wg.Add(2)
		go func() {
			defer wg.Done()
			errChan <- vscode.VSCODE_PACKAGE_MANAGER.Install(DEFAULT_VSCODE_EXTENSIONS)
		}()

		go func() {
			defer wg.Done()
			errChan <- []error{vscode.SystemVSCode.UpdateSettings(DEFAULT_VSCODE_SETTINGS)}
		}()
	}

	// Export the generic development binaries to the user's environment
	if !args.NoExport {
		wg.Add(2)
		go func() {
			defer wg.Done()
			errChan <- utils.ExportDistroboxBinaries(DEFAULT_DEV_BINARIES)
		}()

		// Export the generic development applications to the user's environment
		go func() {
			defer wg.Done()
			errChan <- utils.ExportDistroboxApplications(DEFAULT_DEV_APPS)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	errs = utils.MergeErrors(errChan)
	if len(errs) == 0 {
		zap.L().Info("DevBox setup complete! Please make sure to restart your shell or run 'source ~/.zshrc' to apply the changes.")
	}
	return errs
}
