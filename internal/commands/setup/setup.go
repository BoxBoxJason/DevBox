package setup

var (
	// DEFAULT_DEV_BINARIES contains the default development binaries to be exported
	DEFAULT_DEV_BINARIES = []string{
		"git",
		"podman",
		"markdownlint",
		"tree",
		"curl",
		"wget",
		"vim",
		"jq",
		"yq",
		"bat",
	}

	// DEFAULT_DEV_APPS contains the default development applications to be exported
	DEFAULT_DEV_APPS = []string{
		"code",
	}

	// DEFAULT_VSCODE_EXTENSIONS contains the default VSCode extensions to be installed
	DEFAULT_VSCODE_EXTENSIONS = []string{
		"avidanson.vscode-markdownlint",
		"bierner.markdown-mermaid",
	}

	// DEFAULT_ENVIRONMENT contains the default environment variables for development
	DEFAULT_ENVIRONMENT = map[string]string{
		"XDG_CONFIG_HOME": "${XDG_CONFIG_HOME:-$HOME/.config}",
		"XDG_DATA_HOME":   "${XDG_DATA_HOME:-$HOME/.local/share}",
		"XDG_CACHE_HOME":  "${XDG_CACHE_HOME:-$HOME/.cache}",
		"XDG_STATE_HOME":  "${XDG_STATE_HOME:-$HOME/.local/state}",
		"XDG_RUNTIME_DIR": "${XDG_RUNTIME_DIR:-/run/user/$(id -u)}",
		"XDG_CONFIG_DIRS": "${XDG_CONFIG_DIRS:-/etc/xdg}",
		"XDG_DATA_DIRS":   "${XDG_DATA_DIRS:-/usr/local/share:/usr/share:$XDG_DATA_HOME}",
		"XDG_STATE_DIRS":  "${XDG_STATE_DIRS:-/var/lib/xdg}",
		"XDG_CACHE_DIRS":  "${XDG_CACHE_DIRS:-/var/cache/xdg:$XDG_CACHE_HOME}",
		"LANG":            "en_US.UTF-8",
		"CLICOLOR":        "1",
		"EDITOR":          "vim",
		"ARCHFLAGS":       "-arch $(uname -m)",
	}
)
