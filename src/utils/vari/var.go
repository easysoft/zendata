package vari

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/const"
)

var (
	Config      = model.Config{Version: 1, Language: "en"}

	RunMode      constant.RunMode

	ExeDir       string
	WorkDir      string
	InputDir     string

	LogDir       string
	ScreenWidth  int
	ScreenHeight int

	RequestType  string
	Verbose     bool
	Interpreter string

	WithHead bool
	HeadSep string

	Length int
	LeftPad string
	RightPad string

	JsonResp string = "[]"
	Ip string
	Port int
)
