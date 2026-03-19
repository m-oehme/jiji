// Package version holds build-time version information set via ldflags.
package version

// Set via -ldflags at build time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
