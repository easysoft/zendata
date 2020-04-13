package constant

import (
	"fmt"
	"os"
)

var (
	ConfigVer  = 1
	ConfigFile = fmt.Sprintf("conf%sztf.conf", string(os.PathSeparator))

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

	RequestTypePathInfo = "PATH_INFO"

	UnitTestType []string = []string{"junit", "testng", "phpunit", "pytest", "jtest", "cppunit", "gtest", "qtest"}
	AutoTestType []string = []string{"selenium", "appium"}
)
