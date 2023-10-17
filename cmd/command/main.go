package main

import (
	"flag"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/easysoft/zendata/internal/command"
	commandConfig "github.com/easysoft/zendata/internal/command/config"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/helper"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
)

var (
	defaultFile  string
	configFile   string
	exportFields string

	root  string
	input string

	parse    bool
	listData bool
	listRes  bool
	view     string
	md5      string
	salt     string
	mock     bool

	example bool
	help    bool
	version bool
	set     bool

	AppVersion string
	BuildTime  string
	GoVersion  string
	GitHash    string

	flagSet *flag.FlagSet
)

func main() {
	channel := make(chan os.Signal, 1) // fix for sigchanyzer: misuse of unbuffered os.Signal channel as argument to signal.Notify
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

	flagSet.IntVar(&vari.GlobalVars.Total, "n", -1, "")
	flagSet.IntVar(&vari.GlobalVars.Total, "lines", -1, "")

	flagSet.StringVar(&exportFields, "F", "", "")
	flagSet.StringVar(&exportFields, "field", "", "")

	flagSet.StringVar(&vari.GlobalVars.OutputFormat, "f", consts.FormatText, "")
	flagSet.StringVar(&vari.GlobalVars.OutputFormat, "format", consts.FormatText, "")

	flagSet.StringVar(&vari.GlobalVars.Output, "o", "", "")
	flagSet.StringVar(&vari.GlobalVars.Output, "output", "", "")

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

	// gen yaml from sql or table
	flagSet.BoolVar(&parse, "parse", false, "")

	flagSet.StringVar(&input, "i", "", "")
	flagSet.StringVar(&input, "input", "", "")

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
	flagSet.BoolVar(&version, "version", false, "")

	flagSet.Parse(os.Args[1:])

	execCommand()
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

	flagSet.SetOutput(io.Discard)
	if err := flagSet.Parse(os.Args[1:]); err == nil {
		opts(files)
	} else {
		logUtils.PrintUsage()
	}
}

func opts(files []string) {
	if exportFields != "" {
		vari.GlobalVars.ExportFields = strings.Split(exportFields, ",")
	}

	if version {
		logUtils.PrintVersion(AppVersion, BuildTime, GoVersion, GitHash)
		return

	} else if example {
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

func genMock(input string) {
	mainCtrl, _ := command.InitCtrl()
	mainCtrl.GenMock(input)
}

func genData(files []string) {
	command.PrintStartInfo()

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
