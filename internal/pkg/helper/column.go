package helper

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GenerateFieldDefByMetadata(metadata string, param string, name string, records []interface{}) (info FieldTypeInfo) {
	GetColumnType(metadata, name, records, &info)

	if info.ColumnType == consts.Varchar && info.VarcharType != "" {
		generateDefByVarcharType(param, &info)
		return
	}

	GenDefByColumnType(param, &info)

	return
}

func generateDefByVarcharType(param string, ret *FieldTypeInfo) {

	return
}

func GetColumnType(metaType string, name string, records []interface{}, info *FieldTypeInfo) {
	metaType = strings.ToLower(metaType)

	if metaType == "integer" {
		metaType = "int"
	}

	info.ColumnType = consts.ColumnType(metaType)

	if info.ColumnType == consts.Varchar {
		info.VarcharType = GetVarcharTypeByName(name)
	}

	if info.ColumnType == consts.Varchar && info.VarcharType == "" && len(records) > 0 {
		info.VarcharType = GetVarcharTypeByRecords(records)
	}

	if info.ColumnType == consts.Float || info.ColumnType == consts.Double {
		GetPrecisionByRecords(records, info)
	}

	SetIsNum(info)
	if info.IsNum {
		GetSignByRecords(records, info)
	}

	return
}

func SetIsNum(info *FieldTypeInfo) {
	info.IsNum = info.ColumnType == consts.Tinyint || info.ColumnType == consts.Smallint ||
		info.ColumnType == consts.Mediumint || info.ColumnType == consts.Int || info.ColumnType == consts.Bigint ||
		info.ColumnType == consts.Float || info.ColumnType == consts.Double
}

func GetSignByRecords(records []interface{}, info *FieldTypeInfo) {
	if len(records) == 0 {
		return
	}
	num := getRandNum(len(records))
	computerPrecision(records[num], info)

	num = getRandNum(len(records))
	computerPrecision(records[num], info)

	num = getRandNum(len(records))
	computerPrecision(records[num], info)
}
func GetPrecisionByRecords(records []interface{}, info *FieldTypeInfo) {
	if len(records) == 0 {
		return
	}

	num := getRandNum(len(records))
	computerSign(records[num], info)

	num = getRandNum(len(records))
	computerSign(records[num], info)

	num = getRandNum(len(records))
	computerSign(records[num], info)
}
func computerPrecision(val interface{}, info *FieldTypeInfo) {
	str := fmt.Sprintf("%v", val)

	index := strings.LastIndex(str, ".")
	subStr := str[index+1:]

	info.Precision = len(subStr)
}
func computerSign(val interface{}, info *FieldTypeInfo) {
	str := fmt.Sprintf("%v", val)
	float, _ := strconv.ParseFloat(str, 64)

	info.HasSign = info.HasSign || float < 0
}

func GetVarcharTypeByName(name string) (ret consts.VarcharType) {
	ret = consts.Empty

	if regexp.MustCompile(`(user).*(name)`).MatchString(name) {
		ret = consts.Username
	} else if regexp.MustCompile(`(phone|number)`).MatchString(name) {
		if regexp.MustCompile(`(tel)`).MatchString(name) {
			ret = consts.TelNumber
		} else { // default use mobile phone
			ret = consts.MobileNumber
		}
	} else if regexp.MustCompile(`(email)`).MatchString(name) {
		ret = consts.Email
	} else if regexp.MustCompile(`(url)`).MatchString(name) {
		ret = consts.Url
	} else if regexp.MustCompile(`(ip)`).MatchString(name) {
		ret = consts.Ip
	} else if regexp.MustCompile(`(creditcard)`).MatchString(name) {
		ret = consts.CreditCard
	} else if regexp.MustCompile(`(idcard)`).MatchString(name) {
		ret = consts.IdCard
	} else if regexp.MustCompile(`(token)`).MatchString(name) {
		ret = consts.Token
	}

	return
}

