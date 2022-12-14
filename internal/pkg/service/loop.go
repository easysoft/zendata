package service

import (
	"errors"
	"github.com/easysoft/zendata/internal/pkg/gen/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
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

func (s *LoopService) LoopFieldValues(field *model.DefField, withFix bool) {
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
		isNotRandomAndValOver := !(*field).IsRand && indexOfRow >= len(values)
		if count >= vari.GenVars.Total || isRandomAndLoopEnd || isNotRandomAndValOver {
			break
		}

		(*field).LoopIndex = (*field).LoopIndex + 1
		if (*field).LoopIndex > (*field).LoopEnd {
			(*field).LoopIndex = (*field).LoopStart
		}
	}

	return
}

func (s *LoopService) LoopFieldValueToSingleStr(field *model.DefField, indexOfRow *int, count int, withFix bool) (loopStr string) {
	for j := 0; j < (*field).LoopIndex; j++ {
		if loopStr != "" {
			loopStr = loopStr + field.Loopfix
		}

		str, err := s.getFieldValByIndex(*field, indexOfRow)
		if err != nil {
			str = "N/A"
		}
		loopStr += str

		*indexOfRow++
	}

	loopStr = s.FixService.AddFix(loopStr, field, count, withFix)

	return
}

func (s *LoopService) ComputerLoopTimes(field *model.DefField) {
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

func (s *LoopService) getFieldValByIndex(field model.DefField, index *int) (val string, err error) {
	// 叶节点
	if len(field.Values) == 0 {
		if helper.SelectExcelWithExpr(field) {
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
