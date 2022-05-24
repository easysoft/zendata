package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/action"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/server"
	serverConfig "github.com/easysoft/zendata/internal/server/config"
	serverUtils "github.com/easysoft/zendata/internal/server/utils"
	serverConst "github.com/easysoft/zendata/internal/server/utils/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/easysoft/zendata/res"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/fatih/color"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	flagSet.StringVar(&vari.Ip, "b", "", "")
	flagSet.StringVar(&vari.Ip, "bind", "", "")
	flagSet.IntVar(&vari.Port, "p", 0, "")
	flagSet.IntVar(&vari.Port, "port", 0, "")

	configUtils.InitConfig(root)
	vari.DB, _ = serverConfig.NewGormDB()

	vari.AgentLogDir = vari.ZdPath + serverConst.AgentLogDir + constant.PthSep
	err := fileUtils.MkDirIfNeeded(vari.AgentLogDir)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("perm_deny", vari.AgentLogDir), color.FgRed)
		os.Exit(1)
	}

	startServer()
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
		Handler: handler(server),
	}

	httpServer.ListenAndServe()
}

func handler(s *server.Server) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer( // client static
		&assetfs.AssetFS{Asset: res.Asset, AssetDir: res.AssetDir, AssetInfo: res.AssetInfo, Prefix: "ui/dist"}))
	mux.HandleFunc("/admin", s.Admin)    // data admin
	mux.HandleFunc("/data", dataHandler) // data gen

	return mux
}

func dataHandler(writer http.ResponseWriter, req *http.Request) {
	logUtils.HttpWriter = writer

	if req.Method == http.MethodGet {
		defaultFile, configFile, fields, vari.Total,
			vari.Format, vari.Trim, vari.Table, decode, input, vari.Out = serverUtils.ParseGenParams(req)
	} else if req.Method == http.MethodPost {
		defaultDefContent, configDefContent, fields, vari.Total,
			vari.Format, vari.Trim, vari.Table, decode, input, vari.Out = serverUtils.ParseGenParamsToByte(req)
	}

	if decode {
		files := []string{defaultFile, configFile}
		gen.Decode(files, fields, input)

	} else if defaultDefContent != nil || configDefContent != nil {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintToWithoutNewLine(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))

		genData()
		// Avoid variable affecting the results of request.
		defaultDefContent = nil
		configDefContent = nil

	} else if defaultFile != "" || configFile != "" {
		vari.RunMode = constant.RunModeServerRequest
		logUtils.PrintToWithoutNewLine(i118Utils.I118Prt.Sprintf("server_request", req.Method, req.URL))

		genData()
		// Avoid variable affecting the results of request.
		defaultFile = ""
		configFile = ""
	}
}

func genData() {
	tmStart := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("Start at %s.", tmStart.Format("2006-01-02 15:04:05")))
	}

	vari.Format = constant.FormatJson
	if defaultFile != "" || configFile != "" {
		files := []string{defaultFile, configFile}
		action.Generate(files, fields, vari.Format, vari.Table)
	} else {
		contents := [][]byte{defaultDefContent, configDefContent}
		action.GenerateByContent(contents, fields, vari.Format, vari.Table)
	}

	tmEnd := time.Now()
	if vari.Verbose {
		logUtils.PrintTo(fmt.Sprintf("End at %s.", tmEnd.Format("2006-01-02 15:04:05")))

		dur := tmEnd.Unix() - tmStart.Unix()
		logUtils.PrintTo(fmt.Sprintf("Duriation %d sec.", dur))
	}
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
