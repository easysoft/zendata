package controller

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverService "github.com/easysoft/zendata/internal/server/service"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type MockCtrl struct {
	MockService *serverService.MockService `inject:""`
	BaseCtrl
}

func (c *MockCtrl) List(ctx iris.Context) {
	req := model.ReqData{}

	err := ctx.ReadQuery(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(consts.CommErr, err.Error()))
		return
	}

	list, total, _ := c.MockService.List(req.Keywords, req.Page)

	ctx.JSON(c.SuccessResp(iris.Map{"list": list, "total": total}))
}

func (c *MockCtrl) Get(ctx iris.Context) {
	id, _ := ctx.Params().GetInt("id")

	data, _ := c.MockService.Get(id)

	ctx.JSON(c.SuccessResp(data))
}

func (c *MockCtrl) Mock(ctx iris.Context) {
	paths := ctx.Params().Get("paths")
	mediaType := ctx.Params().Get("paths")
	code := ctx.URLParam("code")
	if code == "" {
		code = "200"
	}

	if vari.GlobalVars.MockData == nil {
		c.MockService.Init()
	}

	resp, _ := c.MockService.GetResp(paths, ctx.Method(), code, mediaType)

	ctx.JSON(resp, context.JSON{Indent: "    "})
}
