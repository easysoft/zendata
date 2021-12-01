package serverConfig

import (
	"github.com/easysoft/zendata/src/utils/vari"
)

type Config struct {
	ServerIP   string
	ServerPort int
}

func NewConfig() *Config {
	return &Config{
		ServerIP:   vari.Ip,
		ServerPort: vari.Port,
	}
}
