package testHelper

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	"github.com/fatih/color"
	"log"
	"os"
)

func BeforeAll() {
	configUtils.InitConfig("")

	//vari.Config.Language = "zh"
	//i118Utils.InitI118(vari.Config.Language)
}

func PreCase() {

	log.SetOutput(&consts.Buf)
	color.Output = &consts.Buf
}

func PostCase() {
	consts.Buf.Reset()
	log.SetOutput(os.Stdout)
}