func GetVarcharTypeByRecords(records []interface{}) (ret consts.VarcharType) {
	ret = consts.Empty

	if len(records) == 0 {
		return
	}

	val := fmt.Sprintf("%v", records[0])

	if govalidator.IsEmail(val) {
		ret = consts.Email
	} else if govalidator.IsCreditCard(val) {
		ret = consts.CreditCard
	} else if govalidator.IsMAC(val) {
		ret = consts.Mac
	} else if govalidator.IsUUID(val) {
		ret = consts.Uuid
	} else if govalidator.IsJSON(val) {
		ret = consts.JsonStr
		//} else if govalidator.IsUnixTime(val) {
		//	ret = consts.UnixTime
	} else if govalidator.IsMD5(val) {
		ret = consts.Md5

	} else if IsMobilePhone(val) {
		ret = consts.MobileNumber
	} else if IsTelPhone(val) {
		ret = consts.TelNumber
	} else if IsIDCard(val) {
		ret = consts.IdCard
	} else if govalidator.IsURL(val) {
		ret = consts.Url
	}

	return
}

type FieldTypeInfo struct {
	ColumnType  consts.ColumnType
	VarcharType consts.VarcharType

	IsNum     bool
	Precision int
	HasSign   bool

	Note string
	Rang string
}

func GenDefByColumnType(param string, ret *FieldTypeInfo) {
	rang := ""
	note := ""

	switch ret.ColumnType {
	// integer
	case "bit":
		rang = "0,1"
	case "tinyint":
		rang = "0-255"
	case "smallint":
		rang = "0-65535"
	case "mediumint":
		rang = "0-65535"
		note = `"mediumint [0,2^24-1]"`
	case "int":
		rang = "0-100000"
		note = `"ini [0,2^32-1]"`
	case "bigint":
		rang = "0-100000"
		note = `"bigint [0,2^64-1]"`
	// floating-point
	case "float":
		rang = "1.01-99.99:0.01"
		note = `"float"`
	case "double":
		rang = "1.01-99.99:0.01"
		note = `"double"`
	// fixed-point
	case "decimal":
		rang = "123.45"
		note = `"decimal"`
	// character string
	case "char":
		rang = genChar(param)
	case "tinytext":
		rang = `"tinytext"`
	case "text":
		rang = `"text"`
	case "mediumtext":
		rang = `"mediumtext"`
	case "longtext":
		rang = `"longtext"`
	// binary data
	case "tinyblob":
		rang = `"tinyblob"`
	case "blob":
		rang = `"blob"`
	case "mediumblob":
		rang = `"mediumblob"`
	case "longblob":
		rang = `"longblob"`
	case "binary":
		rang = "binary"
	case "varbinary":
		rang = "varbinary"
	// date and time type
	case "date":
		rang = `"date"`
	case "time":
		rang = `"time"`
	case "year":
		rang = `"year"`
	case "datetime":
		rang = `"datetime"`
	case "timestamp":
		rang = `"timestamp"`
	// other type
	case "geometry":
		rang = `"geometry"`
	case "point":
		rang = `"point"`
	case "linestring":
		rang = `"linestring"`
	case "polygon":
		rang = `"polygon"`
	case "multipoint":
		rang = `"multipoint"`
	case "multilinestring":
		rang = `"multilinestring"`
	case "multipolygon":
		rang = `"multipolygon"`
	case "geometrycollection":
		rang = `"geometrycollection"`
	case "enum":
		rang = getEnumValue(param)
	case "set":
		rang = getSetValue(param)
	default:
	}

	ret.Rang = rang
	ret.Note = note

	return
}

func genChar(param string) (ret string) {
	rang := `a-z`

	paramInt, _ := strconv.Atoi(param)

	if paramInt > 0 {
		rang += fmt.Sprintf("{%d!}", paramInt)
	}

	return
}

func getEnumValue(param string) (ret string) {
	arr := strings.Split(param, ",")

	num := getRandNum(len(arr))
	ret = strings.Trim(arr[num], "'")

	return
}

func getSetValue(param string) (ret string) {
	//arr := strings.Split(param, ",")
	//ret = strings.Join(arr, ",")

	ret = strings.ReplaceAll(param, "'", "")

	return
}

func getRandNum(num int) (ret int) {
	ret = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)

	return
}
