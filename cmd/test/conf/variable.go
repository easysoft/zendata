package testConsts

import (
	"os"
	"runtime"
	"strings"
)

const (
	FilePthSep = string(os.PathSeparator)
)

var (
	NewLine  = "\n"
	RootPath = ""

	Token = ""
)

func init() {
	if runtime.GOOS == "windows" {
		NewLine = "\r\n"
	}
	RootPath, _ = os.Getwd()
	if strings.Index(RootPath, "test") != -1 {
		RootPath = RootPath[:strings.Index(RootPath, "test")]
	}
	if RootPath[len(RootPath)-1:] != FilePthSep {
		RootPath += FilePthSep
	}
}
