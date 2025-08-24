package install

import (
	"devbox/internal/commands"
	"devbox/pkg/utils"
)

var (
	GOLANG_INSTALLABLE_TOOLCHAIN = &commands.InstallableToolchain{
		Name:        "golang",
		Description: "Golang development environment",
		ExportedBinaries: []string{
			"go",
			"make",
			"gofmt",
		},
		ExportedApplications: []string{},
		EnvironmentVariables: map[string]string{
			"GOMAXPROCS":  "${GOMAXPROCS:-$(nproc)}",
			"GOPATH":      "${GOPATH:-${XDG_DATA_HOME}/go}",
			"GOCACHE":     "${GOCACHE:-${XDG_CACHE_HOME}/go}",
			"GO111MODULE": "${GO111MODULE:-on}",
			"CGO_ENABLED": "${CGO_ENABLED:-0}",
			"GOFLAGS":     "${GOFLAGS:--trimpath -modcacherw}",
			"GOPROXY":     "${GOPROXY:-https://proxy.golang.org,direct}",
			"GOSUMDB":     "${GOSUMDB:-sum.golang.org}",
			"PATH":        "${GOPATH}/bin:${PATH}",
		},
		PackageManager: &utils.PackageManager{
			Name:         "go",
			InstallCmd:   "install",
			MultiInstall: false,
			SudoRequired: false,
		},
		PackageManagerPackages: []string{
			"github.com/securego/gosec/v2/cmd/gosec@latest",
			"github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest",
			"honnef.co/go/tools/cmd/staticcheck@latest",
			"github.com/axw/gocov/gocov@latest",
			"golang.org/x/tools/gopls@latest",
			"golang.org/x/tools/cmd/godoc@latest",
			"golang.org/x/tools/cmd/cover@latest",
			"github.com/go-delve/delve/cmd/dlv@latest",
			"github.com/cweill/gotests/gotests@latest",
			"github.com/fatih/gomodifytags@latest",
			"github.com/josharian/impl@latest",
		},
		VSCodeExtensions: []string{
			"golang.go",
		},
		VSCodeSettings: map[string]any{
			"go.coverMode":                  "atomic",
			"go.coverOnSingleTestFile":      true,
			"go.coverOnSingleTest":          true,
			"go.diagnostic.vulncheck":       "Imports",
			"go.lintTool":                   "golangci-lint",
			"go.toolsManagement.autoUpdate": true,
		},
	}
)
