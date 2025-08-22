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
			if err := installBash(args); err != nil {
				errorsChannel <- err
			}
		case "c":
			errorsChannel <- installC(args)
		case "container":
			errorsChannel <- installContainer(args)
		case "cpp":
			errorsChannel <- installCpp(args)
		case "github":
			errorsChannel <- installGitHub(args)
		case "gitlab":
			errorsChannel <- installGitLab(args)
		case "golang":
			errorsChannel <- installGolang(args)
		case "java":
			errorsChannel <- installJava(args)
		case "kubernetes":
			errorsChannel <- installKubernetes(args)
		case "node":
			errorsChannel <- installNode(args)
		case "python":
			errorsChannel <- installPython(args)
		case "rust":
			errorsChannel <- installRust(args)
		default:
			errorsChannel <- []error{fmt.Errorf("unknown toolchain: %s", toolchain)}
		}
	}
	close(errorsChannel)
	return utils.MergeErrors(errorsChannel)
}
