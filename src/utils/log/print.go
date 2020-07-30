package logUtils

import (
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
	exampleFile  = fmt.Sprintf("res%sen%ssample.yaml", string(os.PathSeparator), string(os.PathSeparator))
	usageFile  = fmt.Sprintf("res%sen%susage.txt", string(os.PathSeparator), string(os.PathSeparator))
)

func PrintExample() {
	if vari.Config.Language == "zh" {
		exampleFile = strings.Replace(exampleFile, "en", "zh", 1)
		usageFile = strings.Replace(usageFile, "en", "zh", 1)
	}

	content := fileUtils.ReadResData(exampleFile)
	fmt.Printf("%s\n", content)
}

func PrintUsage() {
	if vari.Config.Language == "zh" {
		exampleFile = strings.Replace(exampleFile, "en", "zh", 1)
		usageFile = strings.Replace(usageFile, "en", "zh", 1)
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
