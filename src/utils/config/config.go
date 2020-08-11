package configUtils

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	shellUtils "github.com/easysoft/zendata/src/utils/shell"
	stdinUtils "github.com/easysoft/zendata/src/utils/stdin"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"gopkg.in/ini.v1"
	"os"
	"reflect"
)

func InitConfig() {
	vari.ExeDir = fileUtils.GetExeDir()
	vari.WorkDir = fileUtils.GetWorkDir()

	CheckConfigPermission()
	constant.ConfigFile = vari.ExeDir + constant.ConfigFile

	if commonUtils.IsWin() {
		shellUtils.ExeShell("chcp 65001")
	}

	vari.Config = getInst()

	i118Utils.InitI118(vari.Config.Language)
}

func SaveConfig(conf model.Config) error {
	fileUtils.MkDirIfNeeded(vari.ExeDir + "conf")

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

	ini.MapTo(&config, configPath)

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
	//err := syscall.Access(vari.ExeDir, syscall.O_RDWR)
	err := fileUtils.MkDirIfNeeded(vari.ExeDir + "conf")
	if err != nil {
		logUtils.PrintToWithColor(
			fmt.Sprintf("Permission denied, please change the dir %s.", vari.ExeDir), color.FgRed)
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

	//logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("begin_config"), color.FgCyan)

	enCheck := ""
	var numb string
	if conf.Language == "zh" {
		enCheck = "*"
		numb = "1"
	}
	zhCheck := ""
	if conf.Language == "en" {
		zhCheck = "*"
		numb = "2"
	}

	numbSelected := stdinUtils.GetInput("(1|2)", numb, "enter_language", enCheck, zhCheck)

	if numbSelected == "1" {
		conf.Language = "zh"
	} else {
		conf.Language = "en"
	}

	SaveConfig(conf)
	PrintCurrConfig()
}
