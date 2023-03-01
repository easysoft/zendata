package web

import (
	stdContext "context"
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/server"
	"github.com/easysoft/zendata/internal/server/core/module"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/facebookgo/inject"
	"github.com/fatih/color"
	"strconv"
	"sync"
	"time"

	"github.com/kataras/iris/v12/context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type WebServer struct {
	app               *iris.Application
	modules           []module.WebModule
	idleConnsClosed   chan struct{}
	addr              string
	timeFormat        string
	globalMiddlewares []context.Handler
	wg                sync.WaitGroup
}

// Init 初始化web服务
func Init() *WebServer {
	app := iris.New()

	level := "info"
	if vari.Verbose {
		level = "debug"
	}
	app.Logger().SetLevel(level)

	idleConnClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnClosed)
	})

	mvc.New(app)

	addr := fmt.Sprintf(":%d", vari.Port)

	webServer := &WebServer{
		app:               app,
		addr:              addr,
		timeFormat:        "2006-01-02 15:04:05",
		idleConnsClosed:   idleConnClosed,
		globalMiddlewares: []context.Handler{},
	}

	injectModule(webServer)

	return webServer
}

func injectModule(ws *WebServer) {
	var g inject.Graph

	indexModule := server.NewIndexModule()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: vari.DB},
		&inject.Object{Value: indexModule},
	); err != nil {
		panic(fmt.Sprintf("provide usecase objects to the Graph: %v", err))
	}
	err := g.Populate()
	if err != nil {
		panic(fmt.Sprintf("populate the incomplete Objects: %v", err))
	}

	ws.AddModule(indexModule.Party())
	ws.AddModule(indexModule.PartyData())
	ws.AddModule(indexModule.PartyMock())

	logUtils.PrintTo("start server")
}

// GetAddr 获取web服务地址
func (webServer *WebServer) GetAddr() string {
	return webServer.addr
}

// AddModule 添加模块
func (webServer *WebServer) AddModule(module ...module.WebModule) {
	webServer.modules = append(webServer.modules, module...)
}

// GetModules 获取模块
func (webServer *WebServer) GetModules() []module.WebModule {
	return webServer.modules
}

// Run 启动web服务
func (webServer *WebServer) Run() {
	webServer.app.UseGlobal(webServer.globalMiddlewares...)
	err := webServer.InitRouter()
	if err != nil {
		panic(fmt.Sprintf("初始化路由错误： %v", err))
	}

	port := strconv.Itoa(vari.Port)
	logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("start_server",
		consts.Localhost, port, consts.Localhost, port), color.FgCyan)

	webServer.app.Listen(
		webServer.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(webServer.timeFormat),
	)

	<-webServer.idleConnsClosed
}
