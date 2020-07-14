package constant

import (
	"fmt"
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

	Total = 10
	MaxNumb = 100000 // max number in array

	FormatText = "text"
	FormatJson = "json"
	FormatXml = "xml"
	FormatSql = "sql"
	Formats = []string{FormatText, FormatJson, FormatXml, FormatSql}

	LeftChar rune = '('
	RightChar rune = ')'

	DefaultPort = 8848
	DefaultRoot = "./"

	ResDir  = "data/"

	SqliteDriver  = "sqlite3"
	SqliteSource  = "file:" + ResDir + ".cache/.data.db"
	SqliteTrackTable  = "excel_update"

	ExcelBorder  = `{"border": [{"type":"left","color":"999999","style":1}, {"type":"top","color":"999999","style":1},
		                              {"type":"bottom","color":"999999","style":1}, {"type":"right","color":"999999","style":1}]}`
	ExcelHeader  = `{"fill":{"type":"pattern","pattern":1,"color":["E0EBF5"]}}`
)
