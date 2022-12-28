package controller

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
)

type BaseCtrl struct {
}

func NewBaseCtrl() *BaseCtrl {
	return &BaseCtrl{}
}

func (c *BaseCtrl) SuccessResp(data interface{}) (ret domain.Response) {
	ret = domain.Response{Code: consts.Success.Code, Data: data}

	return
}

func (c *BaseCtrl) ErrResp(respCode consts.ResponseCode, msg string) (ret domain.Response) {
	ret = domain.Response{Code: respCode.Code, Msg: c.ErrMsg(respCode, msg)}

	return
}

func (c *BaseCtrl) BizErrResp(err *domain.BizError, msg string) (ret domain.Response) {
	ret = domain.Response{Code: err.Code, Msg: msg}

	return
}

func (c *BaseCtrl) ErrMsg(err consts.ResponseCode, msg string) (ret string) {
	ret += msg

	return
}
