package main

import (
	"encoding/json"
	"flag"
	"github.com/easysoft/zendata/internal/command"
	commandConfig "github.com/easysoft/zendata/internal/command/config"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/helper"
	serverConfig "github.com/easysoft/zendata/internal/server/config"
<<<<<<< HEAD
	serverUtils "github.com/easysoft/zendata/internal/server/utils"
	serverConst "github.com/easysoft/zendata/internal/server/utils/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
=======
	"github.com/easysoft/zendata/internal/server/core/web"
	serverConst "github.com/easysoft/zendata/internal/server/utils/const"
>>>>>>> 3.0
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	configs      []string
	defaultFile  string
	configFile   string
	exportFields string

	defaultDefContent []byte
	configDefContent  []byte

<<<<<<< HEAD
	uuid   = ""
	root   string
	input  string
	decode bool
=======
	root  string
	input string
>>>>>>> 3.0

	parse    bool
	decode   bool
	listData bool
	listRes  bool
	view     string
	md5      string
	salt     string
	mock     bool

	example bool
	help    bool
	set     bool

	isStartServer bool
	uuid          = ""

	flagSet *flag.FlagSet
)

func main() {
	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-channel
		cleanup()
		os.Exit(0)
	}()

	flagSet = flag.NewFlagSet("zd", flag.ContinueOnError)

	flagSet.StringVar(&uuid, "uuid", "", "区分服务进程的唯一ID")

	flagSet.StringVar(&defaultFile, "d", "", "")
	flagSet.StringVar(&defaultFile, "default", "", "")

	flagSet.StringVar(&configFile, "c", "", "")
	flagSet.StringVar(&configFile, "config", "", "")

	flagSet.IntVar(&vari.GlobalVars.Total, "n", -1, "")
	flagSet.IntVar(&vari.GlobalVars.Total, "lines", -1, "")

	flagSet.StringVar(&exportFields, "F", "", "")
	flagSet.StringVar(&exportFields, "field", "", "")

	flagSet.StringVar(&vari.GlobalVars.OutputFormat, "f", consts.FormatText, "")
	flagSet.StringVar(&vari.GlobalVars.OutputFormat, "format", consts.FormatText, "")

	flagSet.StringVar(&input, "i", "", "")
	flagSet.StringVar(&input, "input", "", "")

	flagSet.StringVar(&vari.GlobalVars.Output, "o", "", "")
	flagSet.StringVar(&vari.GlobalVars.Output, "output", "", "")

	flagSet.BoolVar(&parse, "parse", false, "")

	flagSet.BoolVar(&listData, "l", false, "")
	flagSet.BoolVar(&listData, "list", false, "")
	flagSet.BoolVar(&listRes, "L", false, "")

	flagSet.StringVar(&view, "v", "", "")
	flagSet.StringVar(&view, "view", "", "")

	flagSet.StringVar(&md5, "md5", "", "")
	flagSet.StringVar(&salt, "salt", "", "")

	flagSet.BoolVar(&mock, "m", false, "")

	flagSet.StringVar(&root, "R", "", "")
	flagSet.StringVar(&root, "root", "", "")

	flagSet.BoolVar(&vari.GlobalVars.Human, "H", false, "")
	flagSet.BoolVar(&vari.GlobalVars.Human, "human", false, "")

	flagSet.BoolVar(&vari.GlobalVars.Trim, "T", false, "")
	flagSet.BoolVar(&vari.GlobalVars.Trim, "trim", false, "")

	flagSet.BoolVar(&vari.GlobalVars.Recursive, "r", false, "")
	flagSet.BoolVar(&vari.GlobalVars.Recursive, "recursive", false, "")

	//flagSet.StringVar(&vari.CacheParam, "C", "", "")
	//flagSet.StringVar(&vari.CacheParam, "cache", "", "")

	flagSet.BoolVar(&example, "e", false, "")
	flagSet.BoolVar(&example, "example", false, "")

	flagSet.BoolVar(&set, "S", false, "")
	flagSet.BoolVar(&set, "set", false, "")

	flagSet.StringVar(&vari.GlobalVars.Table, "t", "", "")
	flagSet.StringVar(&vari.GlobalVars.Table, "table", "", "")
	flagSet.StringVar(&vari.GlobalVars.DBDsn, "dsn", "", "")
	flagSet.StringVar(&vari.GlobalVars.DBType, "db", "db", "")
	flagSet.StringVar(&vari.GlobalVars.DBType, "server", "mysql", "") // TODO: will remove
	flagSet.BoolVar(&vari.GlobalVars.DBClear, "clear", false, "")

	flagSet.StringVar(&vari.ProtoCls, "cls", "", "")
	flagSet.StringVar(&vari.GlobalVars.MockDir, "mock", "", "")

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	flagSet.BoolVar(&help, "h", false, "")
	flagSet.BoolVar(&help, "help", false, "")

	// for server
	flagSet.BoolVar(&isStartServer, "s", false, "启动服务")
	flagSet.StringVar(&uuid, "uuid", "", "区分服务进程的唯一ID")

	flagSet.IntVar(&vari.Port, "p", 8848, "")
	flagSet.IntVar(&vari.Port, "port", 0, "")

	flagSet.Parse(os.Args[1:])
	if isStartServer {
		vari.GlobalVars.RunMode = consts.RunModeServer
		startServer()
	} else {
		execCommand()
	}
}

