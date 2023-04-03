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
		info.Format = "id_card()"

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

	Note   string
	Rang   string
	Type   string
	Format string
	From   string
	Use    string
	Select string
	Prefix string
}

func GenDefByColumnType(param string, ret *FieldTypeInfo) {
	rang := ""
	typ := ""
	format := ""
	from := ""
	selectStr := ""
	note := ""

	switch ret.ColumnType {
	// int
	case "bit":
		rang, note = GenBit()
	case "tinyint":
		rang, note = GenTinyint(ret.HasSign)
	case "smallint":
		rang, note = GenSmallint(ret.HasSign)
	case "mediumint":
		rang, note = GenMediumint(ret.HasSign)
	case "int":
		rang, note = GenInt(ret.HasSign)
	case "bigint":
		rang, note = GenBigint(ret.HasSign)

	// float
	case "float":
		rang, note = GenFloat(ret.HasSign)
	case "double":
		rang, note = GenDouble(ret.HasSign)
	// fixed-point
	case "decimal":
		rang, note = GenDecimal(ret.HasSign)

	// string
	case "char":
		rang = GenChar(param)
	case "tinytext":
		from = "idiom.v1.idiom"
		selectStr = "word"
	case "text":
		from = "xiehouyu.v1.xiehouyu"
		selectStr = "riddle"
	case "mediumtext":
		from = "joke.v1.joke"
		selectStr = "content"
	case "longtext":
		from = "song.v1.song"
		selectStr = "lyric"

	// binary data
	case "tinyblob":
		from, format = GenBin()
	case "blob":
		from, format = GenBin()
	case "mediumblob":
		from, format = GenBin()
	case "longblob":
		from, format = GenBin()
	case "binary":
		from, format = GenBin()
	case "varbinary":
		from, format = GenBin()

	// date and time type
	case "date":
		rang, typ, format = GenDate()
	case "time":
		rang, typ, format = GenTime()
	case "year":
		rang, typ, format = GenYear()
	case "datetime":
		rang, typ, format = GenDatetime()
	case "timestamp":
		rang, typ, format = GenTimestamp()

	// other type
	case "enum":
		rang = getEnumValue(param)
	case "set":
		rang = getSetValue(param)

	//case "geometry":
	//	rang = `"geometry"`
	//case "point":
	//	rang = `"point"`
	//case "linestring":
	//	rang = `"linestring"`
	//case "polygon":
	//	rang = `"polygon"`
	//case "multipoint":
	//	rang = `"multipoint"`
	//case "multilinestring":
	//	rang = `"multilinestring"`
	//case "multipolygon":
	//	rang = `"multipolygon"`
	//case "geometrycollection":
	//	rang = `"geometrycollection"`

	default:
	}

	ret.Rang = rang
	ret.Format = format
	ret.Type = typ
	ret.From = from
	ret.Select = selectStr
	ret.Note = note

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
