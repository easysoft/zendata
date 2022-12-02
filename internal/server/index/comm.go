package index

import (
	"github.com/easysoft/zendata/internal/server/controller"
	"github.com/easysoft/zendata/internal/server/core/module"
	"github.com/kataras/iris/v12"
)

type CommModule struct {
	CommCtrl *controller.CommCtrl `inject:""`
}

// Party 执行
func (m *CommModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Get("/getWorkDir", m.CommCtrl.GetWorkDir).Name = "获取当前工作目录"
		index.Post("/syncData", m.CommCtrl.SyncData).Name = "同步当前工作数据"
	}
	return module.NewModule("/comm", handler)
}
