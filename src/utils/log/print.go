package logUtils

import (
	"encoding/json"
	"fmt"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"os"
	"regexp"
	"strings"
)

var (
	exampleFile  = fmt.Sprintf("res%sdoc%ssample.yaml", string(os.PathSeparator), string(os.PathSeparator))
	usageFile  = fmt.Sprintf("res%sdoc%susage.txt", string(os.PathSeparator), string(os.PathSeparator))
)

func PrintExample() {
	if vari.Config.Language == "en" {
		exampleFile = strings.Replace(exampleFile, ".yaml", "_en.yaml", -1)
		usageFile = strings.Replace(usageFile, ".txt", "_en.txt", -1)
	}

	content := fileUtils.ReadResData(exampleFile)
	fmt.Printf("%s\n", content)
}

func PrintUsage() {
	if vari.Config.Language == "en" {
		exampleFile = strings.Replace(exampleFile, ".yaml", "_en.yaml", -1)
		usageFile = strings.Replace(usageFile, ".txt", "_en.txt", -1)
	}

	usage := fileUtils.ReadResData(usageFile)
	exeFile := "zd"
	if commonUtils.IsWin() {
		exeFile += ".exe"
	}

	if !commonUtils.IsWin() {
		regx, _ := regexp.Compile(`\\`)
		usage = regx.ReplaceAllString(usage, "/")

		regx, _ = regexp.Compile(`zd.exe`)
		usage = regx.ReplaceAllString(usage, "zd")

		regx, _ = regexp.Compile(`d:`)
		usage = regx.ReplaceAllString(usage, "/home/user")
	}
	fmt.Printf("%s\n", usage)
}

func PrintTo(str string) {
	output := color.Output
	fmt.Fprint(output, str+"\n")
}

func PrintToWithColor(msg string, attr color.Attribute) {
	output := color.Output

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		color.New(attr).Fprintf(output, msg+"\n")
	}
}

func PrintToCmd(msg string, attr color.Attribute) {
	output := color.Output

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		clr := color.New(attr)
		clr.Fprint(output, msg+"\n")
	}
}

func PrintUnicode(str []byte) {
	var a interface{}

	temp := strings.Replace(string(str), "\\\\", "\\", -1)

	err := json.Unmarshal([]byte(temp), &a)

	var msg string
	if err == nil {
		msg = fmt.Sprint(a)
	} else {
		msg = temp
	}

	PrintToCmd(msg, -1)
}
