package service

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/mattn/go-runewidth"
	"regexp"
	"strconv"
	"strings"
)

type FixService struct {
}

func (s *FixService) AddFix(str string, field *model.DefField, count int, withFix bool) (ret string) {
	prefix := s.getStrValueFromFixRange(field.PrefixRange, count)
	postfix := s.getStrValueFromFixRange(field.PostfixRange, count)
	divider := field.Divider

	if field.Length > runewidth.StringWidth(str) {
		str = stringUtils.AddPad(str, *field)
	}
	if withFix && !vari.GlobalVars.Trim {
		str = prefix + str + postfix
	}
	if vari.GlobalVars.OutputFormat == consts.FormatText && !vari.GlobalVars.Trim {
		str += divider
	}

	ret = s.RemoveSpaceIfOutputNoTextFile(str)

	return
}

func (s *FixService) getStrValueFromFixRange(rang *model.Range, index int) string {
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
			str, success = stringUtils.FormatStr(format, val.(int64), 0)
		}
		if !success {
			str = strconv.FormatInt(val.(int64), 10)
		}
	case float64:
		precision := 0
		if format != "" {
			str, success = stringUtils.FormatStr(format, val.(float64), precision)
		}
		if !success {
			str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
		}
	case byte:
		str = string(val.(byte))
		if format != "" {
			str, success = stringUtils.FormatStr(format, str, 0)
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
				str, success = stringUtils.FormatStr(format, valInt, 0)
			}
		} else {
			str, success = stringUtils.FormatStr(format, str, 0)
		}
	default:
	}

	return str
}

func (s *FixService) RemoveSpaceIfOutputNoTextFile(str string) (ret string) {
	ret = str

	if vari.GlobalVars.OutputFormat != "" && vari.GlobalVars.OutputFormat != consts.FormatText {
		ret = strings.TrimSpace(ret)
	}

	return
}
