package server

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/server/core/module"
	"github.com/easysoft/zendata/internal/server/index"
	"github.com/kataras/iris/v12"
)

type IndexModule struct {
	DefModule *index.DefModule `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

func (m *IndexModule) Party() module.WebModule {
	handler := func(v1 iris.Party) {}

	modules := []module.WebModule{
		m.DefModule.Party(),
	}
	return module.NewModule(constant.ApiPath, handler, modules...)
}
