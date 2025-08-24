package setup

import (
	"devbox/internal/commands"
	"devbox/pkg/utils"
	"devbox/pkg/vscode"
	"sync"
)

var (
	// DEFAULT_IDE contains the default IDE to be installed
	DEFAULT_IDE = "code"

	// DEFAULT_DEV_BINARIES contains the default development binaries to be exported
	DEFAULT_DEV_BINARIES = []string{
		"git",
		"podman",
		"markdownlint",
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
		"avidanson.vscode-markdownlint",
		"bierner.markdown-mermaid",
	}

	// DEFAULT_VSCODE_SETTINGS contains the default VSCode settings to be applied
	DEFAULT_VSCODE_SETTINGS = map[string]any{
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
		"ARCHFLAGS":       "-arch $(uname -m)",
	}
)

func SetupDevbox(args *commands.SharedCmdArgs) []error {
	// Set the development environment variables
	errs := utils.SystemEnvManager.Set(DEFAULT_ENVIRONMENT)
	if errs != nil {
		return errs
	}

	// Use a WaitGroup to manage parallel installations
	var wg sync.WaitGroup
	errChan := make(chan []error, 3) // Channel to collect errors from goroutines

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Install generic utility software development / unix binaries
		errChan <- utils.SystemPackageManager.Install(DEFAULT_DEV_BINARIES)
	}()

	// Install the VSCode extensions for Go development
	if !args.SkipIde {
		errs := utils.SystemPackageManager.Install([]string{DEFAULT_IDE})
		if errs != nil {
			return errs
		}

		wg.Add(2)
		go func() {
			defer wg.Done()
			errChan <- utils.VSCODE_PACKAGE_MANAGER.Install(DEFAULT_VSCODE_EXTENSIONS)
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

		if !args.SkipIde {
			wg.Add(2)
			go func() {
				defer wg.Done()
				errChan <- utils.ExportDistroboxApplications([]string{DEFAULT_IDE})
			}()
			go func() {
				defer wg.Done()
				errChan <- utils.ExportDistroboxBinaries([]string{DEFAULT_IDE})
			}()
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	return utils.MergeErrors(errChan)
}
