package install

import (
	"devbox/internal/commands"
	"devbox/internal/commands/setup"
	"devbox/pkg/utils"
	"fmt"

	"go.uber.org/zap"
)

func InstallToolchains(args *commands.SharedCmdArgs, toolchains ...string) []error {
	zap.L().Debug("Installing toolchains", zap.Strings("toolchains", toolchains))
	errorsChannel := make(chan []error, len(toolchains)+1)
	err := utils.LoadEnv(setup.DEFAULT_ENVIRONMENT)
	if err != nil {
		zap.L().Warn("Failed to load default environment variables", zap.Errors("errors", err))
		errorsChannel <- err
	}
	for _, toolchain := range toolchains {
		switch toolchain {
		case "bash":
			errorsChannel <- BASH_INSTALLABLE_TOOLCHAIN.Install(args)
		case "c":
			errorsChannel <- installC(args)
		case "container":
			errorsChannel <- CONTAINER_INSTALLABLE_TOOLCHAIN.Install(args)
		case "cpp":
			errorsChannel <- installCpp(args)
		case "github":
			errorsChannel <- GITHUB_INSTALLABLE_TOOLCHAIN.Install(args)
		case "gitlab":
			errorsChannel <- GITLAB_INSTALLABLE_TOOLCHAIN.Install(args)
		case "golang":
			errorsChannel <- GOLANG_INSTALLABLE_TOOLCHAIN.Install(args)
		case "java":
			errorsChannel <- installJava(args)
		case "kubernetes":
			errorsChannel <- installKubernetes(args)
		case "node":
			errorsChannel <- NODE_INSTALLABLE_TOOLCHAIN.Install(args)
		case "python":
			errorsChannel <- PYTHON_INSTALLABLE_TOOLCHAIN.Install(args)
		case "rust":
			errorsChannel <- installRust(args)
		default:
			errorsChannel <- []error{fmt.Errorf("unknown toolchain: %s", toolchain)}
		}
	}
	close(errorsChannel)
	return utils.MergeErrors(errorsChannel)
}
