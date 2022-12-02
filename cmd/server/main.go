package main

import (
	"flag"
	"fmt"
	"github.com/easysoft/zendata/internal/agent"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/server"
	serverConfig "github.com/easysoft/zendata/internal/server/config"
	"github.com/easysoft/zendata/internal/server/core/web"
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
)

var (
	flagSet *flag.FlagSet
	uuid    = ""
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

	flagSet.StringVar(&vari.Ip, "b", "", "")
	flagSet.StringVar(&vari.Ip, "bind", "", "")
	flagSet.IntVar(&vari.Port, "p", 0, "")
	flagSet.IntVar(&vari.Port, "port", 0, "")
	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	configUtils.InitConfig("")
	vari.DB, _ = serverConfig.NewGormDB()

	vari.AgentLogDir = vari.ZdPath + serverConst.AgentLogDir + constant.PthSep
	err := fileUtils.MkDirIfNeeded(vari.AgentLogDir)
	if err != nil {
		logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("perm_deny", vari.AgentLogDir), color.FgRed)
		os.Exit(1)
	}

	go func() {
		startDataServer()
	}()

	startAdminServer()
}

func startAdminServer() {
	webServer := web.Init(vari.Port)
	if webServer == nil {
		return
	}

	webServer.Run()
}

func startDataServer() {
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

func handler(server *server.Server) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(
		&assetfs.AssetFS{Asset: res.Asset, AssetDir: res.AssetDir, AssetInfo: res.AssetInfo, Prefix: "ui/dist"}))

	//mux.HandleFunc("/admin", server.AdminHandler)
	mux.HandleFunc("/data", agent.DataHandler)

	return mux
}

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
