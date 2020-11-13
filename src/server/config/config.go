package serverConfig

import (
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/vari"
)

type Config struct {
	ServerIP string
	ServerPort int

	DBDriver string
	DBPath string
}

func NewConfig() *Config {
	return &Config{
		ServerIP: vari.Ip,
		ServerPort: vari.Port,
		DBDriver: constant.SqliteDriver,
		DBPath: constant.SqliteData,
	}
}
