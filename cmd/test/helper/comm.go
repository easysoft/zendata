package testHelper

import (
	"log"
	"os"

	"github.com/easysoft/zendata/cmd/test/consts"
	commandConfig "github.com/easysoft/zendata/internal/command/config"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
)

func BeforeAll() {
	configUtils.InitConfig("")
	vari.DB, _ = commandConfig.NewGormDB()
}

func PreCase() {
	log.SetOutput(&consts.Buf)
	color.Output = &consts.Buf

	vari.GlobalVars.Total = 10
	vari.GlobalVars.ExportFields = []string{""}
	vari.GlobalVars.Output = ""
	vari.GlobalVars.Trim = false
	vari.GlobalVars.Human = false
	vari.GlobalVars.Recursive = false

	vari.GlobalVars.DBType = ""
	vari.GlobalVars.Table = ""
	vari.ProtoCls = ""
}

func PostCase() {
	consts.Buf.Reset()
	log.SetOutput(os.Stdout)
}
