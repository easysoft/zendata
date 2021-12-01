package serverConfig

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/gorm"
)

func NewGormDB() (gormDb *gorm.DB, err error) {
	gormDb, err = gorm.Open(constant.SqliteDriver, constant.SqliteFile)

	if vari.Verbose {
		gormDb = gormDb.Debug()
	}

	if vari.RunMode == constant.RunModeServer {
		for _, model := range model.Models {
			if err := gormDb.Set("gorm:table_options", "").AutoMigrate(model).Error; err != nil {
				return nil, fmt.Errorf("auto migrate table %+v failure %s", model, err.Error())
			}
		}
	}

	return
}
