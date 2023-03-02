package commandConfig

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB() (gormDb *gorm.DB, err error) {
	gormDb, err = gorm.Open(sqlite.Open(consts.SqliteFile), &gorm.Config{})

	if vari.Verbose {
		gormDb = gormDb.Debug()
	}

	err = gormDb.AutoMigrate(
		model.Models...,
	)

	return
}
