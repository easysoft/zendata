package logUtils

import (
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	exampleFile = fmt.Sprintf("res%sen%ssample.yaml", string(os.PathSeparator), string(os.PathSeparator))
	usageFile   = fmt.Sprintf("res%sen%susage.txt", string(os.PathSeparator), string(os.PathSeparator))

	FileWriter *os.File
	HttpWriter http.ResponseWriter
	FilePath   string // for excel output
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

		regx, _ = regexp.Compile(`d:\\zd\\config        `)
		usage = regx.ReplaceAllString(usage, "/home/user/zd/config")
	}
	fmt.Printf("%s\n", usage)
}

func PrintTo(str string) {
	output := color.Output
	fmt.Fprint(output, str+"\n")
}
func PrintToWithoutNewLine(str string) {
	output := color.Output
	fmt.Fprint(output, str)
}

func PrintToWithColor(msg string, attr color.Attribute) {
	output := color.Output

	if attr == -1 {
		fmt.Fprint(output, msg+"\n")
	} else {
		color.New(attr).Fprintf(output, msg+"\n")
		//color.New(attr).Println(output, msg)
	}
}

func PrintErrMsg(msg string) {
	PrintToWithColor(msg, color.FgCyan)
}

func PrintLine(line string) {
	if vari.DefType == constant.DefTypeText {
		line += "\n"
	}

	if FileWriter != nil {
		PrintToFile(line)
	} else if vari.RunMode == constant.RunModeServerRequest {
		PrintToHttp(line)
	} else {
		PrintToScreen(line)
	}

	return
}
func PrintToFile(line string) {
	fmt.Fprint(FileWriter, line)
}
func PrintToHttp(line string) {
	fmt.Fprint(HttpWriter, line)
}
func PrintToScreen(line string) {
	fmt.Print(line)
}
