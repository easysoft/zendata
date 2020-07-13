package configUtils

import (
	"github.com/easysoft/zendata/src/utils/display"
	"github.com/easysoft/zendata/src/utils/i118"
	"github.com/easysoft/zendata/src/utils/vari"
)

func InitConfig() {
	InitScreenSize()

	i118Utils.InitI118(vari.Config.Language)
}

func InitScreenSize() {
	w, h := display.GetScreenSize()
	vari.ScreenWidth = w
	vari.ScreenHeight = h
}