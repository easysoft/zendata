package serverConfig

import (
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> 3.0
	"gorm.io/gorm/logger"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB() (gormDb *gorm.DB, err error) {
<<<<<<< HEAD
	gormDb, err = gorm.Open(sqlite.Open(constant.SqliteFile), &gorm.Config{
=======
	gormDb, err = gorm.Open(sqlite.Open(consts.SqliteFile), &gorm.Config{
>>>>>>> 3.0
		Logger: logger.Default.LogMode(logger.Info),
	})

	if vari.Verbose {
		gormDb = gormDb.Debug()
	}

	err = gormDb.AutoMigrate(
		model.Models...,
	)

	return
}
