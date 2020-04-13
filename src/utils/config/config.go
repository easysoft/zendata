package configUtils

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	"github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/display"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stdinUtils "github.com/easysoft/zendata/src/utils/stdin"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"gopkg.in/ini.v1"
	"os"
	"reflect"
)

func InitConfig() {
	vari.ZDataDir = fileUtils.GetZTFDir()
	CheckConfigPermission()

	constant.ConfigFile = vari.ZDataDir + constant.ConfigFile
	vari.Config = getInst()

	InitScreenSize()

	i118Utils.InitI118(vari.Config.Language)
}

func InitScreenSize() {
	w, h := display.GetScreenSize()
	vari.ScreenWidth = w
	vari.ScreenHeight = h
}

func SaveConfig(conf model.Config) error {
	fileUtils.MkDirIfNeeded(fileUtils.GetZTFDir() + "conf")

	conf.Version = constant.ConfigVer

	cfg := ini.Empty()
	cfg.ReflectFrom(&conf)

	cfg.SaveTo(constant.ConfigFile)

	vari.Config = ReadCurrConfig()
	return nil
}

func PrintCurrConfig() {
	logUtils.PrintToWithColor("\n"+i118Utils.I118Prt.Sprintf("current_config"), color.FgCyan)

	val := reflect.ValueOf(vari.Config)
	typeOfS := val.Type()
	for i := 0; i < reflect.ValueOf(vari.Config).NumField(); i++ {
		if !commonUtils.IsWin() && i > 4 {
			break
		}

		val := val.Field(i)
		name := typeOfS.Field(i).Name

		fmt.Printf("  %s: %v \n", name, val.Interface())
	}
}

func ReadCurrConfig() model.Config {
	config := model.Config{}

	configPath := constant.ConfigFile

	if !fileUtils.FileExist(configPath) {
		config.Language = "en"
		i118Utils.InitI118("en")

		return config
	}

	ini.MapTo(&config, constant.ConfigFile)

	return config
}

func getInst() model.Config {
	isSetAction := len(os.Args) > 1 && (os.Args[1] == "set" || os.Args[1] == "-set")
	if !isSetAction {
		CheckConfigReady()
	}

	ini.MapTo(&vari.Config, constant.ConfigFile)

	if vari.Config.Version != constant.ConfigVer { // old config file, re-init
		if vari.Config.Language != "en" && vari.Config.Language != "zh" {
			vari.Config.Language = "en"
		}

		SaveConfig(vari.Config)
	}

	return vari.Config
}

func CheckConfigPermission() {
	//err := syscall.Access(vari.ZDataDir, syscall.O_RDWR)

	err := fileUtils.MkDirIfNeeded(vari.ZDataDir + "conf")
	if err != nil {
		logUtils.PrintToWithColor(
			fmt.Sprintf("Permission denied to open %s for write. Please change work dir.", vari.ZDataDir), color.FgRed)
		os.Exit(0)
	}
}

func CheckConfigReady() {
	if !fileUtils.FileExist(constant.ConfigFile) {
		InputForSet()
	}
}

func InputForSet() {
	conf := ReadCurrConfig()

	logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("begin_config"), color.FgCyan)

	enCheck := ""
	var numb string
	if conf.Language == "en" {
		enCheck = "*"
		numb = "1"
	}
	zhCheck := ""
	if conf.Language == "zh" {
		zhCheck = "*"
		numb = "2"
	}

	numbSelected := stdinUtils.GetInput("(1|2)", numb, "enter_language", enCheck, zhCheck)

	if numbSelected == "1" {
		conf.Language = "en"
	} else {
		conf.Language = "zh"
	}

	SaveConfig(conf)
	PrintCurrConfig()
}