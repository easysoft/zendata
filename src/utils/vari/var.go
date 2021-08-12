package vari

import (
	"database/sql"
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/const"
)

var (
	Config = model.Config{Version: 1, Language: "en"}
	DB     *sql.DB

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

	Total     int
	WithHead  bool
	Human     bool
	Trim      bool
	Recursive bool
	Format    = constant.FormatText

	Out      string
	Table    string
	Type     string
	Server   string // database type
	DBDsn    string
	DBClear  bool
	ProtoCls string

	JsonResp string = "[]"
	Ip       string
	Port     int

	ResLoading               = false
	Def                      = model.DefData{}
	Res                      = map[string]map[string][]string{}
	RandFieldNameToValuesMap = map[string]map[string]interface{}{}
	TopFieldMap              = map[string]model.DefField{}

	CacheResFileToMap  = map[string]map[string][]string{}
	CacheResFileToName = map[string]string{}

	DefaultFileDir string
	ConfigFileDir  string

	AgentLogDir string
)
