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
			if err := installC(args); err != nil {
				errorsChannel <- err
			}
		case "container":
			if err := installContainer(args); err != nil {
				errorsChannel <- err
			}
		case "cpp":
			if err := installCpp(args); err != nil {
				errorsChannel <- err
			}
		case "github":
			if err := installGitHub(args); err != nil {
				errorsChannel <- err
			}
		case "gitlab":
			if err := installGitLab(args); err != nil {
				errorsChannel <- err
			}
		case "golang":
			if err := installGolang(args); err != nil {
				errorsChannel <- err
			}
		case "java":
			if err := installJava(args); err != nil {
				errorsChannel <- err
			}
		case "kubernetes":
			if err := installKubernetes(args); err != nil {
				errorsChannel <- err
			}
		case "node":
			if err := installNode(args); err != nil {
				errorsChannel <- err
			}
		case "python":
			if err := installPython(args); err != nil {
				errorsChannel <- err
			}
		case "rust":
			if err := installRust(args); err != nil {
				errorsChannel <- err
			}

		default:
			errorsChannel <- []error{fmt.Errorf("unknown toolchain: %s", toolchain)}
		}
	}
	close(errorsChannel)
	return utils.MergeErrors(errorsChannel)

}