func execCommand() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-help")
	}

	files, count := fileUtils.GetFilesFromParams(os.Args[1:]) // zd.exe demo\default.yaml demo\test.yaml
	flagSet.Parse(os.Args[1+count:])
	if count == 0 { // not has a def file list, use -d and -c files
		files = []string{defaultFile, configFile}
	}

	configUtils.InitConfig(root)
	vari.DB, _ = commandConfig.NewGormDB()
	//defer vari.DB.Close()

	flagSet.SetOutput(ioutil.Discard)
	if err := flagSet.Parse(os.Args[1:]); err == nil {
		opts(files)
	} else {
		logUtils.PrintUsage()
	}
}

<<<<<<< HEAD
func toGen(files []string) {
	tmStart := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("Start at %s.", tmStart.Format("2006-01-02 15:04:05")))
	}

	if vari.RunMode == constant.RunModeParse {
		ext := filepath.Ext(input)
		if ext == ".sql" {
			action.ParseSql(input, vari.Out)
		} else if ext == ".txt" {
			action.ParseArticle(input, vari.Out)
		}
	} else if vari.RunMode == constant.RunModeServer {
		vari.AgentLogDir = vari.ZdPath + serverConst.AgentLogDir + constant.PthSep
		err := fileUtils.MkDirIfNeeded(vari.AgentLogDir)
		if err != nil {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("perm_deny", vari.AgentLogDir), color.FgRed)
			os.Exit(1)
		}

		startServer() // will init its own db

	} else if vari.RunMode == constant.RunModeServerRequest {
		vari.Format = constant.FormatJson
		if defaultFile != "" || configFile != "" {
			files := []string{defaultFile, configFile}
			action.Generate(files, fields, vari.Format, vari.Table)
		} else {
			contents := [][]byte{defaultDefContent, configDefContent}
			action.GenerateByContent(contents, fields, vari.Format, vari.Table)
		}

	} else if vari.RunMode == constant.RunModeGen {
		if vari.Human {
			vari.WithHead = true
		}

		if vari.Out != "" {
			fileUtils.MkDirIfNeeded(filepath.Dir(vari.Out))
			fileUtils.RemoveExist(vari.Out)

			ext := strings.ToLower(filepath.Ext(vari.Out))
			if len(ext) > 1 {
				ext = strings.TrimLeft(ext, ".")
			}
			if stringUtils.InArray(ext, constant.Formats) {
				vari.Format = ext
			}

			if vari.Format == constant.FormatExcel {
				logUtils.FilePath = vari.Out
			} else {
				logUtils.FileWriter, _ = os.OpenFile(vari.Out, os.O_RDWR|os.O_CREATE, 0777)
				defer logUtils.FileWriter.Close()
			}
		}
		if vari.DBDsn != "" {
			vari.Format = constant.FormatSql
		}

		if vari.Format == constant.FormatSql && vari.Table == "" {
			logUtils.PrintErrMsg(i118Utils.I118Prt.Sprintf("miss_table_name"))
			return
		}
=======
func startServer() {
	configUtils.InitConfig(root)
	vari.DB, _ = serverConfig.NewGormDB()
>>>>>>> 3.0

	vari.AgentLogDir = vari.ZdDir + serverConst.AgentLogDir + consts.PthSep
	err := fileUtils.MkDirIfNeeded(vari.AgentLogDir)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("perm_deny", vari.AgentLogDir), color.FgRed)
		os.Exit(1)
	}

	if vari.Port == 0 {
		vari.Port = consts.DefaultDataServicePort
	}

	webServer := web.Init()
	if webServer == nil {
		return
	}

	webServer.Run()
}

func opts(files []string) {
	if exportFields != "" {
		vari.GlobalVars.ExportFields = strings.Split(exportFields, ",")
	}

	if example {
		logUtils.PrintExample()
		return
	} else if help {
		logUtils.PrintUsage()
		return
	} else if set {
		helper.Set()
		return
	} else if listData {
		helper.ListData()
		return
	} else if listRes {
		helper.ListRes()
		return
	} else if view != "" {
		helper.View(view)
		return
	} else if md5 != "" {
		helper.AddMd5(md5, salt)
		return
	} else if decode {
		gen.Decode(files, input)
		return
	} else if parse {
		genYaml(input)
		return
	} else if mock {
		if input == "" {
			return
		}
		genMock(input)

		return
	}

	genData(files)
}

func genYaml(input string) {
	mainCtrl, _ := command.InitCtrl()
	mainCtrl.GenYaml(input)
}

<<<<<<< HEAD
	mux.HandleFunc("/api/v1/heartbeat", heartbeat) // test

	return mux
}

func heartbeat(writer http.ResponseWriter, req *http.Request) {
	serverUtils.SetupCORS(&writer, req)

	ret := map[string]interface{}{"code": 0, "msg": "ok"}

	bytes, _ := json.Marshal(ret)
	io.WriteString(writer, string(bytes))
}

func DataHandler(writer http.ResponseWriter, req *http.Request) {
	logUtils.HttpWriter = writer
=======
func genMock(input string) {
	mainCtrl, _ := command.InitCtrl()
	mainCtrl.GenMock(input)
}

func genData(files []string) {
	command.PrintStartInfo()
>>>>>>> 3.0

	if command.ClearCache() {
		return
	}

	err := command.SetOutFormat()
	defer logUtils.OutputFileWriter.Close()
	if err != nil {
		return
	}

	mainCtrl, _ := command.InitCtrl()
	mainCtrl.Generate(files)

	command.PrintEndInfo()
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
