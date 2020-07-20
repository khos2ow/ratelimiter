package version

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// current version
const dev = "v0.2.0"

// Provisioned by ldflags
var (
	version    string
	commitHash string
	buildDate  string
)

func init() {
	// Load defaults for info variables
	if version == "" {
		version = dev
	}
	if commitHash == "" {
		commitHash = dev
	}
	if buildDate == "" {
		buildDate = time.Now().Format(time.RFC3339)
	}
}

// Full return the full version of the binary including commit hash and build date
func Full() string {
	if !strings.HasSuffix(version, commitHash) {
		version += " " + commitHash
	}
	return fmt.Sprintf("Version   : %s\nBuild Date: %s", version, buildDate)
}

// String return the full version of the binary including commit hash and build date
func String() string {
	if !strings.HasSuffix(version, commitHash) {
		version += " " + commitHash
	}
	return fmt.Sprintf("%s %s/%s BuildDate: %s", version, runtime.GOOS, runtime.GOARCH, buildDate)
}
