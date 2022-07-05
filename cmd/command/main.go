package main

import (
	"errors"
	"flag"
	"fmt"
	commandConfig "github.com/easysoft/zendata/internal/command/config"
	"github.com/easysoft/zendata/internal/pkg/action"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/helper"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	configs     []string
	defaultFile string
	configFile  string

	defaultDefContent []byte
	configDefContent  []byte

	fields string

	root   string
	input  string
	decode bool

	listData bool
	listRes  bool
	view     string
	md5      string

	example bool
	help    bool
	set     bool

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
	flagSet.StringVar(&configFile, "Config", "", "")

	flagSet.StringVar(&input, "i", "", "")
	flagSet.StringVar(&input, "input", "", "")

	flagSet.IntVar(&vari.Total, "n", -1, "")
	flagSet.IntVar(&vari.Total, "lines", -1, "")

	flagSet.StringVar(&fields, "F", "", "")
	flagSet.StringVar(&fields, "field", "", "")

	flagSet.StringVar(&vari.Out, "o", "", "")
	flagSet.StringVar(&vari.Out, "output", "", "")

	flagSet.BoolVar(&listData, "l", false, "")
	flagSet.BoolVar(&listData, "list", false, "")
	flagSet.BoolVar(&listRes, "L", false, "")

	flagSet.StringVar(&view, "v", "", "")
	flagSet.StringVar(&view, "view", "", "")

	flagSet.StringVar(&md5, "md5", "", "")

	flagSet.BoolVar(&vari.Human, "H", false, "")
	flagSet.BoolVar(&vari.Human, "human", false, "")

	flagSet.BoolVar(&decode, "D", false, "")
	flagSet.BoolVar(&decode, "decode", false, "")

	flagSet.StringVar(&vari.Ip, "b", "", "")
	flagSet.StringVar(&vari.Ip, "bind", "", "")
	flagSet.IntVar(&vari.Port, "p", 0, "")
	flagSet.IntVar(&vari.Port, "port", 0, "")
	flagSet.StringVar(&root, "R", "", "")
	flagSet.StringVar(&root, "root", "", "")

	flagSet.BoolVar(&vari.Trim, "T", false, "")
	flagSet.BoolVar(&vari.Trim, "trim", false, "")

	flagSet.BoolVar(&vari.Recursive, "r", false, "")
	flagSet.BoolVar(&vari.Recursive, "recursive", false, "")

	flagSet.StringVar(&vari.CacheParam, "C", "", "")
	flagSet.StringVar(&vari.CacheParam, "cache", "", "")

	flagSet.BoolVar(&example, "e", false, "")
	flagSet.BoolVar(&example, "example", false, "")

	flagSet.BoolVar(&help, "h", false, "")
	flagSet.BoolVar(&help, "help", false, "")

	flagSet.BoolVar(&set, "S", false, "")
	flagSet.BoolVar(&set, "set", false, "")

	flagSet.StringVar(&vari.Table, "t", "", "")
	flagSet.StringVar(&vari.Table, "table", "", "")
	flagSet.StringVar(&vari.Server, "s", "mysql", "")
	flagSet.StringVar(&vari.Server, "server", "mysql", "")
	flagSet.StringVar(&vari.DBDsn, "dsn", "", "")
	flagSet.BoolVar(&vari.DBClear, "clear", false, "")

	flagSet.StringVar(&vari.ProtoCls, "cls", "", "")

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

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

	switch os.Args[1] {
	default:
		flagSet.SetOutput(ioutil.Discard)
		if err := flagSet.Parse(os.Args[1:]); err == nil {
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
				helper.AddMd5(md5)
				return
			} else if decode {
				gen.Decode(files, fields, input)
				return
			}

			if input != "" {
				vari.RunMode = constant.RunModeParse
			}

			toGen(files)
		} else {
			logUtils.PrintUsage()
		}
	}
}

func toGen(files []string) {
	tmStart := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("Start at %s.", tmStart.Format("2006-01-02 15:04:05")))
	}

	if vari.RunMode == constant.RunModeGen {
		if clearCache() {
			return
		}

		if err := getFormat(); err != nil {
			return
		}
		action.Generate(files, fields, vari.Format, vari.Table)

	} else if vari.RunMode == constant.RunModeParse {
		ext := filepath.Ext(input)
		if ext == ".sql" {
			action.ParseSql(input, vari.Out)
		} else if ext == ".txt" {
			action.ParseArticle(input, vari.Out)
		}
	}

	tmEnd := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("End at %s.", tmEnd.Format("2006-01-02 15:04:05")))
		dur := tmEnd.Unix() - tmStart.Unix()
		logUtils.PrintTo(fmt.Sprintf("Duriation %d sec.", dur))
	}
}

func clearCache() (ret bool) {
	cacheKey, cacheOpt, hasCache := gen.ParseCache()
	if cacheOpt == "clear" {
		if cacheKey == "all" {
			gen.ClearAllCache()
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("success_to_clear_all_cache"))
		} else {
			if hasCache {
				gen.ClearCache(cacheKey)
			}
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("success_to_clear_cache", cacheKey))
		}

		ret = true
	}

	return
}

func getFormat() (err error) {
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
		msg := i118Utils.I118Prt.Sprintf("miss_table_name")
		logUtils.PrintErrMsg(msg)
		err = errors.New(msg)
	}

	return
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
