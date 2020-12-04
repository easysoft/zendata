package constant

import (
	"fmt"
	"os"
	"os/user"
)

var (
	PthSep = string(os.PathSeparator)

	ConfigVer      = 1
	userProfile, _ = user.Current()

	ConfigFile = fmt.Sprintf("%s%s.zd.conf", userProfile.HomeDir, PthSep)

	LanguageDefault = "en"
	LanguageEN      = "en"
	LanguageZH      = "zh"

	EnRes = fmt.Sprintf("res%sen%smessages.json", string(os.PathSeparator), string(os.PathSeparator))
	ZhRes = fmt.Sprintf("res%szh%smessages.json", string(os.PathSeparator), string(os.PathSeparator))

	LogDir = fmt.Sprintf("log%s", string(os.PathSeparator))

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight = 10

	MaxNumb = 100000 // max number in array

	FormatText  = "text"
	FormatJson  = "json"
	FormatXml   = "xml"
	FormatSql   = "sql"
	FormatExcel = "xlsx"
	FormatCsv   = "csv"
	FormatData  = "data"
	Formats     = []string{FormatText, FormatJson, FormatXml, FormatSql, FormatExcel, FormatCsv}

	ModeParallel       = "parallel"
	ModeRecursive      = "recursive"
	ModeParallelShort  = "p"
	ModeRecursiveShort = "r"
	Modes              = []string{ModeParallel, ModeRecursive, ModeParallelShort, ModeRecursiveShort}

	ConfigTypeText    = "text"
	ConfigTypeArticle = "article"
	ConfigTypeImage   = "image"
	ConfigTypeVoice   = "voice"
	ConfigTypeVideo   = "video"

	FieldTypeList      = "list"
	FieldTypeTimestamp = "timestamp"

	LeftBrackets  rune = '('
	RightBrackets rune = ')'
	Backtick      rune = '`'

	DefaultPort   = 8848
	DefaultRoot   = "./"
	DefaultNumber = 10

	ResDirData  = "data"
	ResDirYaml  = "yaml"
	ResDirUsers = "users"
	ResKeys     = []string{ResDirData, ResDirYaml, ResDirUsers}

	ResTypeDef       = "def" // only used for refer type
	ResTypeConfig    = "config"
	ResTypeRanges    = "ranges"
	ResTypeInstances = "instances"
	ResTypeYaml      = "yaml"
	ResTypeExcel     = "excel"
	ResTypeText      = "text"
	ResTypeValue     = "value"
	ResTypes         = []string{ResTypeConfig, ResTypeRanges, ResTypeInstances, ResTypeExcel, ResTypeText, ResTypeValue}

	TmpDir = "tmp/"

	SqliteDriver = "sqlite3"
	SqliteData   = "file:" + TmpDir + "cache/.data.db"

	SqliteTrackTable = "excel_update"

	ExcelBorder = `{"border": [{"type":"left","color":"999999","style":1}, {"type":"top","color":"999999","style":1},
		                              {"type":"bottom","color":"999999","style":1}, {"type":"right","color":"999999","style":1}]}`
	ExcelHeader = `{"fill":{"type":"pattern","pattern":1,"color":["E0EBF5"]}}`

	TablePrefix = "zd_"
	PageSize    = 15
)
