package server

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/vari"
)

func ListData() (defs []model.Data) {
	vari.GormDB.Find(&defs)

	return
}
