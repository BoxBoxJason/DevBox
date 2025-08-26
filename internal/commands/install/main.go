package install

import (
	"devbox/internal/commands"
	"devbox/pkg/utils"
	"fmt"

	"go.uber.org/zap"
)

func InstallToolchains(args *commands.SharedCmdArgs, toolchains ...string) []error {
	zap.L().Debug("Installing toolchains", zap.Strings("toolchains", toolchains))
	errorsChannel := make(chan []error, len(toolchains))
	for _, toolchain := range toolchains {
		switch toolchain {
		case "bash":
			errorsChannel <- BASH_INSTALLABLE_TOOLCHAIN.Install(args)
		case "c":
			errorsChannel <- C_INSTALLABLE_TOOLCHAIN.Install(args)
		case "container":
			errorsChannel <- CONTAINER_INSTALLABLE_TOOLCHAIN.Install(args)
		case "cpp":
			errorsChannel <- CPP_INSTALLABLE_TOOLCHAIN.Install(args)
		case "github":
			errorsChannel <- GITHUB_INSTALLABLE_TOOLCHAIN.Install(args)
		case "gitlab":
			errorsChannel <- GITLAB_INSTALLABLE_TOOLCHAIN.Install(args)
		case "golang":
			errorsChannel <- GOLANG_INSTALLABLE_TOOLCHAIN.Install(args)
		case "java":
			errorsChannel <- JAVA_INSTALLABLE_TOOLCHAIN.Install(args)
		case "kubernetes":
			errorsChannel <- KUBERNETES_INSTALLABLE_TOOLCHAIN.Install(args)
		case "node":
			errorsChannel <- NODE_INSTALLABLE_TOOLCHAIN.Install(args)
		case "python":
			errorsChannel <- PYTHON_INSTALLABLE_TOOLCHAIN.Install(args)
		case "rust":
			errorsChannel <- RUST_INSTALLABLE_TOOLCHAIN.Install(args)
		default:
			errorsChannel <- []error{fmt.Errorf("unknown toolchain: %s", toolchain)}
		}
	}
	close(errorsChannel)
	return utils.MergeErrors(errorsChannel)
}
