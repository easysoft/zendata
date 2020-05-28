package configUtils

import (
	"github.com/easysoft/zendata/src/utils/display"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/i118"
	"github.com/easysoft/zendata/src/utils/vari"
)

func InitConfig() {
	vari.ExeDir = fileUtils.GetExeDir()
	vari.WorkDir = fileUtils.GetWorkDir()

	InitScreenSize()

	i118Utils.InitI118(vari.Config.Language)
}

func InitScreenSize() {
	w, h := display.GetScreenSize()
	vari.ScreenWidth = w
	vari.ScreenHeight = h
}