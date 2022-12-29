package serverConfig

import (
	"github.com/easysoft/zendata/pkg/utils/vari"
)

type Config struct {
	ServerPort int
}

func NewConfig() *Config {
	return &Config{
		ServerPort: vari.DataServicePort,
	}
}
