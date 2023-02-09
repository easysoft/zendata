package service

import (
	"errors"
	"fmt"
	"github.com/easysoft/zendata/internal/pkg/domain"
	genHelper "github.com/easysoft/zendata/internal/pkg/gen/helper"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"strconv"
	"strings"
)

type LoopService struct {
	FixService    *FixService    `inject:""`
	FormatService *FormatService `inject:""`
}

func (s *LoopService) LoopAndFixFieldValues(field *domain.DefField, withFix bool) {
	s.ComputerLoopTimes(field)

	values := make([]interface{}, 0)

	indexOfRow := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := s.LoopFieldValueToSingleStr(field, &indexOfRow, count, withFix)
		values = append(values, str)

		count++
		isRandomAndLoopEnd := (*field).IsRand && (*field).LoopIndex == (*field).LoopEnd
		isNotRandomAndValOver := !(*field).IsRand && indexOfRow >= len(field.Values)
		if count >= vari.GlobalVars.Total || isRandomAndLoopEnd || isNotRandomAndValOver {
			break
		}

		(*field).LoopIndex = (*field).LoopIndex + 1
		if (*field).LoopIndex > (*field).LoopEnd {
			(*field).LoopIndex = (*field).LoopStart
		}
	}

	field.Values = values

	return
}

func (s *LoopService) LoopFieldValueToSingleStr(field *domain.DefField, indexOfRow *int, count int, withFix bool) (
	ret interface{}) {

	if (*field).LoopIndex <= 1 && field.Loopfix == "" {
		ret, _ = s.getFieldValByIndex(*field, indexOfRow)
		ret = s.FixService.AddFix(ret, field, count, withFix)

		*indexOfRow++
		return
	}

	for j := 0; j < (*field).LoopIndex; j++ {
		if ret != nil && ret != "" {
			ret = fmt.Sprintf("%v", ret) + field.Loopfix
		}

		str, err := s.getFieldValByIndex(*field, indexOfRow)
		if err != nil {
			str = "N/A"
		}

		temp := fmt.Sprintf("%v", str)
		if ret == nil || ret == "" {
			ret = temp
		} else {
			ret = fmt.Sprintf("%v", ret) + temp
		}

		*indexOfRow++
	}

	ret = s.FixService.AddFix(ret, field, count, withFix)

	return
}

func (s *LoopService) ComputerLoopTimes(field *domain.DefField) {
	if (*field).LoopIndex != 0 {
		return
	}

	arr := strings.Split(field.Loop, "-")
	(*field).LoopStart, _ = strconv.Atoi(arr[0])
	if len(arr) > 1 {
		field.LoopEnd, _ = strconv.Atoi(arr[1])
	}

	if (*field).LoopStart == 0 {
		(*field).LoopStart = 1
	}
	if (*field).LoopEnd == 0 {
		(*field).LoopEnd = 1
	}

	(*field).LoopIndex = (*field).LoopStart
}

func (s *LoopService) getFieldValByIndex(field domain.DefField, index *int) (val interface{}, err error) {
	// 叶节点
	if len(field.Values) == 0 {
		if genHelper.IsSelectExcelWithExpr(field) {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_generate_field", field.Field), color.FgCyan)
			err = errors.New("")
		}
		return
	}

	idx := *index % len(field.Values)
	str := field.Values[idx]

	val = s.FormatService.GetFieldValStr(field, str)

	return
}
