package install

import (
	"devbox/internal/commands"
	"devbox/pkg/utils"
)

var (
	NODE_INSTALLABLE_TOOLCHAIN = &commands.InstallableToolchain{
		Name:        "node",
		Description: "Node.js development environment",
		InstalledPackages: []string{
			"npm",
			"npx",
			"node",
			"yarn",
		},
		ExportedBinaries: []string{
			"npm",
			"npx",
			"node",
			"yarn",
		},
		ExportedApplications: []string{},
		EnvironmentVariables: map[string]string{
			"npm_config_prefix":   "${XDG_DATA_HOME}/npm",
			"NPM_CONFIG_REGISTRY": "https://registry.npmjs.org/",
			"NPM_CONFIG_CACHE":    "${XDG_CACHE_HOME}/npm",
			"PATH":                "${npm_config_prefix}/bin:${YARN_GLOBAL_FOLDER}/bin:${PATH}",
			"YARN_CACHE_FOLDER":   "${XDG_CACHE_HOME}/yarn",
			"YARN_GLOBAL_FOLDER":  "${XDG_DATA_HOME}/yarn",
			"YARN_CONFIG_FOLDER":  "${XDG_CONFIG_HOME}/yarn",
			"YARN_REGISTRY":       "https://registry.yarnpkg.com",
		},
		PackageManager: &utils.PackageManager{
			Name:             "npm",
			InstallCmd:       "install",
			NoInteractiveArg: utils.StrPtr("--user"),
		},
		PackageManagerPackages: []string{
			"eslint",
			"prettier",
			"typescript",
			"jest",
			"ts-node",
			"esbuild",
		},
		VSCodeExtensions: []string{
			"ms-vscode.vscode-typescript-next",
			"dbaeumer.vscode-eslint",
			"esbenp.prettier-vscode",
			"orta.vscode-jest",
			"christian-kohler.npm-intellisense",
			"rvest.vs-code-prettier-eslint",
		},
		VSCodeSettings: map[string]any{
			"eslint.format.enable":            true,
			"eslint.lintTask.enable":          true,
			"eslint.run":                      "onSave",
			"eslint.useESLintClass":           true,
			"prettier.bracketSameLine":        true,
			"prettier.experimentalTernaries":  true,
			"prettier.singleAttributePerLine": true,
		},
	}
)
