package index

import (
	"github.com/easysoft/zendata/internal/server/controller"
	"github.com/easysoft/zendata/internal/server/core/module"
	"github.com/kataras/iris/v12"
)

type MockModule struct {
	MockCtrl *controller.MockCtrl `inject:""`
}

func NewMockModule() *DataModule {
	return &DataModule{}
}

// Party 执行
func (m *MockModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Get("/", m.MockCtrl.List).Name = "Mock列表"
		index.Get("/{id:uint}", m.MockCtrl.Get).Name = "Mock详情"
		index.Post("/upload", m.MockCtrl.Upload).Name = "上传Spec"

		index.Any("/{paths:path}", m.MockCtrl.Mock) // mock data url
	}

	return module.NewModule("/mocks", handler)
}
