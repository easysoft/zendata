package testHelper

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"log"
	"os"
)

func PreCase() {
	log.SetOutput(&consts.Buf)
	color.Output = &consts.Buf

	vari.Config.Language = "zh"
	i118Utils.InitI118(vari.Config.Language)
}

func PostCase() {
	consts.Buf.Reset()
	log.SetOutput(os.Stdout)
}
