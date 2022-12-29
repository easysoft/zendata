package server

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/server/core/module"
	"github.com/easysoft/zendata/internal/server/index"
	"github.com/kataras/iris/v12"
)

type IndexModule struct {
	CommModule  *index.CommModule  `inject:""`
	DefModule   *index.DefModule   `inject:""`
	DataModule  *index.DataModule  `inject:""`
	AdminModule *index.AdminModule `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

func (m *IndexModule) Party() module.WebModule {
	handler := func(v1 iris.Party) {}

	modules := []module.WebModule{
		m.DefModule.Party(),
		m.CommModule.Party(),
		m.DataModule.Party(),
		m.AdminModule.Party(),
	}
	return module.NewModule(consts.ApiPath, handler, modules...)
}
