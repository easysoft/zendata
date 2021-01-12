package serverConfig

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/jinzhu/gorm"
)

func NewGormDB(config *Config) (gormDb *gorm.DB, err error) {
	gormDb, err = gorm.Open(config.DBDriver, config.DBPath)
	gormDb = gormDb.Debug()
	for _, model := range model.Models {
		if err := gormDb.Set("gorm:table_options", "").AutoMigrate(model).Error; err != nil {
			return nil, fmt.Errorf("auto migrate table %+v failure %s", model, err.Error())
		}
	}

	return
}
