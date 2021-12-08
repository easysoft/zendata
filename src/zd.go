package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zendata/res"
	"github.com/easysoft/zendata/src/action"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/server"
	serverConfig "github.com/easysoft/zendata/src/server/config"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	serverConst "github.com/easysoft/zendata/src/server/utils/const"
	"github.com/easysoft/zendata/src/service"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	configUtils "github.com/easysoft/zendata/src/utils/config"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	configs     []string
	defaultFile string
	configFile  string
	// content
	defaultDefContent []byte
	configDefContent  []byte
	//count       int
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
	flagSet.StringVar(&vari.DBDsn, "dns", "", "")
	flagSet.BoolVar(&vari.DBClear, "clear", false, "")

	flagSet.StringVar(&vari.ProtoCls, "cls", "", "")

	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	if len(os.Args) == 1 {
		os.Args = append(os.Args, "-help")
	}

	files, count := fileUtils.GetFilesFromParams(os.Args[1:])
	flagSet.Parse(os.Args[1+count:])
	if count == 0 {
		files = []string{defaultFile, configFile}
	}

	if vari.Ip != "" || vari.Port != 0 {
		vari.RunMode = constant.RunModeServer
	}

	configUtils.InitConfig(root)
	vari.DB, _ = serverConfig.NewGormDB()
	defer vari.DB.Close()

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
			} else if listData {
				service.ListData()
				return
			} else if listRes {
				service.ListRes()
				return
			} else if view != "" {
				service.View(view)
				return
			} else if md5 != "" {
				service.AddMd5(md5)
				return
			} else if decode {
				gen.Decode(files, fields, input)
				return
			}

			if vari.Ip != "" || vari.Port != 0 {
				vari.RunMode = constant.RunModeServer
			} else if input != "" {
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
		//  use the files from post data
		// todo 直接从流中读取，不生成临时文件
		vari.Format = constant.FormatJson
		//action.Generate(files, fields, vari.Format, vari.Table)
		action.Generate2(files, fields, vari.Format, vari.Table)

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

		action.Generate(files, fields, vari.Format, vari.Table)
	}

	tmEnd := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("End at %s.", tmEnd.Format("2006-01-02 15:04:05")))

		dur := tmEnd.Unix() - tmStart.Unix()
		logUtils.PrintTo(fmt.Sprintf("Duriation %d sec.", dur))
	}
}

func startServer() {
	if vari.Ip == "" {
		vari.Ip = commonUtils.GetIp()
	}
	if vari.Port == 0 {
		vari.Port = constant.DefaultPort
	}

	port := strconv.Itoa(vari.Port)
	logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server",
		vari.Ip, port, vari.Ip, port, vari.Ip, port), color.FgCyan)

	// start admin server
	config := serverConfig.NewConfig()
	server, err := server.InitServer(config)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server_fail", port), color.FgRed)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", server.Config.ServerPort),
		Handler: Handler(server),
	}

	httpServer.ListenAndServe()
}

func Handler(s *server.Server) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer( // client static
		&assetfs.AssetFS{Asset: res.Asset, AssetDir: res.AssetDir, AssetInfo: res.AssetInfo, Prefix: "ui/dist"}))
	mux.HandleFunc("/admin", s.Admin)    // data admin
	mux.HandleFunc("/data", DataHandler) // data gen

	return mux
}

func DataHandler(writer http.ResponseWriter, req *http.Request) {
	logUtils.HttpWriter = writer

	defaultDefContent, configDefContent, fields, vari.Total,
		vari.Format, vari.Trim, vari.Table, decode, input, vari.Out = serverUtils.ParseGenParamsToByte(req)

	if decode {
		files := []string{defaultFile, configFile}
		gen.Decode(files, fields, input)
	} else if defaultDefContent != nil || configDefContent != nil {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintToWithoutNewLine(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))
		files := []string{string(defaultDefContent), string(configDefContent)}

		toGen(files)
	}
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
