package commands

import (
	"devbox/pkg/packagemanager"
	"devbox/pkg/utils"
	"devbox/pkg/vscode"
	"errors"
	"maps"
	"sync"
)

var (
	ErrNoToolchain = errors.New("no toolchains specified, use --help to see available toolchains")
)

// InstallToolchains installs the given toolchains by performing the following steps:
// 1. Install the system packages specified by the toolchains.
// 2. Export the binaries and applications specified by the toolchains.
// 3. Install the IDE tools (like VSCode extensions) specified by the toolchains.
// 4. Install the recommended development packages using the toolchains' package managers.
// 5. Run any extra installation steps specified by the toolchains.
// It uses goroutines to perform steps 1, 2, and 3 in parallel for better performance.
func InstallToolchains(args *SharedCmdArgs, toolchains ...*Toolchain) []error {
	if len(toolchains) == 0 {
		return []error{ErrNoToolchain}
	}
	errChan := make(chan []error, 4)
	wgPackages := sync.WaitGroup{}
	wgOverall := sync.WaitGroup{}

	wgPackages.Add(1)
	wgOverall.Add(1)
	go func() {
		defer wgPackages.Done()
		defer wgOverall.Done()
		errChan <- InstallToolchainsBinaries(toolchains...)
	}()

	wgOverall.Add(1)
	go func() {
		defer wgOverall.Done()
		wgPackages.Wait()
		errChan <- ExportToolchainsPackages(args, toolchains...)
	}()

	wgOverall.Add(1)
	go func() {
		defer wgOverall.Done()
		wgPackages.Wait()
		errChan <- InstallToolchainsPackages(toolchains...)
	}()

	wgOverall.Add(1)
	go func() {
		defer wgOverall.Done()
		errChan <- InstallToolchainsIDETools(args, toolchains...)
	}()

	wgOverall.Wait()
	close(errChan)
	return utils.MergeErrors(errChan)
}

// InstallToolchainsBinaries installs the system packages specified by the given toolchains.
// It uses the system package manager to install the packages.
func InstallToolchainsBinaries(toolchains ...*Toolchain) []error {
	if len(toolchains) == 0 {
		return []error{ErrNoToolchain}
	}

	toolChainsRawMergedBinaries := make([][]string, len(toolchains))
	for i, tc := range toolchains {
		toolChainsRawMergedBinaries[i] = tc.InstalledPackages
	}
	return packagemanager.SystemPackageManager.Install(utils.MergeStringSlices(toolChainsRawMergedBinaries...))
}

// ExportToolchainsPackages exports the binaries and applications specified by the given toolchains.
func ExportToolchainsPackages(args *SharedCmdArgs, toolchains ...*Toolchain) []error {
	if args.NoExport {
		return nil
	}

	if len(toolchains) == 0 {
		return []error{ErrNoToolchain}
	}

	mergedExportedBinaries := make([][]string, len(toolchains))
	mergedExportedApplications := make([][]string, len(toolchains))

	for i, tc := range toolchains {
		mergedExportedBinaries[i] = append(tc.ExportedBinaries, tc.ExportedApplications...)
		mergedExportedApplications[i] = tc.ExportedApplications
	}

	errChan := make(chan []error, 2)
	errChan <- utils.ExportDistroboxBinaries(utils.MergeStringSlices(mergedExportedBinaries...))
	errChan <- utils.ExportDistroboxApplications(utils.MergeStringSlices(mergedExportedApplications...))

	close(errChan)
	return utils.MergeErrors(errChan)
}

// InstallToolchainsIDETools installs the IDE tools (like VSCode extensions) for the given toolchains.
// It also creates or updates the IDE settings as specified by the toolchains.
func InstallToolchainsIDETools(args *SharedCmdArgs, toolchains ...*Toolchain) []error {
	if args.SkipIde {
		return nil
	}
	if len(toolchains) == 0 {
		return []error{ErrNoToolchain}
	}

	idePluginsRawMerged := make([][]string, len(toolchains))
	ideSettingsMerged := make(map[string]any)
	for i, tc := range toolchains {
		idePluginsRawMerged[i] = tc.VSCodeExtensions
		maps.Copy(ideSettingsMerged, tc.VSCodeSettings)
	}
	var wg sync.WaitGroup
	errChan := make(chan []error, 2)

	unDuplicatedPlugins := utils.MergeStringSlices(idePluginsRawMerged...)
	if len(unDuplicatedPlugins) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- vscode.VSCODE_PACKAGE_MANAGER.Install(unDuplicatedPlugins)
		}()
	}

	if len(ideSettingsMerged) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errChan <- []error{vscode.SystemVSCode.UpdateSettings(ideSettingsMerged)}
		}()
	}

	wg.Wait()
	close(errChan)
	return utils.MergeErrors(errChan)
}

// InstallToolchainsPackages installs the recommended development packages using the package managers specified by the given toolchains.
// Every installation is done in parallel using goroutines.
func InstallToolchainsPackages(toolchains ...*Toolchain) []error {
	if len(toolchains) == 0 {
		return []error{ErrNoToolchain}
	}

	// Group toolchains by their package manager to avoid duplicate installations
	packageManagerToPackages := make(map[*packagemanager.PackageManager][]string)
	for _, tc := range toolchains {
		if tc.PackageManagers != nil {
			for pkgManager, packages := range *tc.PackageManagers {
				existingPackages, exists := packageManagerToPackages[pkgManager]
				if exists {
					packageManagerToPackages[pkgManager] = utils.MergeStringSlices(existingPackages, packages)
				} else {
					packageManagerToPackages[pkgManager] = packages
				}
			}
		}
	}

	errChan := make(chan []error, len(packageManagerToPackages))
	var wg sync.WaitGroup

	for pkgManager, packages := range packageManagerToPackages {
		wg.Add(1)
		go func(pm *packagemanager.PackageManager, pkgs []string) {
			defer wg.Done()
			errChan <- pm.Install(pkgs)
		}(pkgManager, packages)
	}

	wg.Wait()
	close(errChan)
	return utils.MergeErrors(errChan)
}
