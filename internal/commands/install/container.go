package install

import "devbox/internal/commands"

var (
	CONTAINER_INSTALLABLE_TOOLCHAIN = &commands.InstallableToolchain{
		Name:        "container",
		Description: "Container development environment",
		ExportedBinaries: []string{
			"podman",
			"hadolint",
			"trivy",
			"podman-compose",
		},
		ExportedApplications: []string{
			"podman-desktop",
		},
		EnvironmentVariables: map[string]string{},
		PackageManager:       nil,
		VSCodeExtensions: []string{
			"ms-azuretools.vscode-docker",
			"docker.docker",
			"exiasr.hadolint",
		},
		VSCodeSettings: map[string]any{
			"containers.composeCommand":                     "podman compose",
			"containers.containerClient":                    "com.microsoft.visualstudio.containers.podman",
			"containers.containerCommand":                   "podman",
			"containers.orchestratorClient":                 "com.microsoft.visualstudio.orchestrators.podmancompose",
			"docker.lsp.experimental.scout.notPinnedDigest": true,
			"docker.lsp.experimental.scout.recommendedTag":  true,
			"docker.lsp.telemetry":                          "off",
			"hadolint.outputLevel":                          "hint",
		},
	}
)
