package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/service"
)

type ResService struct {
}

func (s *ResService) LoadRes(resType string) (res map[string][]model.ResFile) {
	res, _, _ = service.LoadRes(resType)
	return
}
