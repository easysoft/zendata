package configUtils

import (
	"database/sql"
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
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
)

func InitDB() (db *sql.DB, err error) {
	db, err = sql.Open(constant.SqliteDriver, constant.SqliteData)
	err = db.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
	if err != nil {
		logUtils.PrintErrMsg(
			fmt.Sprintf("Error on opening db %s, error is %s", constant.SqliteData, err.Error()))
	}
	return
}

func InitConfig() {
	vari.ExeDir = fileUtils.GetExeDir()
	vari.WorkDir = fileUtils.GetWorkDir()

	CheckConfigPermission()

	if commonUtils.IsWin() {
		shellUtils.ExeShell("chcp 65001")
	}

	vari.Config = getInst()

	i118Utils.InitI118(vari.Config.Language)
}

func SaveConfig(conf model.Config) error {
	fileUtils.MkDirIfNeeded(filepath.Dir(constant.ConfigFile))

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
	err := fileUtils.MkDirIfNeeded(filepath.Dir(constant.ConfigFile))
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

	// set lang
	langNo := stdinUtils.GetInput("(1|2)", numb, "enter_language", enCheck, zhCheck)
	if langNo == "1" {
		conf.Language = "zh"
	} else {
		conf.Language = "en"
	}

	// set PATH environment vari
	var addToPath bool
	if commonUtils.IsWin() {
		addToPath = true
		// stdinUtils.InputForBool(&addToPath, true, "add_to_path_win")
	} else {
		stdinUtils.InputForBool(&addToPath, true, "add_to_path_linux")
	}

	if addToPath {
		AddZdToPath()
	}

	SaveConfig(conf)
	PrintCurrConfig()
}

func AddZdToPath() {
	userProfile, _ := user.Current()
	home := userProfile.HomeDir

	if commonUtils.IsWin() {
		addZdToPathWin(home)
	} else {
		addZdToPathLinux(home)
	}
}

func addZdToPathWin(home string) {
	pathVar := os.Getenv("PATH")
	if strings.Contains(pathVar, vari.ExeDir) { return }

	cmd := `setx Path "%%Path%%;` + vari.ExeDir + `"`
	logUtils.PrintToWithColor("\n" + i118Utils.I118Prt.Sprintf("add_to_path_tips_win", cmd), color.FgRed)

	// TODO: fix the space issue
	//out, err := shellUtils.ExeShell(cmd)
	//
	//if err == nil {
	//	msg := i118Utils.I118Prt.Sprintf("add_to_path_success_win")
	//	logUtils.PrintToWithColor(msg, color.FgRed)
	//} else {
	//	logUtils.PrintToWithColor(
	//		i118Utils.I118Prt.Sprintf("fail_to_exec_cmd", cmd, err.Error() + ": " + out), color.FgRed)
	//}
}

func addZdToPathLinux(home string) {
	path := fmt.Sprintf("%s%s%s", home, constant.PthSep, ".bash_profile")

	content := fileUtils.ReadFile(path)
	if strings.Contains(content, vari.ExeDir) { return }

	cmd := fmt.Sprintf("echo 'export PATH=$PATH:%s' >> %s", vari.ExeDir, path)
	out, err := shellUtils.ExeShell(cmd)

	if err == nil {
		msg := i118Utils.I118Prt.Sprintf("add_to_path_success_linux", path)
		logUtils.PrintToWithColor(msg, color.FgRed)
	} else {
		logUtils.PrintToWithColor(
			i118Utils.I118Prt.Sprintf("fail_to_exec_cmd", cmd, err.Error() + ": " + out), color.FgRed)
	}
}

