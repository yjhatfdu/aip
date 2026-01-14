package version

// Set these at build time with -ldflags.
var (
	Version      = "dev"
	Commit       = "none"
	BuildTime    = "unknown"
	DefaultModel = "unknown"
)
