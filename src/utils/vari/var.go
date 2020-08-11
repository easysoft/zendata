package vari

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/const"
)

var (
	Config      = model.Config{Version: 1, Language: "en"}

	RunMode      constant.RunMode

	ExeDir       string

	LogDir       string
	ScreenWidth  int
	ScreenHeight int

	RequestType  string
	Verbose     bool
	Interpreter string

	Total int
	WithHead bool
	Human bool
	Trim bool
	Recursive bool

	JsonResp string = "[]"
	Ip string
	Port int

	Def = model.DefData{}
	Res = map[string]map[string][]string{}
	RandFieldNameToValuesMap = map[string]map[string]interface{}{}
	TopFiledMap  = map[string]model.DefField{}

	DefaultDir string
	ConfigDir string
)
