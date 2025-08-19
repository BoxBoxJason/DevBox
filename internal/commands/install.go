package commands

var (
	// BASH_BINARIES contains the binaries to be exported for Bash
	BASH_BINARIES = []string{
		"bash",
		"shfmt",
		"shellcheck",
		"zsh",
		"fish",
	}

	// GOLANG_EXPORTED_BINARIES contains the binaries to be exported for Go
	GOLANG_EXPORTED_BINARIES = []string{
		"go",
	}
	// GOLANG_PACKAGES contains the Go packages to be installed using go install
	// They will not be exported as they should already installed in the user's PATH
	GOLANG_PACKAGES = []string{
		"github.com/securego/gosec/v2/cmd/gosec@latest",
		"github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest",
		"honnef.co/go/tools/staticcheck@latest",
		"github.com/axw/gocov/gocov@latest",
		"golang.org/x/tools/gopls",
		"golang.org/x/tools/cmd/godoc",
		"golang.org/x/tools/cmd/cover",
		"github.com/go-delve/delve/cmd/dlv@latest",
	}

	// RUST_BINARIES contains the binaries to be exported for Rust
	RUST_BINARIES = []string{
		"cargo",
		"rustc",
		"rustup",
	}
	// RUST_PACKAGES contains the Rust packages to be installed using cargo install
	// They will not be exported as they should already installed in the user's PATH
	RUST_PACKAGES = []string{
		"cargo-clippy",
		"cargo-auditable",
		"cargo-edit",
		"cargo-audit",
		"cargo-watch",
		"cargo-fmt",
	}

	// PYTHON_BINARIES contains the binaries to be exported for Python
	PYTHON_BINARIES = []string{
		"python",
		"pip",
	}
	// PYTHON_PACKAGES contains the Python packages to be installed using pip
	// They will not be exported as they should already installed in the user's PATH
	PYTHON_PACKAGES = []string{
		"pylint",
		"black",
		"bandit",
		"pytest",
		"mypy",
		"flake8",
		"autopep8",
	}

	// NODE_BINARIES contains the binaries to be exported for Node.js
	NODE_BINARIES = []string{
		"npm",
		"npx",
		"node",
		"yarn",
	}
	// NODE_PACKAGES contains the Node.js packages to be installed using npm install
	// They will not be exported as they should already installed in the user's PATH
	NODE_PACKAGES = []string{
		"eslint",
		"prettier",
		"typescript",
		"jest",
		"ts-node",
	}

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

	// CONTAINER_BINARIES contains the binaries to be exported for Container tools
	CONTAINER_BINARIES = []string{
		"podman",
		"hadolint",
		"trivy",
	}

	// JAVA_BINARIES contains the binaries to be exported for Java
	JAVA_BINARIES = []string{
		"java",
		"javac",
		"mvn",
		"junit",
		"jacoco",
	}

	// TODO JAVA_BINARIES_DOWNLOAD contains the list of binaries to be installed via download
	JAVA_BINARIES_DOWNLOAD = []string{
		"checkstyle",
		"spotbugs",
	}

	// C_BINARIES contains the binaries to be exported for C
	C_BINARIES = []string{
		"gcc",
		"g++",
		"clang",
		"make",
		"cmake",
		"clang-tidy",
		"cppcheck",
		"clang-format",
		"valgrind",
		"lcov",
		"gcovr",
		"gdb",
	}

	// CPP_BINARIES contains the binaries to be exported for C++
	CPP_BINARIES = []string{
		"g++",
		"clang++",
		"make",
		"cmake",
		"clang-tidy",
		"cppcheck",
		"clang-format",
		"valgrind",
		"lcov",
		"gcovr",
		"gdb",
	}

	// GITLAB_BINARIES contains the binaries to be exported for GitLab
	GITLAB_BINARIES = []string{
		"yamllint",
		"glab",
		"git-lfs",
	}

	// GITHUB_BINARIES contains the binaries to be exported for GitHub
	GITHUB_BINARIES = []string{
		"gh",
		"hub",
		"git-lfs",
	}
)
