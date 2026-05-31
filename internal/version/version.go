// Package version provides build-time version information for Spectra.
package version

import (
	"fmt"
	"runtime"
)

// Version is set at build time via -ldflags.
var Version = "dev"

// Commit is set at build time via -ldflags.
var Commit = "unknown"

// Date is set at build time via -ldflags.
var Date = "unknown"

// Info returns a formatted version string suitable for display.
func Info() string {
	return fmt.Sprintf("Spectra %s (commit: %s, built: %s, go: %s)",
		Version, Commit, Date, runtime.Version())
}
