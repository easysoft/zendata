package service

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FormatService struct {
}

func (s *FormatService) GetFieldValStr(field model.DefField, val interface{}) string {
	str := "n/a"
	success := false

	format := strings.TrimSpace(field.Format)

	if field.Type == consts.FieldTypeTimestamp && format != "" {
		str = time.Unix(val.(int64), 0).Format(format)
		return str
	}

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
		if field.Precision > 0 {
			precision = field.Precision
		}
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

	return str
}
