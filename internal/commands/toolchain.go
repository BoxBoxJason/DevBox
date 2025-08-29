package commands

import (
	"devbox/pkg/packagemanager"
	"devbox/pkg/utils"
	"devbox/pkg/vscode"
	"sync"
)

type SharedCmdArgs struct {
	SkipIde  bool
	NoExport bool
}

type Toolchain struct {
	Name                 string
	Description          string
	InstalledPackages    []string
	ExportedBinaries     []string
	ExportedApplications []string
	PackageManagers      *map[*packagemanager.PackageManager][]string
	VSCodeExtensions     []string
	VSCodeSettings       map[string]any
	EnvironmentVariables map[string]string
	PostInstallHooks     *func(args *SharedCmdArgs) []error
}

// Install installs the toolchain by performing the following steps:
// 1. Install the system packages specified by the toolchain.
// 2. Export the binaries and applications specified by the toolchain.
// 3. Install the IDE tools (like VSCode extensions) specified by the toolchain
// 4. Install the recommended development packages using the toolchain's package manager.
// 5. Run any extra installation steps specified by the toolchain.
// It uses goroutines to perform steps 1, 2, and 3 in parallel for better performance.
func (it *Toolchain) Install(args *SharedCmdArgs) []error {
	wgPackages := sync.WaitGroup{}
	wgOverall := sync.WaitGroup{}
	errChan := make(chan []error, 5)

	wgPackages.Add(1)
	wgOverall.Add(1)
	go func() {
		defer wgPackages.Done()
		errChan <- it.InstallSystemPackages(args)
	}()

	// Export the toolchain binaries to the user's environment
	wgOverall.Add(1)
	go func() {
		defer wgOverall.Done()
		wgPackages.Wait()
		errChan <- it.ExportSystemPackages(args)
	}()

	// Install the VSCode extensions for development
	wgOverall.Add(1)
	go func() {
		defer wgOverall.Done()
		errChan <- it.InstallIDETools(args)
	}()

	// Install the recommended development packages using go install
	if it.PackageManagers != nil && len(*it.PackageManagers) > 0 {
		for pkgManager, packages := range *it.PackageManagers {
			wgOverall.Add(1)
			go func(pm *packagemanager.PackageManager, pkgs []string) {
				defer wgOverall.Done()
				wgPackages.Wait()
				errChan <- pm.Install(pkgs)
			}(pkgManager, packages)
		}
	}

	// Run any extra installation steps for the toolchain
	if it.PostInstallHooks != nil {
		wgOverall.Wait()
		wgOverall.Add(1)
		go func() {
			defer wgOverall.Done()
			errChan <- (*it.PostInstallHooks)(args)
		}()
	}

	wgOverall.Wait()
	close(errChan)
	return utils.MergeErrors(errChan)
}

// InstallSystemPackages installs the system packages specified by the toolchain.
// It uses the system package manager to install the packages.
func (it *Toolchain) InstallSystemPackages(args *SharedCmdArgs) []error {
	if len(it.InstalledPackages) > 0 {
		return packagemanager.SystemPackageManager.Install(it.InstalledPackages)
	}
	return nil
}

// ExportSystemPackages exports the binaries and applications specified by the toolchain.
func (it *Toolchain) ExportSystemPackages(args *SharedCmdArgs) []error {
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

// InstallIDETools installs the IDE tools (like VSCode extensions) for the toolchain.
// It also creates or updates the IDE settings as specified by the toolchain.
func (it *Toolchain) InstallIDETools(args *SharedCmdArgs) []error {
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
			errChan <- vscode.VSCODE_PACKAGE_MANAGER.Install(it.VSCodeExtensions)
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
