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
	fileUtils "github.com/easysoft/zendata/src/utils/file"
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
	"path/filepath"
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
	set   bool

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

	flagSet.StringVar(&table, "t", "", "")
	flagSet.StringVar(&table, "table", "", "")

	flagSet.BoolVar(&listRes, "l", false, "")
	flagSet.BoolVar(&listRes, "list", false, "")

	flagSet.StringVar(&viewRes, "v", "", "")
	flagSet.StringVar(&viewRes, "view", "", "")

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

	flagSet.BoolVar(&example, "e", false, "")
	flagSet.BoolVar(&example, "example", false, "")

	flagSet.BoolVar(&help, "h", false, "")
	flagSet.BoolVar(&help, "help", false, "")

	flagSet.BoolVar(&set, "s", false, "")
    flagSet.BoolVar(&set, "set", false, "")

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-help")
	}

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
                service.Set()
                return
            } else if listRes {
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
		if root != "" {
			if fileUtils.IsAbosutePath(root) {
				vari.ExeDir = root
			} else {
				vari.ExeDir = vari.ExeDir + root
			}
			vari.ExeDir = fileUtils.AddSepIfNeeded(vari.ExeDir)
		}
		constant.SqliteSource = strings.Replace(constant.SqliteSource, "file:", "file:" + vari.ExeDir, 1)

		StartServer()
	} else if vari.RunMode == constant.RunModeServerRequest {
		format = constant.FormatJson
		action.Generate(defaultFile, configFile, count, fields, output, format, table)
	} else if vari.RunMode == constant.RunModeParse {
		action.ParseSql(input, output)
	} else if vari.RunMode == constant.RunModeGen {
		if vari.Human {
			vari.WithHead = true
		}

		if output != "" {
			fileUtils.MkDirIfNeeded(filepath.Dir(output))
			fileUtils.RemoveExist(output)
			action.FileWriter, _ = os.OpenFile(output, os.O_RDWR | os.O_CREATE, 0777)
			defer action.FileWriter.Close()

			ext := strings.ToLower(path.Ext(output))
			if len(ext) > 1 {
				ext = strings.TrimLeft(ext,".")
			}
			if stringUtils.InArray(ext, constant.Formats) {
				format = ext
			}
		}

		if format == constant.FormatSql && table == "" {
			action.PrintErrMsg(i118Utils.I118Prt.Sprintf("miss_table_name"))
			return
		}

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

	http.HandleFunc("/", DataHandler)

	port := strconv.Itoa(vari.Port)

	logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server", vari.Ip, port, vari.Ip, port), color.FgCyan)
	err := http.ListenAndServe(fmt.Sprintf(":%d", vari.Port), nil)

	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server_fail", port), color.FgRed)
	}
}

func DataHandler(writer http.ResponseWriter, req *http.Request) {
	action.HttpWriter = writer

	defaultFile, configFile, fields, count,
		format, table, decode, input, output = service.ParseRequestParams(req)

	if decode {
		gen.Decode(defaultFile, configFile, fields, input, output)
		fmt.Fprintln(writer, vari.JsonResp)
	} else if defaultFile != "" || configFile != "" {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))

		toGen()
	}
}

func init() {
	cleanup()
	configUtils.InitConfig()
}

func cleanup() {
	color.Unset()
}
