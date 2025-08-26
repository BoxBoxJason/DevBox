package install

import (
	"devbox/internal/commands"
	"devbox/pkg/utils"
	"sync"

	"go.uber.org/zap"
)

var (
	KUBERNETES_INSTALLABLE_TOOLCHAIN = &commands.InstallableToolchain{
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
		ExportedBinaries: []string{
			"kubectl",
			"kustomize",
			"helm",
			"yamllint",
			"dot",
			"k9s",
		},
		PackageManager: &utils.PackageManager{
			Name:         "krew",
			InstallCmd:   "install",
			SudoRequired: false,
			MultiInstall: true,
		},
		PackageManagerPackages: []string{
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
		ExtraSteps: &InstallKubernetesExternalDependenciesFunc,
	}

	KUBERNETES_PIP_PACKAGES = []string{
		"KubeDiagrams",
	}

	KUBERNETES_GO_PACKAGES = []string{
		"sigs.k8s.io/krew/cmd/krew@latest",
		"github.com/norwoodj/helm-docs/cmd/helm-docs@latest",
		"sigs.k8s.io/kind@latest",
	}

	InstallKubernetesExternalDependenciesFunc = InstallKubernetesExternalDependencies
)

func InstallKubernetesExternalDependencies(args *commands.SharedCmdArgs) []error {
	errChan := make(chan []error, 2)
	var wg sync.WaitGroup

	if len(KUBERNETES_GO_PACKAGES) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			zap.L().Warn("Installing krew using go install, which may fail if go is not properly set up.")
			zap.L().Warn("If krew is not already installed, the first run may fail to install krew plugins. Please rerun the command to ensure all packages are installed.")
			errChan <- GOLANG_INSTALLABLE_TOOLCHAIN.PackageManager.Install(KUBERNETES_GO_PACKAGES)
		}()
	}

	if len(KUBERNETES_PIP_PACKAGES) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			zap.L().Warn("Installing Python packages using pip, which may fail if Python is not properly set up.")
			errChan <- PYTHON_INSTALLABLE_TOOLCHAIN.PackageManager.Install(KUBERNETES_PIP_PACKAGES)
		}()
	}

	wg.Wait()
	close(errChan)
	return utils.MergeErrors(errChan)
}
