package controller

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverService "github.com/easysoft/zendata/internal/server/service"
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

func (c *MockCtrl) Upload(ctx iris.Context) {
	f, fh, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(c.ErrResp(consts.CommErr, err.Error()))
		return
	}
	defer f.Close()

	name, spec, mockConf, dataConf, pth, err := c.MockService.Upload(ctx, fh)
	if err != nil {
		ctx.JSON(c.ErrResp(consts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(iris.Map{"name": name, "spec": spec, "mock": mockConf, "data": dataConf, "path": pth}))
}

func (c *MockCtrl) GetPreviewData(ctx iris.Context) {
	id, _ := ctx.URLParamInt("id")

	data, _ := c.MockService.GetPreviewData(id)

	ctx.JSON(c.SuccessResp(data))
}

func (c *MockCtrl) GetPreviewResp(ctx iris.Context) {
	req := domain.MockPreviewReq{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(consts.ParamErr, err.Error()))
	}

	data, _ := c.MockService.GetPreviewResp(req)

	ctx.JSON(c.SuccessResp(data))
}

func (c *MockCtrl) Save(ctx iris.Context) {
	req := model.ZdMock{}
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.JSON(c.ErrResp(consts.ParamErr, err.Error()))
	}

	err := c.MockService.Save(&req)
	if err != nil {
		ctx.JSON(c.ErrResp(consts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(iris.Map{"id": req.ID}))
}

func (c *MockCtrl) Remove(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(consts.CommErr, err.Error()))
		return
	}

	err = c.MockService.Remove(id)
	if err != nil {
		ctx.JSON(c.ErrResp(consts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *MockCtrl) StartMockService(ctx iris.Context) {
	id, err := ctx.URLParamInt("id")
	if err != nil {
		ctx.JSON(c.ErrResp(consts.CommErr, err.Error()))
		return
	}

	act := ctx.URLParam("act")

	if act != "stop" {
		err = c.MockService.StartMockService(id)
	} else {
		err = c.MockService.StopMockService(id)
	}

	if err != nil {
		ctx.JSON(c.ErrResp(consts.CommErr, err.Error()))
		return
	}

	ctx.JSON(c.SuccessResp(nil))
}

func (c *MockCtrl) Mock(ctx iris.Context) {
	paths := ctx.Params().Get("paths")
	mediaType := ctx.URLParam("contentType")
	code := ctx.URLParam("code")
	if code == "" {
		code = "200"
	}
	if mediaType == "" {
		mediaType = ctx.GetHeader("Content-Type")
		if mediaType == "" {
			mediaType = "application/json"
		}
	}

	resp, err := c.MockService.GetResp(paths, ctx.Method(), code, mediaType)
	if err != nil {
		resp = iris.Map{"err": err.Error()}
	}

	ctx.JSON(resp, context.JSON{Indent: "    "})
}
