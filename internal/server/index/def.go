package index

import (
	"github.com/easysoft/zendata/internal/server/controller"
	"github.com/easysoft/zendata/internal/server/core/module"
	"github.com/kataras/iris/v12"
)

type DefModule struct {
	DefCtrl *controller.DefCtrl `inject:""`
}

func NewDefModule() *DefModule {
	return &DefModule{}
}

// Party 执行
func (m *DefModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Get("/", m.DefCtrl.List).Name = "列表"
		index.Get("/{id:int}", m.DefCtrl.Get).Name = "详情"
		index.Post("/", m.DefCtrl.Create).Name = "新建"
		index.Put("/{id:int}", m.DefCtrl.Update).Name = "更新"
		index.Delete("/{id:int}", m.DefCtrl.Delete).Name = "删除"

		index.Post("/sync", m.DefCtrl.Create).Name = "同步"
	}
	return module.NewModule("/sites", handler)
}
