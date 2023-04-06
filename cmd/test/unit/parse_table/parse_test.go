package parse_table

import (
	"flag"
	"github.com/easysoft/zendata/internal/command"
	commandConfig "github.com/easysoft/zendata/internal/command/config"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var (
	Dsn = "" // "root:P2ssw0rd@(127.0.0.1:3306)/zd_test?charset=utf8&parseTime=True&loc=Local"
	Out = "" // "demo/out/sql"

	TableByFieldType = "by_field_type"
	TableByFieldName = "by_field_name"
	TableByRecords   = "by_records"
)

func setup() {
	initArgs()

	os.Chdir("../../../../")
	configUtils.InitConfig("")
	vari.DB, _ = commandConfig.NewGormDB()

	initTables()
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestGenDefFromColumnDesc(t *testing.T) {
	vari.GlobalVars.Table = "by_field_type"
	gen(t)
}

func TestGenDefFromFieldName(t *testing.T) {
	vari.GlobalVars.Table = "by_field_name"
	gen(t)
}

func TestGenDefFromRecords(t *testing.T) {
	vari.GlobalVars.Table = "by_records"
	gen(t)
}

func gen(t *testing.T) {
	mainCtrl, _ := command.InitCtrl()
	mainCtrl.GenYaml("")
}

func initTables() {
	initTable(TableByFieldType)
	initTable(TableByFieldName)
	initTable(TableByRecords)
}

func initTable(name string) {
	dir := filepath.Join("cmd", "test", "unit", "parse_table")
	sqls := fileUtils.ReadFile(filepath.Join(dir, name+".sql"))

	arr := strings.Split(sqls, ";")

	db, err := gorm.Open(mysql.Open(vari.GlobalVars.DBDsn))
	if err != nil {
		panic(err)
	}

	for _, sql := range arr {
		if sql == "" {
			continue
		}
		err = db.Exec(sql).Error
		if err != nil {
			panic(err)
		}
	}
}

func initArgs() {
	flag.Parse()
	argList := flag.Args()

	vari.GlobalVars.DBDsn = argList[0]
	vari.GlobalVars.Output = argList[1]
	log.Printf("DBDsn = %s", vari.GlobalVars.DBDsn)
	log.Printf("Output = %s", vari.GlobalVars.Output)
}
