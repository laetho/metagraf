package params

var(
	// CreateSecrets Boolean flag to indicate if we should create empty secrets.
	CreateSecrets bool = false

	// LocalBuild Boolean flag to indicate we want to build from local directory.
	LocalBuild bool = false

	// FileBuild Give path to a file that should be copied to the container image being built (i.e. my.war).
	FileBuild string = ""
)