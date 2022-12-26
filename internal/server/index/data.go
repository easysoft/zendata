package index

import (
	"github.com/easysoft/zendata/internal/server/controller"
	"github.com/easysoft/zendata/internal/server/core/module"
	"github.com/kataras/iris/v12"
)

type DataModule struct {
	DataCtrl *controller.DataCtrl `inject:""`
}

func NewDataModule() *DataModule {
	return &DataModule{}
}

// Party 执行
func (m *DataModule) Party() module.WebModule {
	handler := func(index iris.Party) {
		index.Get("/generate", m.DataCtrl.GenerateByFile).Name = "通过执行文件路径生成数据"
		index.Post("/generate", m.DataCtrl.GenerateByContent).Name = "通过推送文件内容生成数据"
		index.Post("/decode", m.DataCtrl.Decode).Name = "反向解析数据"
	}

	return module.NewModule("/data", handler)
}
