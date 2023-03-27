package helper

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"strings"
)

func GenFieldDefByMetadata(metadata string, name string, records []string) (ret FieldTypeInfo) {
	columnType, varcharType := GetColumnType(metadata, name, records)

	if columnType == consts.Varchar && varcharType != "" {
		ret = GenDefByVarcharType(varcharType)
	}

	ret = GenDefByColumnType(columnType)

	return
}

func GenDefByVarcharType(varcharType consts.VarcharType) (ret FieldTypeInfo) {
	return
}

func GetColumnType(metaType string, name string, records []string) (
	columnType consts.ColumnType, varcharType consts.VarcharType) {
	metaType = strings.ToLower(metaType)

	if metaType == "integer" {
		metaType = "int"
	}

	columnType = consts.ColumnType(metaType)

	if columnType == consts.Varchar {
		varcharType = GetVarcharTypeByName(name)
	}

	if columnType == consts.Varchar {
		varcharType = GetVarcharTypeByRecords(records)
	}

	return
}

func GetColumnTypeByMetadata(metadata string) (typ consts.ColumnType) {
	return
}

func GetVarcharTypeByName(name string) (ret consts.VarcharType) {
	ret = consts.Empty
	return
}

func GetVarcharTypeByRecords(records []string) (ret consts.VarcharType) {
	ret = consts.Empty
	return
}

type FieldTypeInfo struct {
	FieldType consts.ColumnType
	Note      string
	Rang      string
}

func GenDefByColumnType(fieldType consts.ColumnType) (ret FieldTypeInfo) {
	ran := ""
	note := ""

	switch fieldType {

	// integer
	case "BIT":
		ran = "0,1"
	case "TINYINT":
		ran = "0-255"
	case "SMALLINT":
		ran = "0-65535"
	case "MEDIUMINT":
		ran = "0-65535"
		note = `"MEDIUMINT [0,2^24-1]"`
	case "INT":
		ran = "0-100000"
		note = `"INI [0,2^32-1]"`
	case "BIGINT":
		ran = "0-100000"
		note = `"BIGINT [0,2^64-1]"`
	// floating-point
	case "FLOAT":
		ran = "1.01-99.99:0.01"
		note = `"FLOAT"`
	case "DOUBLE":
		ran = "1.01-99.99:0.01"
		note = `"DOUBLE"`
	// fixed-point
	case "DECIMAL":
		ran = "123.45"
		note = `"DECIMAL"`
	// character string
	case "CHAR":
		ran = `"a-z"`
	case "TINYTEXT":
		ran = `"TINYTEXT"`
	case "TEXT":
		ran = `"TEXT"`
	case "MEDIUMTEXT":
		ran = `"MEDIUMTEXT"`
	case "LONGTEXT":
		ran = `"LONGTEXT"`
	// binary data
	case "TINYBLOB":
		ran = `"TINYBLOB"`
	case "BLOB":
		ran = `"BLOB"`
	case "MEDIUMBLOB":
		ran = `"MEDIUMBLOB"`
	case "LONGBLOB":
		ran = `"LONGBLOB"`
	case "BINARY":
		ran = "BINARY"
	case "VARBINARY":
		ran = "VARBINARY"
	// Date and time type
	case "DATE":
		ran = `"DATE"`
	case "TIME":
		ran = `"TIME"`
	case "YEAR":
		ran = `"YEAR"`
	case "DATETIME":
		ran = `"DATETIME"`
	case "TIMESTAMP":
		ran = `"TIMESTAMP"`
	// other type
	case "ENUM":
		ran = `"ENUM"`
	case "SET":
		ran = `"SET"`
	case "GEOMETRY":
		ran = `"GEOMETRY"`
	case "POINT":
		ran = `"POINT"`
	case "LINESTRING":
		ran = `"LINESTRING"`
	case "POLYGON":
		ran = `"POLYGON"`
	case "MULTIPOINT":
		ran = `"MULTIPOINT"`
	case "MULTILINESTRING":
		ran = `"MULTILINESTRING"`
	case "MULTIPOLYGON":
		ran = `"MULTIPOLYGON"`
	case "GEOMETRYCOLLECTION":
		ran = `"GEOMETRYCOLLECTION"`
	default:
	}

	ret.Rang = ran
	ret.Note = note
	ret.FieldType = fieldType

	return ret
}
