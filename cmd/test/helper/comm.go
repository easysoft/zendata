package testHelper

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	commandConfig "github.com/easysoft/zendata/internal/command/config"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"log"
	"os"
	"strings"
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

func SetFields(fields string) {
	arr := strings.Split(fields, ",")
	vari.GlobalVars.ExportFields = arr
}

func SetTotal(total int) {
	vari.GlobalVars.Total = total
}

func SetTrim(val bool) {
	vari.GlobalVars.Trim = val
}

func SetHuman(val bool) {
	vari.GlobalVars.Human = val
}

func SetRecursive(val bool) {
	vari.GlobalVars.Recursive = val
}
