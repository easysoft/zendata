package comm

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/model"
	"github.com/fatih/color"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"time"
)

const (
	dbUser   = "root"
	dbUrl    = "127.0.0.1:3306"
	dbPasswd = "P2ssw0rd"
	dbName   = "zendata"
)

var (
	DB *gorm.DB
)

func GetDB() *gorm.DB {
	if DB != nil {
		return DB
	}

	DB = GormMySQL()
	_ = DB.Use(
		dbresolver.Register(
			dbresolver.Config{}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)

	DB.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: false})

	err := DB.AutoMigrate(
		&model.DataCategory1{},
		&model.DataCategory2{},
	)
	if err != nil {
		fmt.Printf(color.RedString("migrate models failed, error: %s.", err.Error()))
	}

	return DB
}

func GormMySQL() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", dbUser, dbPasswd, dbUrl, dbName, "charset=utf8mb4&parseTime=True&loc=Local")

	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(0)
		sqlDB.SetMaxOpenConns(0)
		return db
	}
}
