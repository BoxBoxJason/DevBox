package install

import "devbox/internal/commands"

var (
	GITLAB_INSTALLABLE_TOOLCHAIN = &commands.Toolchain{
		Name:        "gitlab",
		Description: "GitLab development environment",
		InstalledPackages: []string{
			"yamllint",
			"glab",
			"git-lfs",
		},
		ExportedBinaries: []string{
			"yamllint",
			"glab",
			"git-lfs",
		},
		ExportedApplications: []string{},
		VSCodeExtensions: []string{
			"gitlab.gitlab-workflow",
			"redhat.vscode-yaml",
		},
		VSCodeSettings: map[string]any{
			"gitlab.showPipelineUpdateNotifications": true,
			"yaml.format.printWidth":                 100,
		},
	}
)
