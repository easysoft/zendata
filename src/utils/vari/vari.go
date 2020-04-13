package vari

import (
	"github.com/awesome-gocui/gocui"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/const"
)

var (
	Config         = model.Config{}
	Cui            *gocui.Gui
	MainViewHeight int

	RunMode        constant.RunMode
	ZTFDir         string
	LogDir         string
	RunFromCui     bool
	UnitTestType   string
	UnitTestTool   string
	UnitTestResult string
	ProductId      string

	SessionVar  string
	SessionId   string
	RequestType string
	RequestFix  string = ""

	ScriptExtToNameMap map[string]string
	CurrScriptFile     string // scripts/tc-001.py
	CurrResultDate     string // 2019-08-15T173802
	CurrCaseId         int    // 2019-08-15T173802

	ScreenWidth     int
	ScreenHeight    int
	ZentaoBugFileds model.ZentaoBugFileds

	//ZentaoCaseFileds model.ZentaoCaseFileds

	CurrBug        model.Bug
	CurrBugStepIds string

	Verbose     bool
	Interpreter string
)
