package testHelper

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	commandConfig "github.com/easysoft/zendata/internal/command/config"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"log"
	"os"
)

func BeforeAll() {
	configUtils.InitConfig("")
	vari.DB, _ = commandConfig.NewGormDB()

	vari.GlobalVars.Total = 10
}

func PreCase() {
	log.SetOutput(&consts.Buf)
	color.Output = &consts.Buf
}

func PostCase() {
	consts.Buf.Reset()
	log.SetOutput(os.Stdout)
}
