package install

import (
	"devbox/pkg/utils"
	"fmt"
)

var (
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

	// GOLANG_ENVIRONMENT contains the environment variables to be set for Golang development
	GOLANG_ENVIRONMENT = map[string]string{
		"GOMAXPROCS":  "${GOMAXPROCS:-$(nproc)}",
		"GOPATH":      "${GOPATH:-${XDG_DATA_HOME}/go}",
		"GOCACHE":     "${GOCACHE:-${XDG_CACHE_HOME}/go}",
		"GO111MODULE": "${GO111MODULE:-on}",
		"CGO_ENABLED": "${CGO_ENABLED:-0}",
		"GOFLAGS":     "${GOFLAGS:--trimpath -modcacherw}",
		"GOPROXY":     "${GOPROXY:-https://proxy.golang.org,direct}",
		"GOSUMDB":     "${GOSUMDB:-sum.golang.org}",
	}

	GoPackageManager = &utils.PackageManager{
		Name:         "go",
		InstallCmd:   "install",
		MultiInstall: false,
	}
)

// installGolang installs the entire Golang development toolchain and environment.
// It installs the Go binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for Go development.
func installGolang() error {
	// Set the Go development environment variables
	err := utils.SystemEnvManager.Set(GOLANG_ENVIRONMENT)
	if err != nil {
		return fmt.Errorf("failed to set Go development environment variables: %w", err)
	}

	// Install the Go toolchain binaries
	err = utils.SystemPackageManager.Install(GOLANG_EXPORTED_BINARIES)
	if err != nil {
		panic(fmt.Sprintf("Error during golang toolchain installation: %v", err))
	}

	// Export the Go toolchain binaries to the user's environment
	err = utils.ExportDistroboxBinaries(GOLANG_EXPORTED_BINARIES)
	if err != nil {
		panic(fmt.Sprintf("Error during golang binary export: %v", err))
	}

	// Install the recommended development packages using go install
	err = GoPackageManager.Install(GOLANG_PACKAGES)
	if err != nil {
		panic(fmt.Sprintf("Error during golang package installation: %v", err))
	}

	return nil
}
