package catdog_pidfile

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/xerror"
)

const pidPerm os.FileMode = 0755

func GetPid() (pid int, _ error) {
	pidData, err := ioutil.ReadFile(GetPidPath())
	if err != nil {
		return 0, xerror.Wrap(err)
	}

	pid, err = strconv.Atoi(string(pidData))
	return pid, xerror.Wrap(err)
}

func SavePid() error {
	pidBytes := []byte(strconv.Itoa(os.Getpid()))
	return xerror.Wrap(ioutil.WriteFile(GetPidPath(), pidBytes, pidPerm))
}

func GetPidPath() string {
	return filepath.Join(catdog_config.Home, "pidfile", catdog_config.Domain+"."+catdog_config.Project+".pid")
}

func Init() {}
