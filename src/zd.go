package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zendata/src/action"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/service"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	configUtils "github.com/easysoft/zendata/src/utils/config"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
)

var (
	defaultFile string
	configFile  string
	count       int
	fields      string

	root string
	input  string
	output string
	table  string
	format = constant.FormatText
	decode bool

	listRes bool
	viewRes string
	viewDetail string
	md5 string

	example bool
	help   bool

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

	flagSet.IntVar(&count, "n", 10, "")
	flagSet.IntVar(&count, "lines", 10, "")

	flagSet.StringVar(&fields, "F", "", "")
	flagSet.StringVar(&fields, "field", "", "")

	flagSet.StringVar(&output, "o", "", "")
	flagSet.StringVar(&output, "output", "", "")

	flagSet.StringVar(&table, "t", "table_name", "")
	flagSet.StringVar(&table, "table", "table_name", "")

	flagSet.BoolVar(&listRes, "l", false, "")
	flagSet.BoolVar(&listRes, "list", false, "")

	flagSet.StringVar(&viewRes, "v", "", "")
	flagSet.StringVar(&viewRes, "view", "", "")

	flagSet.StringVar(&md5, "md5", "", "")

	flagSet.StringVar(&vari.HeadSep, "H", "", "")
	flagSet.StringVar(&vari.HeadSep, "human", "", "")

	flagSet.BoolVar(&decode, "D", false, "")
	flagSet.BoolVar(&decode, "decode", false, "")

	flagSet.StringVar(&vari.Ip, "b", "", "")
	flagSet.StringVar(&vari.Ip, "bind", "", "")
	flagSet.IntVar(&vari.Port, "p", 0, "")
	flagSet.IntVar(&vari.Port, "port", 0, "")
	flagSet.StringVar(&root, "R", "", "")
	flagSet.StringVar(&root, "root", "", "")

	flagSet.BoolVar(&example, "e", false, "")
	flagSet.BoolVar(&example, "example", false, "")

	flagSet.BoolVar(&help, "h", false, "")
	flagSet.BoolVar(&help, "help", false, "")

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-help")
	}

	switch os.Args[1] {
	case "-e", "-example":
		logUtils.PrintExample()
	case "-h", "-help":
		logUtils.PrintUsage()
	default:
		flagSet.SetOutput(ioutil.Discard)
		if err := flagSet.Parse(os.Args[1:]); err == nil {
			if listRes {
				service.ListRes()
				return
			} else if viewRes != "" {
				service.ViewRes(viewRes)
				return
			} else if md5 != "" {
				service.AddMd5(md5)
				return
			} else if decode {
				gen.Decode(defaultFile, configFile, fields, input, output)
				return
			}

			if vari.Ip != "" || vari.Port != 0 {
				vari.RunMode = constant.RunModeServer
			} else if input != "" {
				vari.RunMode = constant.RunModeParse
			}

			toGen()
		} else {
			logUtils.PrintUsage()
		}
	}
}

func toGen() {
	if vari.RunMode == constant.RunModeServer {
		StartServer()
	} else if vari.RunMode == constant.RunModeParse {
		action.ParseSql(input, output)
	} else if vari.RunMode == constant.RunModeGen {
		if root != "" {
			vari.ExeDir = root
		}
		if vari.HeadSep != "" {
			vari.WithHead = true
		}

		if output != "" {
			ext := strings.ToLower(path.Ext(output))
			if len(ext) > 1 {
				ext = strings.TrimLeft(ext,".")
			}
			if stringUtils.InArray(ext, constant.Formats) {
				format = ext
			}
		}

		action.Generate(defaultFile, configFile, count, fields, output, format, table)
	} else if vari.RunMode == constant.RunModeServerRequest {
		action.Generate(defaultFile, configFile, count, fields, output, format, table)
	}
}

func StartServer() {
	if vari.Ip == "" {
		vari.Ip = commonUtils.GetIp()
	}
	if vari.Port == 0 {
		vari.Port = constant.DefaultPort
	}

	port := strconv.Itoa(vari.Port)
	logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server", vari.Ip, port, vari.Ip, port),
		color.FgCyan)

	http.HandleFunc("/", DataHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", vari.Port), nil)
}

func DataHandler(w http.ResponseWriter, req *http.Request) {


	root, defaultFile, configFile, fields, count, vari.HeadSep,
		format, table, decode, input, output = service.ParseRequestParams(req)

	if decode {
		gen.Decode(defaultFile, configFile, fields, input, output)
		fmt.Fprintln(w, vari.JsonResp)
	} else if defaultFile != "" || configFile != "" {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))
		toGen()
		fmt.Fprintln(w, vari.JsonResp)
	}
}

func init() {
	cleanup()

	logUtils.InitLogger()
	configUtils.InitConfig()
}

func cleanup() {
	color.Unset()
}
