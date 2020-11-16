package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/service"
	constant "github.com/easysoft/zendata/src/utils/const"
)

type ResService struct {
}

func (s *ResService) LoadRes(resType string) (ret []model.ResFile) {
	res, _, _ := service.LoadRes(resType)

	for _, key := range constant.ResKeys {
		ret = append(ret,  res[key]...)
	}

	return
}
