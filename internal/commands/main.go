package commands

import (
	"devbox/internal/envmanager"
	"devbox/pkg/utils"
	"devbox/pkg/vscode"
	"sync"
)

type SharedCmdArgs struct {
	SkipIde  bool
	NoExport bool
}

type InstallableToolchain struct {
	Name                   string
	Description            string
	InstalledPackages      []string
	ExportedBinaries       []string
	ExportedApplications   []string
	PackageManager         *utils.PackageManager
	PackageManagerPackages []string
	VSCodeExtensions       []string
	VSCodeSettings         map[string]any
	EnvironmentVariables   map[string]string
	ExtraSteps             *func(args *SharedCmdArgs) []error
}

func (it *InstallableToolchain) Install(args *SharedCmdArgs) []error {
	if len(it.EnvironmentVariables) > 0 {
		// Set the Go development environment variables in the system env file
		errs := envmanager.SystemEnvManager(envmanager.DEFAULT_SYS_ENV_FILE).Set(it.EnvironmentVariables)
		if errs != nil {
			return errs
		}
	}

	// Install the system packages for the toolchain
	errs := it.InstallSystemPackages(args)
	if errs != nil {
		return errs
	}

	var wg sync.WaitGroup
	errChan := make(chan []error, 4)

	// Export the toolchain binaries to the user's environment
	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- it.ExportSystemPackages(args)
	}()

	// Install the VSCode extensions for development
	wg.Add(1)
	go func() {
		defer wg.Done()
		errChan <- it.InstallIDETools(args)
	}()

	// Install the recommended development packages using go install
	if it.PackageManager != nil && len(it.PackageManagerPackages) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- it.PackageManager.Install(it.PackageManagerPackages)
		}()
	}

	// Run any extra installation steps for the toolchain
	if it.ExtraSteps != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- (*it.ExtraSteps)(args)
		}()
	}

	wg.Wait()
	close(errChan)

	return utils.MergeErrors(errChan)
}

func (it *InstallableToolchain) InstallSystemPackages(args *SharedCmdArgs) []error {
	if len(it.InstalledPackages) > 0 {
		return utils.SystemPackageManager.Install(it.InstalledPackages)
	}
	return nil
}

func (it *InstallableToolchain) ExportSystemPackages(args *SharedCmdArgs) []error {
	mergedSystemPackages := append(it.ExportedBinaries, it.ExportedApplications...)
	if !args.NoExport && len(mergedSystemPackages) > 0 {
		var wg sync.WaitGroup
		errChan := make(chan []error, 2)

		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- utils.ExportDistroboxBinaries(mergedSystemPackages)
		}()

		if len(it.ExportedApplications) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				errChan <- utils.ExportDistroboxApplications(it.ExportedApplications)
			}()
		}

		wg.Wait()
		close(errChan)

		return utils.MergeErrors(errChan)
	}
	return nil
}

func (it *InstallableToolchain) InstallIDETools(args *SharedCmdArgs) []error {
	if args.SkipIde {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan []error, 2)

	// Install the VSCode extensions for development
	if len(it.VSCodeExtensions) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- utils.VSCODE_PACKAGE_MANAGER.Install(it.VSCodeExtensions)
		}()
	}

	if len(it.VSCodeSettings) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- []error{vscode.SystemVSCode.UpdateSettings(it.VSCodeSettings)}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	return utils.MergeErrors(errChan)
}
