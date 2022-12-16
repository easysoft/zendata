package vari

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"gorm.io/gorm"
	"time"
)

type GenVarType struct {
	Total int

	OutputFile   string
	OutputFormat string
	ExportFields []string
	ColIsNumArr  []bool

	Table        string
	DBType       string // database type
	DBDsn        string
	DBDsnParsing DBDsnData
	DBClear      bool

	Human     bool
	Trim      bool
	Recursive bool

	ConfigFileDir string

	DefData model.DefData
	ResData map[string]map[string][]interface{}

	CacheResFileToMap                  map[string]map[string][]interface{}
	RandFieldSectionPathToValuesMap    map[int]map[string]interface{}
	RandFieldSectionShortKeysToPathMap map[int]string

	StartTime time.Time
	EndTime   time.Time
}

var (
	GlobalVars = GenVarType{
		DefData:      model.DefData{},
		OutputFormat: consts.FormatText,

		CacheResFileToMap:                  map[string]map[string][]interface{}{},
		RandFieldSectionPathToValuesMap:    map[int]map[string]interface{}{},
		RandFieldSectionShortKeysToPathMap: map[int]string{},
	}
)

var (
	Config = model.Config{Version: 1, Language: "en"}
	DB     *gorm.DB

	RunMode consts.RunMode

	WorkDir      string
	ZdPath       string
	CurrFilePath string

	CfgFile      string
	LogDir       string
	ScreenWidth  int
	ScreenHeight int

	RequestType string
	Verbose     bool
	Interpreter string

	CacheParam string

	DefType string

	ProtoCls string

	JsonResp string = "[]"
	Ip       string
	Port     int

	ResLoading  = false
	Res         = map[string]map[string][]string{}
	TopFieldMap = map[string]model.DefField{}

	CacheResFileToMap  = map[string]map[string][]string{}
	CacheResFileToName = map[string]string{}

	AgentLogDir string
)

// parsing from  DBDsn [added by Leo [2022/5/5]]
type DBDsnData struct {
	Driver   string
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
	Code     string
}
