package commands

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
	}

	// DEFAULT_DEV_APPS contains the default development applications to be exported
	DEFAULT_DEV_APPS = []string{
		"code",
	}
)
