package logUtils

import (
	"encoding/json"
	"fmt"
	"github.com/easysoft/zendata"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
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

	OutputFileWriter *os.File
	OutputHttpWriter http.ResponseWriter
	OutputFilePath   string // for excel output
)

func PrintExample() {
	if vari.Config.Language == "zh" {
		exampleFile = strings.Replace(exampleFile, "en", "zh", 1)
		usageFile = strings.Replace(usageFile, "en", "zh", 1)
	}

	content, _ := zd.ReadResData(exampleFile)
	fmt.Printf("%s\n", content)
}

func PrintUsage() {
	if vari.Config.Language == "zh" {
		exampleFile = strings.Replace(exampleFile, "en", "zh", 1)
		usageFile = strings.Replace(usageFile, "en", "zh", 1)
	}

	usageBytes, _ := zd.ReadResData(usageFile)
	usage := string(usageBytes)
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

func Info(str string) {
	PrintTo(str)
}
func Infof(str string, args ...interface{}) {
	PrintTo(fmt.Sprintf(str, args))
}
func InfofIfVerbose(str string, args ...interface{}) {
	if vari.Verbose {
		PrintTo(fmt.Sprintf(str, args))
	}
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

func PrintRecord(str string) {
	if OutputFileWriter != nil {
		PrintToFile(str)
	} else if vari.GlobalVars.RunMode == constant.RunModeServerRequest {
		PrintToHttp(str)
	} else {
		PrintToScreen(fmt.Sprintf("%s", str))
	}
}

func PrintLine(line string) {
	if vari.DefType == constant.DefTypeText {
		line += "\n"
	}

	if OutputFileWriter != nil {
		PrintToFile(line)
	} else if vari.GlobalVars.RunMode == constant.RunModeServerRequest {
		PrintToHttp(line)
	} else {
		PrintToScreen(line)
	}

	return
}
func PrintToFile(line string) {
	fmt.Fprint(OutputFileWriter, line)
}
func PrintToHttp(line string) {
	fmt.Fprint(OutputHttpWriter, line)
}
func PrintToScreen(line string) {
	fmt.Print(line)
}

func PrintVersion(appVersion, buildTime, goVersion, gitHash string) {
	fmt.Printf("%s \n", appVersion)
	fmt.Printf("Build TimeStamp: %s \n", buildTime)
	fmt.Printf("GoLang Version: %s \n", goVersion)
	fmt.Printf("Git Commit Hash: %s \n", gitHash)
}

func ConvertUnicode(str []byte) string {
	var a interface{}

	temp := strings.Replace(string(str), "\\\\", "\\", -1)

	err := json.Unmarshal([]byte(temp), &a)

	var msg string
	if err == nil {
		bytes, _ := json.Marshal(a)
		msg = string(bytes)
	} else {
		msg = temp
	}

	return msg
}
