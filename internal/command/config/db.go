package commandConfig

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB() (gormDb *gorm.DB, err error) {
	gormDb, err = gorm.Open(sqlite.Open(constant.SqliteFile), &gorm.Config{})

	if vari.Verbose {
		gormDb = gormDb.Debug()
	}

	return
}
