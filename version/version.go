package version

import (
	"runtime"

	"github.com/pubgo/xerror"

	ver "github.com/hashicorp/go-version"
	"github.com/pubgo/catdog/plugins/catdog_version"
)

var BuildTime = ""
var Version = ""
var GoVersion = runtime.Version()
var GoPath = ""
var GoROOT = ""
var CommitID = ""
var Project = ""

func init() {
	if Version == "" {
		Version = "v0.0.1"
	}

	xerror.ExitErr(ver.NewVersion(Version))
	xerror.Exit(catdog_version.Register("catdog_version", catdog_version.M{
		"build_time": BuildTime,
		"version":    Version,
		"go_version": GoVersion,
		"go_path":    GoPath,
		"go_root":    GoROOT,
		"commit_id":  CommitID,
		"project":    Project,
	}))
}
