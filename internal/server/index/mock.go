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
		index.Post("/", m.MockCtrl.Save).Name = "保存Mock"
		index.Delete("/{id:uint}", m.MockCtrl.Remove).Name = "删除Mock"

		index.Post("/upload", m.MockCtrl.Upload).Name = "上传Spec"
		index.Get("/getPreviewData", m.MockCtrl.GetPreviewData).Name = "上传Spec"
		index.Post("/getPreviewResp", m.MockCtrl.GetPreviewResp).Name = "上传Spec"

		index.Post("/startMockService", m.MockCtrl.StartMockService).Name = "启动Mock服务"
		index.Post("/stopMockService", m.MockCtrl.StopMockService).Name = "关闭Mock服务"
	}

	return module.NewModule("/mocks", handler)
}

func (m *MockModule) PartyData() module.WebModule {
	handler := func(index iris.Party) {
		index.Any("/{paths:path}", m.MockCtrl.Mock) // mock data url
	}

	return module.NewModule("/m", handler)
}
