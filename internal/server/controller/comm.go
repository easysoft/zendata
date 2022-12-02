package controller

import (
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/kataras/iris/v12"
)

type CommCtrl struct {
	BaseCtrl
}

func (c *CommCtrl) GetWorkDir(ctx iris.Context) {
	ctx.JSON(c.SuccessResp(vari.ZdPath))
}

func (c *CommCtrl) SyncData(ctx iris.Context) {
	ctx.JSON(c.SuccessResp(vari.ZdPath))
}
