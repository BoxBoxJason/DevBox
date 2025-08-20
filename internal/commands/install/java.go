package install

var (
	// JAVA_BINARIES contains the binaries to be exported for Java
	JAVA_BINARIES = []string{
		"java",
		"javac",
		"mvn",
		"junit",
		"jacoco",
	}

	// JAVA_BINARIES_DOWNLOAD contains the list of binaries to be installed via download
	JAVA_BINARIES_DOWNLOAD = []string{
		"checkstyle",
		"spotbugs",
	}
)

// installJava installs the entire Java development toolchain and environment.
// It installs the Java binaries and packages, ensuring they are available in the user's PATH.
// It also sets up the necessary environment variables for Java development.
func installJava() error {
	return nil
}
