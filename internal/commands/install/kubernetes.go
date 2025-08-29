package install

import (
	"devbox/internal/commands"
	"devbox/pkg/packagemanager"
)

var (
	KUBERNETES_INSTALLABLE_TOOLCHAIN = &commands.Toolchain{
		Name:        "kubernetes",
		Description: "Kubernetes toolchain including kubectl, kustomize, helm, yamllint, and more.",
		InstalledPackages: []string{
			"kubectl",
			"kustomize",
			"helm",
			"yamllint",
			"dot",
			"k9s",
		},
		EnvironmentVariables: map[string]string{
			"KREW_ROOT":                  "${KREW_ROOT:-{XDG_DATA_HOME}/krew}",
			"KIND_EXPERIMENTAL_PROVIDER": "podman",
		},
		ExportedBinaries: []string{
			"kubectl",
			"kustomize",
			"helm",
			"yamllint",
			"dot",
			"k9s",
		},
		PackageManagers: &map[*packagemanager.PackageManager][]string{
			packagemanager.GOLANG_PACKAGE_MANAGER: {
				"sigs.k8s.io/krew/cmd/krew@latest",
				"github.com/norwoodj/helm-docs/cmd/helm-docs@latest",
				"sigs.k8s.io/kind@latest",
			},
			packagemanager.PYTHON_PACKAGE_MANAGER: {
				"KubeDiagrams",
			},
			packagemanager.KREW_PACKAGE_MANAGER: {
				"ai",
				"blame",
				"cost",
				"debug-shell",
				"deprecations",
				"explore",
				"flame",
				"kor",
				"neat",
				"tree",
			},
		},
	}
)
