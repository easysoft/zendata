package vari

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"gorm.io/gorm"
)

type GenVarType struct {
	Total         int
	ConfigFileDir string

	DefData model.DefData
	ResData map[string]map[string][]interface{}

	CacheResFileToMap map[string]map[string][]interface{}
}

var (
	GenVars = GenVarType{
		DefData:           model.DefData{},
		CacheResFileToMap: map[string]map[string][]interface{}{},
	}
)

var (
	Config = model.Config{Version: 1, Language: "en"}
	DB     *gorm.DB

	RunMode constant.RunMode

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

	WithHead   bool
	Human      bool
	Trim       bool
	Recursive  bool
	CacheParam string
	Format     = constant.FormatText

	Out          string
	Table        string
	DefType      string
	Server       string // database type
	DBDsn        string
	DBDsnParsing DBDsnData
	DBClear      bool
	ProtoCls     string

	JsonResp string = "[]"
	Ip       string
	Port     int

	ResLoading                         = false
	Res                                = map[string]map[string][]string{}
	RandFieldSectionPathToValuesMap    = map[int]map[string]interface{}{}
	RandFieldSectionShortKeysToPathMap = map[int]string{}
	TopFieldMap                        = map[string]model.DefField{}

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
