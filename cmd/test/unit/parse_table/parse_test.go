package parse_table

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/easysoft/zendata/internal/command"
	commandConfig "github.com/easysoft/zendata/internal/command/config"
	configUtils "github.com/easysoft/zendata/internal/pkg/config"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Dsn = "" // "root:P2ssw0rd@(127.0.0.1:3306)/zd_test?charset=utf8&parseTime=True&loc=Local"
	Out = "" // "demo/out/sql"

	TableByFieldName = "by_field_name"
	TableByFieldType = "by_field_type"
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

func TestGenDefFromFieldName(t *testing.T) {
	vari.GlobalVars.Table = TableByFieldName
	parse(t)
	checkByFieldName(t)
}

func TestGenDefFromColumnType(t *testing.T) {
	vari.GlobalVars.Table = TableByFieldType
	parse(t)
	checkByFieldType(t)
}

func TestGenDefFromRecords(t *testing.T) {
	vari.GlobalVars.Table = TableByRecords
	parse(t)
	checkByFieldRecords(t)
}

func parse(t *testing.T) {
	mainCtrl, _ := command.InitCtrl()
	mainCtrl.GenYaml("")
}

func checkByFieldName(t *testing.T) {
	def := readYamlObj(TableByFieldName)

	for _, field := range def.Fields {
		pass := true
		if field.Field == "username" {
			pass = field.From == "name.enaccount.v1.yaml"
		} else if field.Field == "telphone" {
			pass = field.From == "phone.v1.yaml"
		} else if field.Field == "mobilephone" {
			pass = field.From == "phone.v1.yaml"
		} else if field.Field == "email" {
			pass = field.From == "email.v1.yaml"
		} else if field.Field == "url" {
			pass = field.From == "domain.domain.v1.yaml"
		} else if field.Field == "ip" {
			pass = field.From == "ip.v1.yaml"
		} else if field.Field == "macaddress" {
			pass = field.Format == "mac()"
		} else if field.Field == "creditcard" {
			pass = field.Format == "credit_card(amex)"
		} else if field.Field == "idcard" {
			pass = field.Format == "id_card()"
		} else if field.Field == "token" {
			pass = field.Format == "token()"
		}

		if !pass {
			t.Errorf(`Wrong Config: field "%s"; from: "%v"; format: "%v""`, field.Field, field.From, field.Format)
		}
	}
}

func checkByFieldType(t *testing.T) {
	def := readYamlObj(TableByFieldType)

	for _, field := range def.Fields {
		pass := true
		if field.Field == "f_bit" {
			pass = field.Range == "0,1:R"
		} else if field.Field == "f_tinyint" {
			pass = field.Range == "0-255:R"
		} else if field.Field == "f_smallint" {
			pass = field.Range == "0-65535:R"
		} else if field.Field == "f_mediumint" {
			pass = field.Range == "0-16777215:R"
		} else if field.Field == "f_int" {
			pass = field.Range == "0-4294967295:R"
		} else if field.Field == "f_bigint" {
			pass = field.Range == "0-9223372036854775807:R"
		} else if field.Field == "f_float" {
			pass = field.Range == "0-99.999:R"
		} else if field.Field == "f_double" {
			pass = field.Range == "0-99.999999:R"
		} else if field.Field == "f_decimal" {
			pass = field.Range == "0-99.99:R"
		} else if field.Field == "f_char" {
			pass = field.Range == "a-z"
		} else if field.Field == "f_tinytext" {
			pass = field.From == "idiom.v1.idiom"
		} else if field.Field == "f_text" {
			pass = field.From == "xiehouyu.v1.xiehouyu"
		} else if field.Field == "f_mediumtext" {
			pass = field.From == "joke.v1.joke"
		} else if field.Field == "f_longtext" {
			pass = field.From == "song.v1.song"

		} else if field.Field == "f_tinyblob" {
			pass = field.Format == "binary()"
		} else if field.Field == "f_blob" {
			pass = field.Format == "binary()"
		} else if field.Field == "f_mediumblob" {
			pass = field.Format == "binary()"
		} else if field.Field == "f_longblob" {
			pass = field.Format == "binary()"
		} else if field.Field == "f_binary" {
			pass = field.Format == "binary()"
		} else if field.Field == "f_varbinary" {
			pass = field.Format == "binary()"

		} else if field.Field == "f_date" {
			pass = field.Range == "(-1M)-(+1w):86400" && field.Format == "YY/MM/DD" && field.Type == "timestamp"
		} else if field.Field == "f_time" {
			pass = field.Range == "(-1M)-(+1w):60" && field.Format == "hh:mm:ss" && field.Type == "timestamp"
		} else if field.Field == "f_year" {
			pass = field.Range == "(-6Y)-(+6Y):31536000" && field.Format == "YYYY" && field.Type == "timestamp"
		} else if field.Field == "f_datetime" {
			pass = field.Range == "(-1M)-(+1w):60" && field.Format == "YY/MM/DD hh:mm:ss" && field.Type == "timestamp"
		} else if field.Field == "f_timestamp" {
			pass = field.Range == "(-1M)-(+1w):60" && field.Type == "timestamp"

		} else if field.Field == "f_enum" {
			pass = field.Range == "a,b,c"
		} else if field.Field == "f_set" {
			pass = field.Range == "a,b,c" && field.Loop == "2-3"
		}
		if !pass {
			t.Errorf(`Wrong Config: field "%s"; from: "%v"; format: "%v""`, field.Field, field.From, field.Format)
		}
	}
}

func checkByFieldRecords(t *testing.T) {
	def := readYamlObj(TableByRecords)

	for _, field := range def.Fields {
		pass := true

		if field.Field == "f1" {
			pass = field.From == "email.v1.yaml"
		} else if field.Field == "f2" {
			pass = field.Format == "credit_card(amex)"
		} else if field.Field == "f3" {
			pass = field.Format == "mac()"
		} else if field.Field == "f4" {
			pass = field.From == "uuid.v1.yaml"
		} else if field.Field == "f5" {
			pass = field.Format == "md5"
		} else if field.Field == "f6" {
			pass = field.From == "phone.v1.yaml" && field.Use == "cellphone"
		} else if field.Field == "f7" {
			pass = field.From == "phone.v1.yaml" && field.Use == "telephone_china"
		} else if field.Field == "f8" {
			pass = field.Format == "id_card()"
		} else if field.Field == "f9" {
			pass = field.From == "domain.domain.v1.yaml"
		} else if field.Field == "f10" {
			pass = field.Format == "json()"
		}

		if !pass {
			t.Errorf(`Wrong Config: field "%s"; from: "%v"; format: "%v""`, field.Field, field.From, field.Format)
		}
	}
}

func initTables() {
	initTable(TableByFieldName)
	initTable(TableByFieldType)
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

func readYamlObj(name string) (def model.ZdDef) {
	pth := filepath.Join(vari.GlobalVars.Output, name+".yaml")
	content := fileUtils.ReadFileBuf(pth)

	def = model.ZdDef{}
	yaml.Unmarshal(content, &def)

	return
}
