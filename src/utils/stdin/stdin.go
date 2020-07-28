package stdinUtils

import (
	"bufio"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
)

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

func Scanf(a *string) {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	*a = string(data)
}
