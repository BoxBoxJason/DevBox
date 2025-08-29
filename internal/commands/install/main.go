package install

import (
	"devbox/internal/commands"
	"fmt"

	"go.uber.org/zap"
)

var (
	// EXISTING_TOOLCHAINS is the list of all installable toolchains
	EXISTING_TOOLCHAINS = map[string]*commands.Toolchain{
		"bash":       BASH_INSTALLABLE_TOOLCHAIN,
		"c":          C_INSTALLABLE_TOOLCHAIN,
		"container":  CONTAINER_INSTALLABLE_TOOLCHAIN,
		"cpp":        CPP_INSTALLABLE_TOOLCHAIN,
		"github":     GITHUB_INSTALLABLE_TOOLCHAIN,
		"gitlab":     GITLAB_INSTALLABLE_TOOLCHAIN,
		"golang":     GOLANG_INSTALLABLE_TOOLCHAIN,
		"java":       JAVA_INSTALLABLE_TOOLCHAIN,
		"kubernetes": KUBERNETES_INSTALLABLE_TOOLCHAIN,
		"node":       NODE_INSTALLABLE_TOOLCHAIN,
		"python":     PYTHON_INSTALLABLE_TOOLCHAIN,
		"rust":       RUST_INSTALLABLE_TOOLCHAIN,
	}
)

func InstallToolchains(args *commands.SharedCmdArgs, toolchains ...string) []error {
	installableToolchains, err := parseToolchains(toolchains)
	if err != nil {
		return []error{err}
	}

	toolChainsNames := make([]string, len(installableToolchains))
	for i, tc := range installableToolchains {
		toolChainsNames[i] = tc.Name
	}
	zap.L().Info("Installing toolchains", zap.Strings("toolchains", toolChainsNames))
	return commands.InstallToolchains(args, installableToolchains...)
}

func parseToolchains(toolchains []string) ([]*commands.Toolchain, error) {
	if len(toolchains) == 0 {
		return nil, fmt.Errorf("no toolchains specified, use --help to see available toolchains")
	}
	var parsedToolchains []*commands.Toolchain
	var toolchainNames = make(map[string]struct{})
	for _, toolchain := range toolchains {
		if _, exists := toolchainNames[toolchain]; exists {
			continue
		}
		installableToolchain, exists := EXISTING_TOOLCHAINS[toolchain]
		if !exists {
			return nil, fmt.Errorf("unknown toolchain: %s", toolchain)
		}
		parsedToolchains = append(parsedToolchains, installableToolchain)
		toolchainNames[toolchain] = struct{}{}
	}
	return parsedToolchains, nil
}
