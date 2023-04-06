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
	if info.ColumnType == consts.Json {
		generateDefForJson(records, &info)
		return
	}

	GenDefByColumnType(param, &info)

	return
}

func generateDefForJson(records []interface{}, info *FieldTypeInfo) {
	rang := Json
	if len(records) > 0 {
		rang = fmt.Sprintf("%v", records[0])
	}

	info.Rang = rang
	info.Format = "json"
}

func generateDefByVarcharType(param string, info *FieldTypeInfo) {
	if info.Rang != "" || info.VarcharType == "" {
		return
	}

	if info.VarcharType == consts.Username {
		info.From = "name.enaccount.v1.yaml"
		info.Use = "common_underline"

	} else if info.VarcharType == consts.Email {
		info.From = "email.v1.yaml"
		info.Use = "western_with_esp"

	} else if info.VarcharType == consts.Url {
		info.From = "domain.domain.v1.yaml"
		info.Use = "letters_at_cn"
		info.Prefix = "https://"

	} else if info.VarcharType == consts.Ip {
		info.From = "ip.v1.yaml"
		info.Use = "privateB"

	} else if info.VarcharType == consts.Mac {
		info.Format = "mac()"

	} else if info.VarcharType == consts.CreditCard {
		info.Format = "credit_card('amex')"

	} else if info.VarcharType == consts.IdCard {
		info.Format = "id_card()"

	} else if info.VarcharType == consts.MobileNumber {
		info.From = "phone.v1.yaml"
		info.Use = "cellphone"

	} else if info.VarcharType == consts.TelNumber {
		info.From = "phone.v1.yaml"
		info.Use = "telephone_china"

	} else if info.VarcharType == consts.Token {
		info.Format = "token()"

	} else if info.VarcharType == consts.Uuid {
		info.From = "uuid.v1.yaml"
		info.Use = "length32_random"

	} else if info.VarcharType == consts.JsonStr {
		info.Format = "json()"

	} else if info.VarcharType == consts.Md5 {
		info.Rang = "1-100"
		info.Format = "md5"
	}

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
	} else if regexp.MustCompile(`(mac).*(address)`).MatchString(name) {
		ret = consts.Mac
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

	} else if govalidator.IsMD5(val) {
		ret = consts.Md5

	} else if IsIDCard(val) {
		ret = consts.IdCard

	} else if IsMobilePhone(val) {
		ret = consts.MobileNumber

	} else if IsTelPhone(val) {
		ret = consts.TelNumber

	} else if govalidator.IsURL(val) {
		ret = consts.Url

	} else if govalidator.IsJSON(val) {
		ret = consts.JsonStr
	}

	return
}

type FieldTypeInfo struct {
	ColumnType  consts.ColumnType
	VarcharType consts.VarcharType

	IsNum     bool
	Precision int
	HasSign   bool

	Note    string
	Rang    string
	Loop    string
	Loopfix string
	Type    string
	Format  string
	From    string
	Use     string
	Select  string
	Prefix  string
}

func GenDefByColumnType(param string, ret *FieldTypeInfo) {
	switch ret.ColumnType {
	// int
	case "bit":
		ret.Rang, ret.Note = GenBit()
	case "tinyint":
		ret.Rang, ret.Note = GenTinyint(ret.HasSign)
	case "smallint":
		ret.Rang, ret.Note = GenSmallint(ret.HasSign)
	case "mediumint":
		ret.Rang, ret.Note = GenMediumint(ret.HasSign)
	case "int":
		ret.Rang, ret.Note = GenInt(ret.HasSign)
	case "bigint":
		ret.Rang, ret.Note = GenBigint(ret.HasSign)

	// float
	case "float":
		ret.Rang, ret.Note = GenFloat(ret.HasSign)
	case "double":
		ret.Rang, ret.Note = GenDouble(ret.HasSign)
	// fixed-point
	case "decimal":
		ret.Rang, ret.Note = GenDecimal(ret.HasSign)

	// string
	case "char":
		ret.Rang, ret.Loop = GenChar(param)

	case "tinytext":
		ret.From = "idiom.v1.idiom"
		ret.Select = "word"
	case "text":
		ret.From = "xiehouyu.v1.xiehouyu"
		ret.Select = "riddle"
	case "mediumtext":
		ret.From = "joke.v1.joke"
		ret.Select = "content"
	case "longtext":
		ret.From = "song.v1.song"
		ret.Select = "lyric"

	// binary data
	case "tinyblob":
		ret.From, ret.Format = GenBin()
	case "blob":
		ret.From, ret.Format = GenBin()
	case "mediumblob":
		ret.From, ret.Format = GenBin()
	case "longblob":
		ret.From, ret.Format = GenBin()
	case "binary":
		ret.From, ret.Format = GenBin()
	case "varbinary":
		ret.From, ret.Format = GenBin()

	// date and time type
	case "date":
		ret.Rang, ret.Type, ret.Format = GenDate()
	case "time":
		ret.Rang, ret.Type, ret.Format = GenTime()
	case "year":
		ret.Rang, ret.Type, ret.Format = GenYear()
	case "datetime":
		ret.Rang, ret.Type, ret.Format = GenDatetime()
	case "timestamp":
		ret.Rang, ret.Type, ret.Format = GenTimestamp()

	// other type
	case "enum":
		ret.Rang = getEnumValue(param)
	case "set":
		ret.Rang, ret.Loop, ret.Loopfix = getSetValue(param)

	//case "geometry":
	//	ret.Rang = `"geometry"`
	//case "point":
	//	ret.Rang = `"point"`
	//case "linestring":
	//	ret.Rang = `"linestring"`
	//case "polygon":
	//	ret.Rang = `"polygon"`
	//case "multipoint":
	//	ret.Rang = `"multipoint"`
	//case "multilinestring":
	//	ret.Rang = `"multilinestring"`
	//case "multipolygon":
	//	ret.Rang = `"multipolygon"`
	//case "geometrycollection":
	//	ret.Rang = `"geometrycollection"`

	default:
	}

	return
}

func getEnumValue(param string) (ret string) {
	ret = strings.ReplaceAll(param, "'", "")

	return
}

func getSetValue(param string) (ret, loop, loopfix string) {
	ret = strings.ReplaceAll(param, "'", "")

	arr := strings.Split(param, ",")

	start := 1
	if len(arr) > 2 {
		start = 2
	}
	loop = fmt.Sprintf("%d-%d", start, len(arr))

	return
}

func getRandNum(num int) (ret int) {
	ret = rand.New(rand.NewSource(time.Now().UnixNano())).Intn(num)

	return
}
