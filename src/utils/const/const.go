package constant

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	"os"
)

var (
	ConfigVer  = 1
	ConfigFile = fmt.Sprintf("def%szdata.conf", string(os.PathSeparator))

	UrlZentaoSettings = "zentaoSettings"
	UrlImportProject  = "importProject"
	UrlSubmitResult   = "submitResults"
	UrlReportBug      = "reportBug"

	ExtNameSuite  = "cs"
	ExtNameJson   = "json"
	ExtNameResult = "txt"

	LanguageDefault = "en"
	LanguageEN      = "en"
	LanguageZH      = "zh"

	EnRes = fmt.Sprintf("res%smessages_en.json", string(os.PathSeparator))
	ZhRes = fmt.Sprintf("res%smessages_zh.json", string(os.PathSeparator))

	//ScriptDir = fmt.Sprintf("scripts%s", string(os.PathSeparator))
	LogDir = fmt.Sprintf("log%s", string(os.PathSeparator))

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight = 10

	MaxNumb = 100000 // max number in array

	UnitTestType []string = []string{"junit", "testng", "phpunit", "pytest", "jtest", "cppunit", "gtest", "qtest"}
	AutoTestType []string = []string{"selenium", "appium"}

	Power3 = 255 * 255 * 255
	Power2 = 255 * 255
	Power1 = 255

	ResDir string = ""
	Definition model.Definition = model.Definition{}
	LoadedFields map[string]model.Field = map[string]model.Field{}

	SqliteDriver string = "sqlite3"
	SqliteSource string = "file:def/data/data.sqlite3"
)
