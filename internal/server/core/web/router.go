package web

import (
	"fmt"
	"github.com/easysoft/zendata/internal/server/core/middleware"
	"github.com/easysoft/zendata/internal/server/core/module"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/snowlyg/helper/arr"
)

// InitRouter 初始化模块路由
func (webServer *WebServer) InitRouter() error {
	webServer.app.UseRouter(middleware.CrsAuth())

	app := webServer.app.Party("/").AllowMethods(iris.MethodOptions)
	{
		app.Use(middleware.InitCheck())

		if vari.Verbose {
			debug := DebugParty()
			app.PartyFunc(debug.RelativePath, debug.Handler)
		}

		webServer.initModule()

		err := webServer.app.Build()
		if err != nil {
			return fmt.Errorf("build router %w", err)
		}

		return nil
	}
}

// GetSources 获取web服务需要认证的权限
func (webServer *WebServer) GetSources() []map[string]string {
	routeLen := len(webServer.app.GetRoutes())
	ch := make(chan map[string]string, routeLen)
	for _, r := range webServer.app.GetRoutes() {
		r := r
		// 去除非接口路径
		handerNames := context.HandlersNames(r.Handlers)
		if !arr.InArrayS([]string{"GET", "POST", "PUT", "DELETE"}, r.Method) ||
			!arr.InArrayS(strings.Split(handerNames, ","), "github.com/snowlyg/multi.(*Verifier).Verify") {
			routeLen--
			continue
		}
		go func(r *router.Route) {
			route := map[string]string{
				"path": r.Path,
				"name": r.Name,
				"act":  r.Method,
			}
			ch <- route
		}(r)
	}

	routes := make([]map[string]string, routeLen)
	for i := 0; i < routeLen; i++ {
		routes[i] = <-ch
	}
	return routes
}

// initModule 初始化web服务模块，包括子模块
func (webServer *WebServer) initModule() {
	if len(webServer.modules) > 0 {
		for _, mod := range webServer.modules {
			mod := mod
			webServer.wg.Add(1)

			sub := webServer.app.PartyFunc(mod.RelativePath, mod.Handler)
			if len(mod.Modules) > 0 {
				for _, subModule := range mod.Modules {
					sub.PartyFunc(subModule.RelativePath, subModule.Handler)
				}
			}
			webServer.wg.Done()
		}
		webServer.wg.Wait()
	}
}

// Party 调试模块
func DebugParty() module.WebModule {
	handler := func(index iris.Party) {
		index.Get("/", func(ctx iris.Context) {
			ctx.HTML("<h1>请点击<a href='/debug/pprof'>这里</a>打开调试页面")
		})
		index.Any("/pprof", pprof.New())
		index.Any("/pprof/{action:path}", pprof.New())
	}
	return module.NewModule("/debug", handler)
}
