package constant

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	"os"
)

var (
	ConfigVer  = 1
	ConfigFile = fmt.Sprintf("conf%szdata.conf", string(os.PathSeparator))

	LanguageDefault = "en"
	LanguageEN      = "en"
	LanguageZH      = "zh"

	EnRes = fmt.Sprintf("res%smessages_en.json", string(os.PathSeparator))
	ZhRes = fmt.Sprintf("res%smessages_zh.json", string(os.PathSeparator))

	LogDir = fmt.Sprintf("log%s", string(os.PathSeparator))

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight = 10

	MaxNumb = 100000 // max number in array

	RootDef model.DefData               = model.DefData{}
	ResMap  map[string]model.FieldValue = map[string]model.FieldValue{}

	ResDir     string = "data/"
	ResBuildIn string = ResDir + "system/buildin.yaml"

	SqliteDriver string = "sqlite3"
	SqliteSource string = "file:" + ResDir + ".cache/.data.db"
	ExcelBorder string = `{"border": [{"type":"left","color":"999999","style":1}, {"type":"top","color":"999999","style":1},
		                              {"type":"bottom","color":"999999","style":1}, {"type":"right","color":"999999","style":1}]}`
	ExcelHeader string = `{"fill":{"type":"pattern","pattern":1,"color":["E0EBF5"]}}`
)
