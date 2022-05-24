package serverConfig

import (
	"fmt"

	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB() (gormDb *gorm.DB, err error) {
	gormDb, err = gorm.Open(sqlite.Open(constant.SqliteFile), &gorm.Config{})

	if vari.Verbose {
		gormDb = gormDb.Debug()
	}

	if vari.RunMode == constant.RunModeServer {
		for _, model := range model.Models {
			if err := gormDb.Set("gorm:table_options", "").AutoMigrate(model); err != nil {
				return nil, fmt.Errorf("auto migrate table %+v failure %s", model, err.Error())
			}
		}
	}

	return
}
