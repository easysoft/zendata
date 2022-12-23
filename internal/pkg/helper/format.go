package helper

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/Chain-Zhang/pinyin"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/mattn/go-runewidth"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func FormatStr(format string, val interface{}, precision int) (ret string, pass bool) {
	if format == "" {
		return val.(string), true
	}

	if strings.Index(format, "md5") == 0 {
		str := interfaceToStr(val, precision)
		ret = Md5(str)
		pass = true
		return
	} else if strings.Index(format, "sha1") == 0 {
		str := interfaceToStr(val, precision)
		ret = Sha1(str)
		pass = true
		return
	} else if strings.Index(format, "base64") == 0 {
		str := interfaceToStr(val, precision)
		ret = Base64(str)
		pass = true
		return
	} else if strings.Index(format, "urlencode") == 0 {
		str := interfaceToStr(val, precision)
		ret = UrlEncode(str)
		pass = true
		return
	} else if strings.Index(format, "uuid") > -1 {
		ret = uuid.NewV4().String()
		sep := ""

		regx := regexp.MustCompile(`uuid\(\s*(\S+)\s*\)`)
		arr := regx.FindStringSubmatch(format)
		if len(arr) > 1 {
			sep = strings.Trim(arr[1], "'")
		}
		ret = strings.ReplaceAll(ret, "-", sep)

		pass = true
		return
	} else if strings.Index(format, "credit_card") > -1 {
		cardType := ""

		regx := regexp.MustCompile(`credit_card\(\s*(\S+)\s*\)`)
		arr := regx.FindStringSubmatch(format)
		if len(arr) > 1 {
			cardType = strings.Trim(arr[1], "'")
		}

		ret = GenerateCreditCard(cardType)

		pass = true
		return
	} else if strings.Index(format, "id_card") > -1 {
		ret = GenerateIdCard()

		pass = true
		return
	}

	str := fmt.Sprintf(format, val)
	if strings.Index(str, "%!") == 0 {
		return "", false
	}

	return str, true
}

func AddPad(str string, field model.DefField) string {
	if field.Length > 0 && field.Length > runewidth.StringWidth(str) {
		gap := field.Length - len(str)
		if field.LeftPad != "" {
			field.LeftPad = field.LeftPad[:1]
			pads := strings.Repeat(field.LeftPad, gap)
			str = pads + str
		} else if field.RightPad != "" {
			field.RightPad = field.RightPad[:1]
			pads := strings.Repeat(field.RightPad, gap)
			str = str + pads
		} else {
			field.LeftPad = " "
			pads := strings.Repeat(field.LeftPad, gap)
			str = pads + str
		}
	}

	return str
}

func ConvertForSql(str string) (ret string) {
	arr := []rune(str)

	count := 0
	for i, item := range arr {
		if count%2 == 1 && string(item) != "'" {
			ret = ret + "'"
		} else if i == len(arr)-1 && count%2 == 0 && string(item) == "'" {
			ret = ret + "'"
		}

		if string(item) != "'" {
			count = 0
		}

		ret = ret + string(item)

		if string(item) == "'" {
			count++
		}
	}

	return
}

func Escape(in string, shouldEscape []rune) (out string) {
	out = ""
	escapeChar := shouldEscape[0]

	for _, v := range in {
		for _, se := range shouldEscape {
			if v == se {
				out += string(escapeChar)
				break
			}
		}

		out += string(v)
	}

	return
}

func EscapeValueOfMysql(in string) string {
	return Escape(in, []rune{'\\', '\'', '"'})
}

func EscapeValueOfSqlServer(in string) string {

	return Escape(in, []rune{'\''})
}

func EscapeValueOfOracle(in string) string {

	return Escape(in, []rune{'\''})
}

func EscapeColumnOfMysql(in string) string {

	return Escape(in, []rune{'`'})
}

func EscapeColumnOfSqlServer(in string) string {
	return Escape(in, []rune{']'})
}

// oracle limit
//func EscapeColumnOfOracle(in string) string

func GetPinyin(word string) string {
	p, _ := pinyin.New(word).Split("").Mode(pinyin.WithoutTone).Convert()

	return p
}

func interfaceToStr(intf interface{}, precision int) (ret string) {
	switch intf.(type) {
	case int64:
		return strconv.FormatInt(intf.(int64), 10)
	case float64:
		return strconv.FormatFloat(intf.(float64), 'f', precision, 64)
	case byte:
		return string(intf.(byte))
	case string:
		return intf.(string)
	default:
		return intf.(string)
	}
}

func Md5(str string) (ret string) {
	h := md5.New()
	h.Write([]byte(str))
	ret = hex.EncodeToString(h.Sum(nil))

	return
}
func Sha1(str string) (ret string) {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	ret = fmt.Sprintf("%x", bs)

	return
}
func Base64(str string) (ret string) {
	ret = base64.StdEncoding.EncodeToString([]byte(str))

	return
}
func UrlEncode(str string) (ret string) {
	ret = url.QueryEscape(str)

	return
}

func ReplaceSpecialChars(bytes []byte) []byte {
	str := string(bytes)

	inRanges := false // for ranges yaml only
	ret := ""
	for _, line := range strings.Split(str, "\n") {
		if strings.Index(strings.TrimSpace(line), "ranges") == 0 {
			inRanges = true
		} else if len(line) > 0 && string(line[0]) != " " { // not begin with space, ranges end
			inRanges = false
		}

		if strings.Index(strings.TrimSpace(line), "range") == 0 || inRanges {
			line = strings.ReplaceAll(line, "[", string(constant.LeftBrackets))
			line = strings.ReplaceAll(line, "]", string(constant.RightBrackets))
		}

		ret += line + "\n"
	}

	return []byte(ret)
}

func ConvertYamlStringToMapFormat(bytes []byte) (ret string) {
	m := yaml.MapSlice{}
	yaml.Unmarshal(bytes, &m)
	bytesReturn, _ := yaml.Marshal(&m)
	ret = string(bytesReturn)

	// replace '"test"' to "test"
	reg := regexp.MustCompile(`([:\s]+?)'"(.*)"'`)
	//if reg.MatchString(ret) {
	ret = reg.ReplaceAllString(ret, `${1}"${2}"`)
	//}
	return
}

func ParseInt(str string) (ret int) {
	ret, _ = strconv.Atoi(str)

	return
}
func ParseBool(str string) (ret bool) {
	ret, _ = strconv.ParseBool(str)

	return
}
