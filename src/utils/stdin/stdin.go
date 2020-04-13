package stdinUtils

import (
	"bufio"
	"fmt"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	langUtils "github.com/easysoft/zentaoatf/src/utils/lang"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func InputForCheckout(productId *string, moduleId *string, suiteId *string, taskId *string,
	independentFile *bool, scriptLang *string) {

	var numb string

	productCheckbox := ""
	suiteCheckbox := ""
	taskCheckbox := ""

	if *productId != "" {
		productCheckbox = "*"
		numb = "1"
	} else if *suiteId != "" {
		suiteCheckbox = "*"
		numb = "2"
	} else if *taskId != "" {
		taskCheckbox = "*"
		numb = "3"
	}

	coType := GetInput("(1|2|3)", numb, "enter_co_type", productCheckbox, suiteCheckbox, taskCheckbox)

	coType = strings.ToLower(coType)
	if coType == "1" {
		*productId = GetInput("\\d+", *productId,
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("product_id")+": "+*productId)

		*moduleId = GetInput("\\d*", *moduleId,
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("module_id")+": "+*moduleId)

	} else if coType == "2" {
		*suiteId = GetInput("\\d+", *suiteId,
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("suite_id")+": "+*suiteId)
	} else if coType == "3" {
		*taskId = GetInput("\\d+", *taskId,
			i118Utils.I118Prt.Sprintf("pls_enter")+" "+i118Utils.I118Prt.Sprintf("task_id")+": "+*taskId)
	}

	InputForBool(independentFile, false, "enter_co_independent")

	numbs, names, labels := langUtils.GetSupportLanguageOptions(nil)
	fmtParam := make([]string, 0)
	dft := ""
	for idx, label := range labels {
		if names[idx] == *scriptLang {
			dft = strconv.Itoa(idx + 1)
			label += " *"
		}
		fmtParam = append(fmtParam, label)
	}

	langStr := GetInput("("+strings.Join(numbs, "|")+")", dft, "enter_co_language", strings.Join(fmtParam, "\n"))
	langNumb, _ := strconv.Atoi(langStr)

	*scriptLang = names[langNumb-1]
}

func InputForDir(dir *string, dft string, i118Key string) {
	*dir = GetInput("is_dir", dft, "enter_dir", i118Utils.I118Prt.Sprintf(i118Key))
}

func InputForBool(in *bool, defaultVal bool, fmtStr string, fmtParam ...interface{}) {
	str := GetInput("(yes|no|y|n|)", "", fmtStr, fmtParam...)

	if str == "" {
		*in = defaultVal

		msg := ""
		if *in {
			msg = "yes"
		} else {
			msg = "no"
		}
		logUtils.PrintTo(msg)
		return
	}

	if str == "y" && str != "yes" {
		*in = true
	} else {
		*in = false
	}
}

func GetInput(regx string, defaultVal string, fmtStr string, params ...interface{}) string {
	var ret string

	msg := i118Utils.I118Prt.Sprintf(fmtStr, params...)

	for {
		logUtils.PrintToWithColor("\n"+msg, color.FgCyan)
		// fmt.Scanln(&ret)
		Scanf(&ret)

		//logUtils.PrintToWithColor(fmt.Sprintf("%v", ret), -1)

		if strings.TrimSpace(ret) == "" && defaultVal != "" {
			ret = defaultVal

			logUtils.PrintTo(ret)
		}

		temp := strings.ToLower(ret)
		if temp == "exit" {
			color.Unset()
			os.Exit(0)
		}

		//logUtils.PrintToWithColor(ret, -1)

		if regx == "" {
			return ret
		}

		var pass bool
		var msg string
		if regx == "is_dir" {
			pass = fileUtils.IsDir(ret)
			msg = "dir_not_exist"
		} else {
			pass, _ = regexp.MatchString("^"+regx+"$", temp)
			msg = "invalid_input"
		}

		if pass {
			return ret
		} else {
			ret = ""
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf(msg), color.FgRed)
		}
	}
}

func GetInputForScriptInterpreter(defaultVal string, fmtStr string, params ...interface{}) string {
	var ret string

	msg := i118Utils.I118Prt.Sprintf(fmtStr, params...)

	for {
		logUtils.PrintToWithColor(msg, color.FgCyan)
		Scanf(&ret)

		ret = strings.TrimSpace(ret)

		if ret == "" && defaultVal != "" {
			ret = defaultVal

			logUtils.PrintToWithColor(ret, -1)
		}

		if ret == "exit" {
			color.Unset()
			os.Exit(0)
		}

		if ret == "" { // ignore to set
			return "-"
		}

		sep := string(os.PathSeparator)
		if sep == `\` {
			sep = `\\`
		}
		reg := fmt.Sprintf(".*%s+[^%s]+", sep, sep)
		pass, _ := regexp.MatchString(reg, ret)
		if pass {
			return ret
		} else {
			ret = ""
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("invalid_input"), color.FgRed)
		}
	}
}

func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}
