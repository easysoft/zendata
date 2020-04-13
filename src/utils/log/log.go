package logUtils

import (
	"fmt"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

var Logger *logrus.Logger

func GetWholeLine(msg string, char string) string {
	prefixLen := (vari.ScreenWidth - utf8.RuneCountInString(msg) - 2) / 2
	if prefixLen <= 0 { // no width in debug mode
		prefixLen = 6
	}
	postfixLen := vari.ScreenWidth - utf8.RuneCountInString(msg) - 2 - prefixLen - 1
	if postfixLen <= 0 { // no width in debug mode
		postfixLen = 6
	}

	preFixStr := strings.Repeat(char, prefixLen)
	postFixStr := strings.Repeat(char, postfixLen)

	return fmt.Sprintf("%s %s %s", preFixStr, msg, postFixStr)
}

func ColoredStatus(status string) string {
	temp := strings.ToLower(status)

	switch temp {
	case "pass":
		return color.GreenString(i118Utils.I118Prt.Sprintf(temp))
	case "fail":
		return color.RedString(i118Utils.I118Prt.Sprintf(temp))
	case "skip":
		return color.YellowString(i118Utils.I118Prt.Sprintf(temp))
	}

	return status
}

func InitLogger() *logrus.Logger {
	vari.LogDir = fileUtils.GetLogDir()

	if Logger != nil {
		return Logger
	}

	Logger = logrus.New()
	Logger.Out = ioutil.Discard

	pathMap := lfshook.PathMap{
		logrus.WarnLevel:  vari.LogDir + "log.txt",
		logrus.ErrorLevel: vari.LogDir + "result.txt",
	}

	Logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&MyFormatter{},
	))

	Logger.SetFormatter(&MyFormatter{})

	return Logger
}

func Screen(msg string) {
	PrintTo(msg)
}
func Log(msg string) {
	Logger.Warnln(msg)
}
func Result(msg string) {
	Logger.Errorln(msg)
}

func ScreenAndResult(msg string) {
	Screen(msg)
	Result(msg)
}

type MyFormatter struct {
	logrus.TextFormatter
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}
