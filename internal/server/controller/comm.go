package controller

import (
	serverService "github.com/easysoft/zendata/internal/server/service"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/kataras/iris/v12"
)

type CommCtrl struct {
	SyncService *serverService.SyncService `inject:""`

	BaseCtrl
}

func (c *CommCtrl) GetWorkDir(ctx iris.Context) {
	ctx.JSON(c.SuccessResp(vari.WorkDir))
}

func (c *CommCtrl) SyncData(ctx iris.Context) {
	c.SyncService.SyncData()

	ctx.JSON(c.SuccessResp(""))
}
