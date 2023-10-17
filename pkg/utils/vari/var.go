package vari

import (
	"time"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/model"
	"gorm.io/gorm"
)

type GenVarType struct {
	RunMode consts.RunMode
	Total   int

	Output           string
	OutputFormat     string
	TopFieldMap      map[string]domain.DefField
	ExportFields     []string
	ExportChildField string
	ColIsNumArr      []bool

	Table        string
	DBType       string // database type
	DBDsn        string
	DBDsnParsing DBDsnData
	DBClear      bool
	MockDir      string

	Human     bool
	Trim      bool
	Recursive bool

	ConfigFileDir string

	DefData domain.DefData
	ResData map[string]map[string][]interface{}

	CacheResFileToMap                  map[string]map[string][]interface{}
	RandFieldSectionPathToValuesMap    map[int]map[string]interface{}
	RandFieldSectionShortKeysToPathMap map[int]string

	FieldNameToValuesMap map[string][]interface{}
	FieldNameToFieldMap  map[string]domain.DefField

	StartTime time.Time
	EndTime   time.Time
}

var (
	GlobalVars = GenVarType{
		DefData:      domain.DefData{},
		OutputFormat: consts.FormatText,
		TopFieldMap:  map[string]domain.DefField{},

		CacheResFileToMap:                  map[string]map[string][]interface{}{},
		RandFieldSectionPathToValuesMap:    map[int]map[string]interface{}{},
		RandFieldSectionShortKeysToPathMap: map[int]string{},

		FieldNameToValuesMap: map[string][]interface{}{},
		FieldNameToFieldMap:  map[string]domain.DefField{},

		RunMode: consts.RunModeGen,
	}
)

var (
	Config = model.Config{Version: 1, Language: "en"}
	DB     *gorm.DB

	WorkDir      string
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

	JsonResp = "[]"
	Port     int

	ResLoading = false
	Res        = map[string]map[string][]string{}

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
