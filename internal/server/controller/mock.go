package controller

import (
	serverService "github.com/easysoft/zendata/internal/server/service"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/kataras/iris/v12"
	"log"
)

type MockCtrl struct {
	MockService *serverService.MockService `inject:""`
	BaseCtrl
}

func (c *MockCtrl) Mock(ctx iris.Context) {
	paths := ctx.Params().Get("paths")
	params := ctx.URLParams()

	if vari.GlobalVars.MockData == nil {
		c.MockService.Init()
	}

	log.Print(paths, ", ", params)
}
