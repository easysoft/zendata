package serverConfig

import (
	"gorm.io/gorm/logger"

	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB() (gormDb *gorm.DB, err error) {
	gormDb, err = gorm.Open(sqlite.Open(constant.SqliteFile), &gorm.Config{
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
