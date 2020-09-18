package stringUtils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/Chain-Zhang/pinyin"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/mattn/go-runewidth"
	"strconv"
	"strings"
	"unicode"
)

func TrimAll(str string) string {
	str = strings.Trim(str, "\n")
	str = strings.TrimSpace(str)

	return str
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}
func UcAll(str string) string {
	ret := ""

	for _, v := range str {
		ret += string(unicode.ToUpper(v))
	}
	return ret
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func U2s(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func BoolToPass(b bool) string {
	if b {
		return constant.PASS.String()
	} else {
		return constant.FAIL.String()
	}
}

func FindInArr(str string, arr []string) (bool,int) {
	for index, s := range arr {
		if str == s {
			return true, index
		}
	}

	return false, -1
}

func StrInArr(str string, arr []string) bool {
	found, _ := FindInArr(str, arr)
	return found
}
func InArray(need interface{}, arr []string) bool {
	for _,v := range arr{
		if need == v{
			return true
		}
	}
	return false
}

func AddPostfix(str string, width int, ch string) string {
	lent := runewidth.StringWidth(str)

	if width > lent {
		postFix := strings.Repeat(ch, width-lent)
		str += postFix
	}

	return str
}

func GetNumbWidth(numb int) int {
	str := strconv.Itoa(numb)
	width := len(str)

	return width
}

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
	}

	str := fmt.Sprintf(format, val)
	if strings.Index(str,"%!") == 0 {
		return "", false
	}

	return str, true
}

func AddPad(str string, field model.DefField) string {
	if  field.Width > 0 && field.Width > runewidth.StringWidth(str) {
		gap := field.Width - len(str)
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

func StartWith(str, sub string) bool {
	return strings.Index(str, sub) == 0
}
func EndWith(str, sub string) bool {
	return strings.Contains(str, sub) && strings.LastIndex(str, sub) == len(str) -len(sub)
}

func ConvertForSql(str string) (ret string) {
	arr := []rune(str)

	count := 0
	for i, item := range arr {
		if count % 2 == 1 && string(item) != "'" {
				ret = ret + "'"
		} else if i == len(arr) - 1 && count % 2 == 0 && string(item) == "'" {
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
	ret = base64.URLEncoding.EncodeToString([]byte(str))

	return
}