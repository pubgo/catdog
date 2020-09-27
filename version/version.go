package version

import (
	"github.com/coreos/etcd/version"
	ver "github.com/hashicorp/go-version"
	"github.com/pubgo/catdog/catdog_version"
	"github.com/pubgo/xerror"
	"runtime"
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

	xerror.ExitErr(ver.NewVersion(version.Version))
}

func Init() {
	catdog_version.Register("catdog_version", catdog_version.M{
		"build_time": BuildTime,
		"version":    Version,
		"go_version": GoVersion,
		"go_path":    GoPath,
		"go_root":    GoROOT,
		"commit_id":  CommitID,
		"project":    Project,
	})
}
