# DevBox

DevBox is a tool designed to simplify the management of development tools and environments. It provides a unified interface for installing, updating, and managing various programming language toolchains and binaries on your system.

It is a perfect tools for developers working with **immutable OS** like NixOS, where you want to keep your system clean and avoid installing tools globally. It relies on the concept of **distrobox**, which allows you to run containers with your host's environment, making it easy to manage dependencies and tools without cluttering your system.

> [!CAUTION]
> DevBox is currently in early development. It is only being tested using distrobox which support `dnf` as the package manager.
> Support for other package managers is not garanteed.

## Features

- **Install Toolchains**: Easily install and manage toolchains for various programming languages like Go, Rust, Python, Node.js, and more.
- **Export Packages**: Export installed packages to your host system, making them available globally.
- **Setup Environment**: Quickly set up your development environment with the necessary tools, dependencies and environment variables.
- **Cross-Platform**: Works on Linux, macOS, and Windows (via WSL).

## Installation

### Prerequisites

Before installing DevBox, ensure you have the following prerequisites

- **Distrobox**: Ensure you have [distrobox](https://distrobox.it/) installed on your system. Distrobox allows you to run containers with your host's environment, making it easy to manage dependencies and tools without cluttering your system.
- **Oh-My-Zsh**: DevBox is designed to work with [Oh-My-Zsh](https://ohmyz.sh/). Ensure you have it installed and configured on your system.

### From Binary

1. You can download the latest release of DevBox from the [releases page](https://github.com/BoxBoxJason/DevBox/releases).
2. Then you can either place the binary in your `PATH` and run it inside of an existing distrobox, or build your own image with the binary included.

### From Source

1. Clone the repository: `git clone https://github.com/boxboxjason/devbox.git`
2. Change into the directory: `cd devbox`
3. Build the binary: `make build`
4. You can then either place the binary in your `PATH` and run it inside of an existing distrobox, or build your own image with the binary included.

### From Docker

You can also run DevBox using the published Docker image. This is useful if you don't plan on building your own image.

1. Run the distrobox using the DevBox image: `distrobox create --image ghcr.io/boxboxjason/devbox:latest --name devbox`
2. Enter the distrobox: `distrobox enter devbox`
3. You can now run DevBox commands inside the distrobox.

## Usage

Once you have DevBox installed and are inside a distrobox, you can start using it to manage your development tools and environments.

### devbox setup

The `devbox setup` command initializes your development environment by setting up necessary environment variables and installing essential tools. It is recommended to run this command after entering a new distrobox.

*This command is idempotent, meaning you can run it multiple times without causing issues.*

> [!WARNING]
> This will install and export Visual Studio Code as an IDE, which may not be suitable for all users. You can use the `--skip-ide` flag to skip this step if you prefer not to use Visual Studio Code. You will then need to manually install your preferred IDE.

```bash
devbox setup
```

### devbox install

The `devbox install` command allows you to install various toolchains and binaries. You can specify the toolchains you want to install, and DevBox will handle the installation process for you. It will also export the installed packages to your host system, making them available globally.

*This command is also idempotent, meaning you can run it multiple times without causing issues.*

```bash
devbox install <toolchain1> <toolchain2> ...
```
