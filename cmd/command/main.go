package main

import (
	"flag"
	"github.com/easysoft/zendata/internal/command"
	commandConfig "github.com/easysoft/zendata/internal/command/config"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/helper"
	serverConfig "github.com/easysoft/zendata/internal/server/config"
	"github.com/easysoft/zendata/internal/server/core/web"
	serverConst "github.com/easysoft/zendata/internal/server/utils/const"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
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

	root  string
	input string

	parse    bool
	decode   bool
	listData bool
	listRes  bool
	view     string
	md5      string
	salt     string

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

	flagSet.StringVar(&defaultFile, "d", "", "")
	flagSet.StringVar(&defaultFile, "default", "", "")

	flagSet.StringVar(&configFile, "c", "", "")
	flagSet.StringVar(&configFile, "config", "", "")

	flagSet.StringVar(&input, "i", "", "")
	flagSet.StringVar(&input, "input", "", "")

	flagSet.IntVar(&vari.GlobalVars.Total, "n", -1, "")
	flagSet.IntVar(&vari.GlobalVars.Total, "lines", -1, "")

	flagSet.StringVar(&exportFields, "F", "", "")
	flagSet.StringVar(&exportFields, "field", "", "")

	flagSet.StringVar(&vari.GlobalVars.OutputFormat, "f", consts.FormatText, "")
	flagSet.StringVar(&vari.GlobalVars.OutputFormat, "format", consts.FormatText, "")

	flagSet.StringVar(&vari.GlobalVars.OutputFile, "o", "", "")
	flagSet.StringVar(&vari.GlobalVars.OutputFile, "output", "", "")

	flagSet.BoolVar(&parse, "parse", false, "")

	flagSet.BoolVar(&listData, "l", false, "")
	flagSet.BoolVar(&listData, "list", false, "")
	flagSet.BoolVar(&listRes, "L", false, "")

	flagSet.StringVar(&view, "v", "", "")
	flagSet.StringVar(&view, "view", "", "")

	flagSet.StringVar(&md5, "md5", "", "")
	flagSet.StringVar(&salt, "salt", "", "")

	flagSet.BoolVar(&decode, "D", false, "")
	flagSet.BoolVar(&decode, "decode", false, "")

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

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	flagSet.BoolVar(&help, "h", false, "")
	flagSet.BoolVar(&help, "help", false, "")

	// for server
	flagSet.BoolVar(&isStartServer, "s", false, "启动服务")
	flagSet.StringVar(&uuid, "uuid", "", "区分服务进程的唯一ID")

	flagSet.IntVar(&vari.DataServicePort, "p", 8848, "")
	flagSet.IntVar(&vari.DataServicePort, "port", 0, "")

	flagSet.Parse(os.Args[1:])
	if isStartServer {
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

func startServer() {
	configUtils.InitConfig("")
	vari.DB, _ = serverConfig.NewGormDB()

	vari.AgentLogDir = vari.ZdPath + serverConst.AgentLogDir + consts.PthSep
	err := fileUtils.MkDirIfNeeded(vari.AgentLogDir)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("perm_deny", vari.AgentLogDir), color.FgRed)
		os.Exit(1)
	}

	if vari.DataServicePort == 0 {
		vari.DataServicePort = consts.DefaultDataServicePort
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
	}

	genData(files)
}

func genYaml(input string) {
	mainCtrl, _ := command.InitCtrl()
	mainCtrl.Parse(input)
}

func genData(files []string) {
	command.PrintStartInfo()

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
