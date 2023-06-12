package service

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/mattn/go-runewidth"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type FixService struct {
}

func (s *FixService) AddFix(val interface{}, field *domain.DefField, count int, withFix bool) (ret interface{}) {
	prefix := s.getStrValueFromFixRange(field.PrefixRange, count)
	postfix := s.getStrValueFromFixRange(field.PostfixRange, count)
	divider := field.Divider
	length := field.Length

	ret = val

	if prefix == "" && postfix == "" && length == 0 && len(divider) == 0 {
		return
	}

	if length > runewidth.StringWidth(fmt.Sprintf("%v", ret)) {
		ret = helper.AddPad(fmt.Sprintf("%v", ret), *field)
	}

	if withFix && !vari.GlobalVars.Trim {
		ret = prefix + fmt.Sprintf("%v", ret) + postfix
	}

	if vari.GlobalVars.OutputFormat == consts.FormatText && !vari.GlobalVars.Trim {
		ret = fmt.Sprintf("%v", ret) + divider
	}

	if commonUtils.GetType(ret) == "string" {
		//ret = s.TrimIfFormatIsNotText(fmt.Sprintf("%v", ret))
	}

	return
}

func (s *FixService) getStrValueFromFixRange(rang *domain.Range, index int) string {
	if rang == nil || len(rang.Values) == 0 {
		return ""
	}

	idx := index % len(rang.Values)
	x := rang.Values[idx]

	return s.convPrefixVal2Str(x, "")
}

func (s *FixService) convPrefixVal2Str(val interface{}, format string) string {
	str := "n/a"
	success := false

	switch val.(type) {
	case int64:
		if format != "" {
			str, success = helper.FormatStr(format, val.(int64), 0)
		}
		if !success {
			str = strconv.FormatInt(val.(int64), 10)
		}
	case float64:
		precision := 0
		if format != "" {
			str, success = helper.FormatStr(format, val.(float64), precision)
		}
		if !success {
			str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
		}
	case byte:
		str = string(val.(byte))
		if format != "" {
			str, success = helper.FormatStr(format, str, 0)
		}
		if !success {
			str = string(val.(byte))
		}
	case string:
		str = val.(string)

		match, _ := regexp.MatchString("%[0-9]*d", format)
		if match {
			valInt, err := strconv.Atoi(str)
			if err == nil {
				str, success = helper.FormatStr(format, valInt, 0)
			}
		} else {
			str, success = helper.FormatStr(format, str, 0)
		}
	default:
	}

	if str == "n/a" {
		log.Println(str)
	}

	return str
}

func (s *FixService) TrimIfFormatIsNotText(str string) (ret string) {
	ret = str

	if vari.GlobalVars.OutputFormat != "" && vari.GlobalVars.OutputFormat != consts.FormatText {
		ret = strings.TrimSpace(ret)
	}

	return
}
