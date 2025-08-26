package install

import "devbox/internal/commands"

var (
	// GITHUB_INSTALLABLE_TOOLCHAIN is the installable toolchain for GitHub
	GITHUB_INSTALLABLE_TOOLCHAIN = &commands.InstallableToolchain{
		Name:        "github",
		Description: "GitHub development environment",
		InstalledPackages: []string{
			"gh",
			"hub",
			"git-lfs",
		},
		ExportedBinaries: []string{
			"gh",
			"hub",
			"git-lfs",
		},
		ExportedApplications: []string{},
		PackageManager:       nil,
		VSCodeExtensions: []string{
			"github.vscode-pull-request-github",
			"github.copilot",
			"github.vscode-github-actions",
			"redhat.vscode-yaml",
		},
		VSCodeSettings: map[string]any{
			"yaml.format.printWidth":                          100,
			"github-actions.workflows.pinned.refresh.enabled": true,
			"githubPullRequests.defaultMergeMethod":           "rebase",
			"githubPullRequests.notifications":                "pullRequests",
		},
	}
)
