package controller

import (
	serverService "github.com/easysoft/zendata/internal/server/service"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type MockCtrl struct {
	MockService *serverService.MockService `inject:""`
	BaseCtrl
}

func (c *MockCtrl) Mock(ctx iris.Context) {
	paths := ctx.Params().Get("paths")
	code := ctx.URLParam("code")
	if code == "" {
		code = "200"
	}

	if vari.GlobalVars.MockData == nil {
		c.MockService.Init()
	}

	resp, _ := c.MockService.GetResp(paths, ctx.Method(), code)

	ctx.JSON(resp, context.JSON{Indent: "    "})
}
