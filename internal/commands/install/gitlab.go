package install

import "devbox/internal/commands"

var (
	GITLAB_INSTALLABLE_TOOLCHAIN = &commands.InstallableToolchain{
		Name:        "gitlab",
		Description: "GitLab development environment",
		ExportedBinaries: []string{
			"yamllint",
			"glab",
			"git-lfs",
		},
		ExportedApplications: []string{},
		PackageManager:       nil,
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
