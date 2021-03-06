package catdog_config

import (
	"path/filepath"

	"github.com/asim/nitro/v3/config"
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/xerror"
	"github.com/spf13/pflag"
)

// 默认的全局配置
var (
	Domain  = "catdog"
	Project = "catdog"
	Debug   = true
	Trace   = false
	Mode    = "dev"
	IsBlock = true
	Home    = filepath.Join(xerror.PanicStr(filepath.Abs(filepath.Dir(""))), "home")
	CfgPath = filepath.Join(Home, "config", "config.yaml")
	cfg     *Config
)

// RunMode 项目运行模式
var RunMode = struct {
	Dev     string
	Test    string
	Stag    string
	Prod    string
	Release string
}{
	Dev:     "dev",
	Test:    "test",
	Stag:    "stag",
	Prod:    "prod",
	Release: "release",
}

type Config struct {
	config.Config
}

func DefaultFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("app", pflag.PanicOnError)
	flags.StringVarP(&Mode, "mode", "m", Mode, "running mode(dev|test|stag|prod|release)")
	flags.StringVarP(&Home, "home", "c", Home, "project config home dir")
	flags.BoolVarP(&Debug, "debug", "d", Debug, "enable debug")
	flags.BoolVarP(&Trace, "trace", "t", Trace, "enable trace")
	flags.BoolVarP(&IsBlock, "block", "b", IsBlock, "enable signal block")
	return flags
}

func Load(path ...string) (reader.Value, error) {
	return cfg.Load(path...)
}

func LoadBytes() []byte {
	return xerror.PanicErr(cfg.Load()).(reader.Value).Bytes()
}

func GetCfg() *Config {
	return cfg
}

const tag = `
          _____                    _____                _____                    _____                   _______                   _____          
         /\    \                  /\    \              /\    \                  /\    \                 /::\    \                 /\    \         
        /::\    \                /::\    \            /::\    \                /::\    \               /::::\    \               /::\    \        
       /::::\    \              /::::\    \           \:::\    \              /::::\    \             /::::::\    \             /::::\    \       
      /::::::\    \            /::::::\    \           \:::\    \            /::::::\    \           /::::::::\    \           /::::::\    \      
     /:::/\:::\    \          /:::/\:::\    \           \:::\    \          /:::/\:::\    \         /:::/~~\:::\    \         /:::/\:::\    \     
    /:::/  \:::\    \        /:::/__\:::\    \           \:::\    \        /:::/  \:::\    \       /:::/    \:::\    \       /:::/  \:::\    \    
   /:::/    \:::\    \      /::::\   \:::\    \          /::::\    \      /:::/    \:::\    \     /:::/    / \:::\    \     /:::/    \:::\    \   
  /:::/    / \:::\    \    /::::::\   \:::\    \        /::::::\    \    /:::/    / \:::\    \   /:::/____/   \:::\____\   /:::/    / \:::\    \  
 /:::/    /   \:::\    \  /:::/\:::\   \:::\    \      /:::/\:::\    \  /:::/    /   \:::\ ___\ |:::|    |     |:::|    | /:::/    /   \:::\ ___\ 
/:::/____/     \:::\____\/:::/  \:::\   \:::\____\    /:::/  \:::\____\/:::/____/     \:::|    ||:::|____|     |:::|    |/:::/____/  ___\:::|    |
\:::\    \      \::/    /\::/    \:::\  /:::/    /   /:::/    \::/    /\:::\    \     /:::|____| \:::\    \   /:::/    / \:::\    \ /\  /:::|____|
 \:::\    \      \/____/  \/____/ \:::\/:::/    /   /:::/    / \/____/  \:::\    \   /:::/    /   \:::\    \ /:::/    /   \:::\    /::\ \::/    / 
  \:::\    \                       \::::::/    /   /:::/    /            \:::\    \ /:::/    /     \:::\    /:::/    /     \:::\   \:::\ \/____/  
   \:::\    \                       \::::/    /   /:::/    /              \:::\    /:::/    /       \:::\__/:::/    /       \:::\   \:::\____\    
    \:::\    \                      /:::/    /    \::/    /                \:::\  /:::/    /         \::::::::/    /         \:::\  /:::/    /    
     \:::\    \                    /:::/    /      \/____/                  \:::\/:::/    /           \::::::/    /           \:::\/:::/    /     
      \:::\    \                  /:::/    /                                 \::::::/    /             \::::/    /             \::::::/    /      
       \:::\____\                /:::/    /                                   \::::/    /               \::/____/               \::::/    /       
        \::/    /                \::/    /                                     \::/____/                 ~~                      \::/____/        
         \/____/                  \/____/                                       ~~                                                            
`
