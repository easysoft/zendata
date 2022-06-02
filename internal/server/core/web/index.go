package web

import (
	stdContext "context"
	"fmt"
	"github.com/easysoft/zendata/internal/server"
	"github.com/easysoft/zendata/internal/server/core/module"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/facebookgo/inject"
	"path/filepath"
	"sync"
	"time"

	"github.com/kataras/iris/v12/context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/snowlyg/helper/dir"
)

type WebServer struct {
	app               *iris.Application
	modules           []module.WebModule
	idleConnsClosed   chan struct{}
	addr              string
	timeFormat        string
	globalMiddlewares []context.Handler
	wg                sync.WaitGroup
	staticPrefix      string
	staticPath        string
	webPath           string
}

// Init 初始化web服务
func Init(port int) *WebServer {
	app := iris.New()
	app.Logger().SetLevel("debug")
	idleConnClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		app.Shutdown(ctx) // close all hosts
		close(idleConnClosed)
	})

	// init grpc
	mvc.New(app)

	// init http
	addr := ":8085"
	if port != 0 {
		addr = fmt.Sprintf(":%d", port)
	}

	webServer := &WebServer{
		app:               app,
		addr:              addr,
		timeFormat:        "2006-01-02 15:04:05",
		staticPrefix:      "/upload",
		staticPath:        "/static/upload",
		webPath:           "./static/dist",
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
		logUtils.PrintErrMsg(fmt.Sprintf("provide usecase objects to the Graph: %v", err))
	}
	err := g.Populate()
	if err != nil {
		logUtils.PrintErrMsg(fmt.Sprintf("populate the incomplete Objects: %v", err))
	}

	ws.AddModule(indexModule.Party())

	logUtils.PrintTo("start server")
}

// GetStaticPath 获取静态路径
func (webServer *WebServer) GetStaticPath() string {
	return webServer.staticPath
}

// GetWebPath 获取前端路径
func (webServer *WebServer) GetWebPath() string {
	return webServer.webPath
}

// GetAddr 获取web服务地址
func (webServer *WebServer) GetAddr() string {
	return webServer.addr
}

// AddModule 添加模块
func (webServer *WebServer) AddModule(module ...module.WebModule) {
	webServer.modules = append(webServer.modules, module...)
}

// AddStatic 添加静态文件
func (webServer *WebServer) AddStatic(requestPath string, fsOrDir interface{}, opts ...iris.DirOptions) {
	webServer.app.HandleDir(requestPath, fsOrDir, opts...)
}

// AddWebStatic 添加前端访问地址
func (webServer *WebServer) AddWebStatic(requestPath string) {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), webServer.webPath))

	webServer.AddStatic(requestPath, fsOrDir, iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	})
}

// AddUploadStatic 添加上传文件访问地址
func (webServer *WebServer) AddUploadStatic() {
	fsOrDir := iris.Dir(filepath.Join(dir.GetCurrentAbPath(), webServer.staticPath))
	webServer.AddStatic(webServer.staticPrefix, fsOrDir)
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
		fmt.Printf("初始化路由错误： %v\n", err)
		panic(err)
	}

	//logUtils.Info(i118Utils.Sprintf("start_server", "localhost",
	//	strings.Replace(webServer.addr, ":", "", -1)))

	webServer.app.Listen(
		webServer.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(webServer.timeFormat),
	)
	<-webServer.idleConnsClosed
}
