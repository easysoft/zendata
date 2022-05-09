package main

import (
	"flag"
	"fmt"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/fatih/color"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"strings"
	"time"
)

const (
	dbUser   = "root"
	dbUrl    = "127.0.0.1:3306"
	dbPasswd = "P2ssw0rd"
	dbName   = "zendata"

	truncateTable = `truncate table %s;`

	createTableTempl = `CREATE TABLE IF NOT EXISTS %s (
		id bigint auto_increment,
		content varchar(1000) not null unique,
		tag varchar(50),
		primary key(id)
	) engine=innodb default charset=utf8 auto_increment=1;`

	//deleteAllSql = "DELETE FROM %s WHERE 1=1;"

	insertTemplate = "INSERT INTO %s (content) VALUES %s;"
)

var (
	db *gorm.DB

	tableName string
	filePath  string
	colNum    int
)

func main() {
	flag.StringVar(&tableName, "t", "", "")
	flag.StringVar(&filePath, "f", "", "")
	flag.IntVar(&colNum, "c", 0, "")

	flag.Parse()

	tableName = "biz_" + strings.TrimLeft(tableName, "biz_")
	db := GetDB()

	truncateTableSql := fmt.Sprintf(truncateTable, tableName)
	err := db.Exec(truncateTableSql).Error
	if err != nil {
		fmt.Printf("create table %s failed, err %s", tableName, err.Error())
		return
	}

	createTableSql := fmt.Sprintf(createTableTempl, tableName)
	err = db.Exec(createTableSql).Error
	if err != nil {
		fmt.Printf("create table %s failed, err %s", tableName, err.Error())
		return
	}

	//deleteDataSql := fmt.Sprintf(deleteAllSql, tableName)
	//err = db.Exec(deleteDataSql).Error
	//if err != nil {
	//	fmt.Printf("insert data failed, err %s", err.Error())
	//	return
	//}

	content := fileUtils.ReadFile(filePath)
	insertSqls := make([]string, 0)

	for _, line := range strings.Split(content, "\n") {
		arr := strings.Split(strings.TrimSpace(line), " ")

		if colNum >= len(arr) {
			continue
		}

		content := arr[colNum]
		insert := fmt.Sprintf("('%s')", content)
		insertSqls = append(insertSqls, insert)
	}

	sql := fmt.Sprintf(insertTemplate, tableName, strings.Join(insertSqls, ","))
	err = db.Exec(sql).Error
	if err != nil {
		fmt.Printf("insert data failed, err %s", err.Error())
		return
	}
}

func GetDB() *gorm.DB {
	if db != nil {
		return db
	}

	db = GormMySQL()
	_ = db.Use(
		dbresolver.Register(
			dbresolver.Config{}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	)

	db.Session(&gorm.Session{FullSaveAssociations: true, AllowGlobalUpdate: false})

	err := db.AutoMigrate(
		&DataCategory1{},
		&DataCategory2{},
	)
	if err != nil {
		fmt.Printf(color.RedString("migrate models failed, error: %s.", err.Error()))
	}

	return db
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
	if db, err := gorm.Open(mysql.New(mysqlConfig)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(0)
		sqlDB.SetMaxOpenConns(0)
		return db
	}
}

type DataCategory1 struct {
	BaseModel

	Name string `json:"name"`
	Desc string `json:"desc" gorm:"column:descr"`
}

func (DataCategory1) TableName() string {
	return "biz_data_category1"
}

type DataCategory2 struct {
	BaseModel

	Name      string `json:"name"`
	Desc      string `json:"desc" gorm:"column:descr"`
	DataTable string `json:"desc"`

	ParentId uint `json:"parentId"`
}

func (DataCategory2) TableName() string {
	return "biz_data_category2"
}

type BaseModel struct {
	ID        uint       `gorm:"primary_key" sql:"type:INT(10) UNSIGNED NOT NULL" json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Deleted   bool       `json:"-" gorm:"default:false"`
	Disabled  bool       `json:"disabled,omitempty" gorm:"default:false"`
}
