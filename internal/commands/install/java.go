package install

import "devbox/internal/commands"

var (
	// JAVA_INSTALLABLE_TOOLCHAIN defines the Java toolchain with its packages, binaries, and settings.
	JAVA_INSTALLABLE_TOOLCHAIN = &commands.Toolchain{
		Name:        "java",
		Description: "Java development environment",
		InstalledPackages: []string{
			"jacoco",
			"java-25-openjdk",
			"junit",
			"maven",
		},
		ExportedBinaries: []string{
			"jar",
			"jarsigner",
			"javac",
			"java",
			"java2html",
			"javadoc",
			"javap",
			"jcmd",
			"jconsole",
			"jdb",
			"jdeprscan",
			"jdeps",
			"jfr",
			"jhsdb",
			"jimage",
			"jlink",
			"jmod",
			"jpackage",
			"jps",
			"jrunscript",
			"jshell",
			"jstat",
			"jstatd",
			"keytool",
			"mvn",
			"jacococli",
		},
	}

	// JAVA_BINARIES_DOWNLOAD contains the list of binaries to be installed via download
	JAVA_BINARIES_DOWNLOAD = []string{
		"checkstyle",
		"spotbugs",
	}
)
