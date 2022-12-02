package constant

import (
	"fmt"
	"os"
	"os/user"
)

const (
	AppName = "zd"
)

var (
	ApiPath = "/api/v1/admin"
	WsPath  = ApiPath + "/ws"

	PthSep = string(os.PathSeparator)

	ConfigVer      = 1.0
	userProfile, _ = user.Current()

	LanguageDefault = "en"
	LanguageEN      = "en"
	LanguageZH      = "zh"

	CachePrefix  = "cache_"
	CachePostfix = "_is_num"
	EnRes        = fmt.Sprintf("res%sen%smessages.json", PthSep, PthSep)
	ZhRes        = fmt.Sprintf("res%szh%smessages.json", PthSep, PthSep)

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight = 10

	MaxNumb         = 100000
	MaxNumbForAsync = 100

	FormatText  = "text"
	FormatJson  = "json"
	FormatXml   = "xml"
	FormatSql   = "sql"
	FormatExcel = "xlsx"
	FormatCsv   = "csv"
	FormatProto = "proto"
	FormatData  = "data"
	Formats     = []string{FormatText, FormatJson, FormatXml, FormatSql, FormatExcel, FormatCsv, FormatProto}

	ModeParallel       = "parallel"
	ModeRecursive      = "recursive"
	ModeParallelShort  = "p"
	ModeRecursiveShort = "r"
	Modes              = []string{ModeParallel, ModeRecursive, ModeParallelShort, ModeRecursiveShort}

	DefTypeText    = "text"
	DefTypeArticle = "article"
	DefTypeImage   = "image"
	DefTypeVoice   = "voice"
	DefTypeVideo   = "video"

	FieldTypeList      = "list"
	FieldTypeTimestamp = "timestamp"
	FieldTypeUlid      = "ulid"
	FieldTypeArticle   = "article"

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
	SqliteFile   = "file:" + TmpDir + "cache/.data.db"

	SqliteTrackTable = "excel_update"

	ExcelBorder = `{"border": [{"type":"left","color":"999999","style":1}, {"type":"top","color":"999999","style":1},
		                              {"type":"bottom","color":"999999","style":1}, {"type":"right","color":"999999","style":1}]}`
	ExcelHeader = `{"fill":{"type":"pattern","pattern":1,"color":["E0EBF5"]}}`

	TablePrefix = "zd_"
	PageSize    = 15

	// database type  [added by leo 2022/05/10]
	DBTypeMysql     = "mysql"
	DBTypeSqlServer = "sqlserver"
	DBTypeOracle    = "oracle"
)
