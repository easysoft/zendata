package controller

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverService "github.com/easysoft/zendata/internal/server/service"
	"github.com/kataras/iris/v12"
)

type DefCtrl struct {
	DefService *serverService.DefService `inject:""`
	BaseCtrl
}

func NewDefCtrl() *DefCtrl {
	return &DefCtrl{}
}

func (c *DefCtrl) List(ctx iris.Context) {
	req := model.ReqData{}

	if err := ctx.ReadQuery(&req); err != nil {
		ctx.JSON(c.ErrResp(constant.CommErr, err.Error()))
		return
	}

	data, total := c.DefService.List(req.Keywords, req.Page)

	ctx.JSON(c.SuccessResp(iris.Map{"data": data, "total": total}))
}

func (c *DefCtrl) Get(ctx iris.Context) {
	//id, err := ctx.Params().GetInt("id")
	//if err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	//	return
	//}
	//
	//po, err := c.DefService.Get(uint(id))
	//if err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	//	return
	//}
	//ctx.JSON(c.SuccessResp(po))
}

func (c *DefCtrl) Create(ctx iris.Context) {
	//req := model.Def{}
	//if err := ctx.ReadJSON(&req); err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
	//}
	//
	//id, err := c.DefService.Create(req)
	//if err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.ErrZentaoConfig, err.Error()))
	//	return
	//}
	//
	//ctx.JSON(c.SuccessResp(iris.Map{"id": id}))
}

func (c *DefCtrl) Update(ctx iris.Context) {
	//req := model.Def{}
	//if err := ctx.ReadJSON(&req); err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.ParamErr, err.Error()))
	//}
	//
	//err := c.DefService.Update(req)
	//if err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.ErrZentaoConfig, err.Error()))
	//	return
	//}
	//
	//ctx.JSON(c.SuccessResp(iris.Map{"id": req.ID}))
}

func (c *DefCtrl) Delete(ctx iris.Context) {
	//id, err := ctx.Params().GetInt("id")
	//if err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	//	return
	//}
	//
	//err = c.DefService.Delete(uint(id))
	//if err != nil {
	//	ctx.JSON(c.ErrResp(commConsts.CommErr, err.Error()))
	//	return
	//}
	//
	//ctx.JSON(c.SuccessResp(nil))
}
