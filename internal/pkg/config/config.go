package configUtils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	shellUtils "github.com/easysoft/zendata/pkg/utils/shell"
	stdinUtils "github.com/easysoft/zendata/pkg/utils/stdin"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"gopkg.in/ini.v1"
	"gorm.io/gorm"
)

func InitConfig(root string) {
	var err error = nil
	vari.WorkDir = fileUtils.GetWorkDir()
	vari.CfgFile = vari.WorkDir + ".zd.conf"

	vari.Config = getInst()

	i118Utils.InitI118(vari.Config.Language)

	if root != "" {
		if !fileUtils.IsAbsPath(root) {
			root, err = filepath.Abs(root)
			if err != nil {
				logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("root_invalid", root), color.FgRed)
				os.Exit(1)
			}
		}
		vari.WorkDir = fileUtils.AddSepIfNeeded(root)
	}

	temp := filepath.Join(vari.WorkDir, "tmp", "cache")
	if !fileUtils.FileExist(temp) {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("root_invalid", root), color.FgRed)
		os.Exit(1)
	}

	CheckConfigPermission()

	if commonUtils.IsWin() {
		shellUtils.Exec("chcp 65001")
	}

	//logUtils.PrintToWithColor("workdir = "+vari.ZdPath, color.FgCyan)
	consts.SqliteFile = strings.Replace(consts.SqliteFile, "file:", "file:"+vari.WorkDir, 1)
	//logUtils.PrintToWithColor("dbfile = "+consts.SqliteFile, color.FgCyan)
}

func SaveConfig(conf model.Config) error {
	fileUtils.MkDirIfNeeded(filepath.Dir(vari.CfgFile))

	if conf.Version == 0 {
		conf.Version = consts.ConfigVer
	}

	cfg := ini.Empty()
	cfg.ReflectFrom(&conf)

	cfg.SaveTo(vari.CfgFile)

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

	if !fileUtils.FileExist(vari.CfgFile) {
		config.Language = "en"
		i118Utils.InitI118("en")

		return config
	}

	ini.MapTo(&config, vari.CfgFile)

	return config
}

func getInst() model.Config {
	isSetAction := len(os.Args) > 1 && (os.Args[1] == "set" || os.Args[1] == "-set")
	if !isSetAction {
		CheckConfigReady()
	}

	ini.MapTo(&vari.Config, vari.CfgFile)

	if vari.Config.Version != consts.ConfigVer { // old config file, re-init
		if vari.Config.Language != "en" && vari.Config.Language != "zh" {
			vari.Config.Language = "en"
		}

		SaveConfig(vari.Config)
	}

	return vari.Config
}

func CheckConfigPermission() {
	//err := syscall.Access(vari.ExeDir, syscall.O_RDWR)
	err := fileUtils.MkDirIfNeeded(filepath.Dir(vari.CfgFile))
	if err != nil {
		logUtils.PrintToWithColor(
			fmt.Sprintf("Permission denied, please change the dir %s.", vari.WorkDir), color.FgRed)
		os.Exit(0)
	}
}

func CheckConfigReady() {
	if !fileUtils.FileExist(vari.CfgFile) {
		//logUtils.PrintTo(vari.CfgFile + "no exist")
		if vari.GlobalVars.RunMode == consts.RunModeServer {
			conf := model.Config{Language: "zh", Version: 1}
			SaveConfig(conf)
		} else {
			InputForSet()
		}
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
		addZdToPathEnvVarWin(home)
	} else {
		addZdToPathEnvVarForLinux(home)
	}
}

func addZdToPathEnvVarWin(home string) {
	pathVar := os.Getenv("PATH")
	if strings.Contains(pathVar, vari.WorkDir) {
		return
	}

	cmd := `setx Path "%%Path%%;` + vari.WorkDir + `"`
	logUtils.PrintToWithColor("\n"+i118Utils.I118Prt.Sprintf("add_to_path_tips_win", cmd), color.FgRed)

	// TODO: fix the space issue
	//out, err := shellUtils.Exec(cmd)
	//
	//if err == nil {
	//	msg := i118Utils.I118Prt.Sprintf("add_to_path_success_win")
	//	logUtils.PrintToWithColor(msg, color.FgRed)
	//} else {
	//	logUtils.PrintToWithColor(
	//		i118Utils.I118Prt.Sprintf("fail_to_exec_cmd", cmd, err.Error() + ": " + out), color.FgRed)
	//}
}

func addZdToPathEnvVarForLinux(home string) {
	path := fmt.Sprintf("%s%s%s", home, consts.PthSep, ".bash_profile")

	content := fileUtils.ReadFile(path)
	if strings.Contains(content, vari.WorkDir) {
		return
	}

	cmd := fmt.Sprintf("echo 'export PATH=$PATH:%s' >> %s", vari.WorkDir, path)
	out, err := shellUtils.Exec(cmd)

	if err == nil {
		msg := i118Utils.I118Prt.Sprintf("add_to_path_success_linux", path)
		logUtils.PrintToWithColor(msg, color.FgRed)
	} else {
		logUtils.PrintToWithColor(
			i118Utils.I118Prt.Sprintf("fail_to_exec_cmd", cmd, err.Error()+": "+out), color.FgRed)
	}
}

func isDataInit(gormDb *gorm.DB) bool {
	def := model.ZdDef{}
	err := gormDb.Find(def).Error

	if err == nil {
		return true
	} else {
		return false
	}
}
