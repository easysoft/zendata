package main

import (
	"flag"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	serverConfig "github.com/easysoft/zendata/internal/server/config"
	"github.com/easysoft/zendata/internal/server/core/web"
	serverConst "github.com/easysoft/zendata/internal/server/utils/const"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"os"
	"os/signal"
	"syscall"
)

var (
	flagSet *flag.FlagSet
	uuid    = ""
	root    string
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

	flagSet.IntVar(&vari.Port, "p", 0, "")
	flagSet.IntVar(&vari.Port, "port", 0, "")
	flagSet.StringVar(&root, "R", "", "")
	flagSet.StringVar(&root, "root", "", "")
	flagSet.BoolVar(&vari.Verbose, "verbose", false, "")

	flagSet.Parse(os.Args[1:])

	vari.GlobalVars.RunMode = consts.RunModeServer

	configUtils.InitConfig(root)
	vari.DB, _ = serverConfig.NewGormDB()

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

func init() {
	cleanup()
}

func cleanup() {
	color.Unset()
}
