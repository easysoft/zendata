package index

import (
	"github.com/easysoft/zendata/internal/server/controller"
	"github.com/easysoft/zendata/internal/server/core/module"
	"github.com/kataras/iris/v12"
)

type AdminModule struct {
	AdminCtrl *controller.AdminCtrl `inject:""`
}

func NewAdminModule() *AdminModule {
	return &AdminModule{}
}

// Party 执行
func (m *AdminModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Post("/", m.AdminCtrl.Handle).Name = "所有管理请求"
	}

	return module.NewModule("/admin", handler)
}
