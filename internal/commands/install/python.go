package install

import (
	"devbox/internal/commands"
	"devbox/pkg/packagemanager"
)

var (
	PYTHON_INSTALLABLE_TOOLCHAIN = &commands.Toolchain{
		Name:        "python",
		Description: "Python development environment",
		InstalledPackages: []string{
			"python",
			"pip",
		},
		ExportedBinaries: []string{
			"python",
			"pip",
		},
		ExportedApplications: []string{},
		EnvironmentVariables: map[string]string{
			"PIP_INDEX_URL":             "${PIP_INDEX_URL:-https://pypi.org/simple}",
			"PIP_BREAK_SYSTEM_PACKAGES": "${PIP_BREAK_SYSTEM_PACKAGES:1}",
			"PIP_CACHE_DIR":             "${PIP_CACHE_DIR:-${XDG_CACHE_HOME}/pip}",
			"PYTHONUSERBASE":            "${PYTHONUSERBASE:-${XDG_DATA_HOME}/python}",
			"PATH":                      "${PYTHONUSERBASE}/bin:${PATH}",
		},
		PackageManagers: &map[*packagemanager.PackageManager][]string{
			packagemanager.PYTHON_PACKAGE_MANAGER: {
				"pylint",
				"black",
				"bandit",
				"pytest",
				"mypy",
				"flake8",
				"autopep8",
			},
		},
		VSCodeExtensions: []string{
			"ms-python.python",
			"ms-python.vscode-pylance",
			"ms-python.debugpy",
			"ms-python.vscode-python-envs",
			"njpwerner.autodocstring",
			"ms-python.pylint",
		},
		VSCodeSettings: map[string]any{
			"python.analysis.autoImportCompletions":                     true,
			"python.analysis.completeFunctionParens":                    true,
			"python.analysis.typeCheckingMode":                          "strict",
			"python.analysis.typeEvaluation.deprecateTypingAliases":     true,
			"python.analysis.typeEvaluation.enableReachabilityAnalysis": true,
			"python.analysis.typeEvaluation.strictDictionaryInference":  true,
			"python.analysis.typeEvaluation.strictListInference":        true,
			"python.analysis.typeEvaluation.strictSetInference":         true,
			"python.terminal.activateEnvInCurrentTerminal":              true,
			"python.testing.pytestEnabled":                              true,
			"python.testing.unittestEnabled":                            true,
			"python.useEnvironmentsExtension":                           true,
			"python-envs.terminal.showActivateButton":                   true,
		},
	}
)
