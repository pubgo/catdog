package catdog_pidfile_plugin

import (
	"fmt"
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const pidPerm os.FileMode = 0755

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name    string
	pidPath string
}

func (p *Plugin) Commands() *cobra.Command {
	return nil
}

func (p *Plugin) Handler() *catdog_handler.Handler {
	return nil
}

func (p *Plugin) String() string {
	return p.name
}

func (p *Plugin) GetPid() (int, error) {
	f, err := p.GetPidFile()
	if err != nil {
		return 0, xerror.Wrap(err)
	}

	pid, err := ioutil.ReadFile(f)
	if err != nil {
		return 0, xerror.Wrap(err)
	}
	return strconv.Atoi(string(pid))
}

func (p *Plugin) GetPidFile() (string, error) {
	if !catdog_util.PathExist(p.pidPath) {
		if err := os.MkdirAll(p.pidPath, pidPerm); err != nil {
			return "", err
		}
	}

	fullPath := filepath.Join(p.pidPath, p.GetPidName())
	return fullPath, nil
}

func (p *Plugin) SavePid() error {
	f, err := p.GetPidFile()
	if err != nil {
		return xerror.Wrap(err)
	}

	pidBytes := []byte(strconv.Itoa(os.Getpid()))
	return xerror.Wrap(ioutil.WriteFile(f, pidBytes, pidPerm))
}

func (p *Plugin) Init(cat catdog_abc.CatDog) error {
	cat.Init(catdog_abc.AfterStart(p.SavePid))
	return xerror.Wrap(dix.Dix(p))
}

func (p *Plugin) GetPidName() string {
	name := catdog_config.Domain
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	return fmt.Sprintf("%s.%s.pid", catdog_util.GetBinName(), name)
}

func (p *Plugin) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(p.name, pflag.PanicOnError)
	flags.StringVar(&p.pidPath, "pidpath", p.pidPath, "pid path")
	return flags
}

func New() *Plugin {
	return &Plugin{
		name:    "pidfile",
		pidPath: filepath.Join(catdog_config.CfgDir, "pidfile"),
	}
}
