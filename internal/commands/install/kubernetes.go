package install

var (
	// KUBERNETES_BINARIES contains the binaries to be exported for Kubernetes
	KUBERNETES_BINARIES = []string{
		"kubectl",
		"kustomize",
		"helm",
		"kics",
		"kubectx",
		"yamllint",
		"minikube",
	}
	// KUBERNETES_PACKAGES contains the Kubernetes packages to be installed using kubectl
	// They will not be exported as they should already installed in the user's PATH
	KUBERNETES_PACKAGES = []string{
		"kube-diagrams",
		"PyYAML",
	}
)

// installKubernetes installs the entire Kubernetes development toolchain and environment.
// It installs the Kubernetes binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for Kubernetes development.
func installKubernetes() error {
	return nil
}
