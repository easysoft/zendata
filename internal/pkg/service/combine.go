package service

import (
	"fmt"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

type CombineService struct {
	ExcelService       *ExcelService       `inject:""`
	ExpressionService  *ExpressionService  `inject:""`
	LoopService        *LoopService        `inject:""`
	OutputService      *OutputService      `inject:""`
	PlaceholderService *PlaceholderService `inject:""`
}

func (s *CombineService) CombineChildrenIfNeeded(field *domain.DefField, isOnTopLevel bool) {
	if len(field.Fields) == 0 {
		return
	}

	// 1. get values for child fields
	if len(field.Values) == 0 {
		for index, child := range field.Fields {
			if len(child.Fields) > 0 && len(child.Values) == 0 { // no need to do if already generated
				s.CombineChildrenIfNeeded(&(field.Fields[index]), false)
			}

			// for text output only
			vari.GlobalVars.FieldNameToValuesMap[field.Fields[index].Field] = field.Fields[index].Values
			vari.GlobalVars.FieldNameToFieldMap[field.Fields[index].Field] = field.Fields[index]
		}
	}

	if !field.Join {
		return
	}

	// 2. deal with expression
	arrByField := make([][]interface{}, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
	for i, child := range field.Fields {
		if child.Value != "" {
			vari.GlobalVars.FieldNameToValuesMap[child.Field] = s.ExpressionService.GenExpressionValues(child)
		}

		// select from excel with expr
		if helper.IsSelectExcelWithExpr(child) {
			vari.GlobalVars.FieldNameToValuesMap[child.Field] = s.ExcelService.genExcelValuesWithExpr(&child, vari.GlobalVars.FieldNameToValuesMap)
		}

		arrByField = append(arrByField, vari.GlobalVars.FieldNameToValuesMap[child.Field])

		// clear child values after combined
		field.Fields[i].Values = nil
	}

	// 3. get combined values for parent field
	isRecursive := vari.GlobalVars.Recursive
	if stringUtils.InArray(field.Mode, consts.Modes) { // set on field level
		isRecursive = field.Mode == consts.ModeRecursive || field.Mode == consts.ModeRecursiveShort
	}

	if len(field.Values) == 0 && field.Fields != nil {
		field.Values = s.combineChildrenValues(arrByField, isRecursive, isOnTopLevel)
	}

	s.LoopService.LoopAndFixFieldValues(field, true)
}

func (s *CombineService) combineChildrenValues(arrByField [][]interface{}, isRecursive, isOnTopLevel bool) (ret []interface{}) {
	arrByRow := s.populateRowsFromTwoDimArr(arrByField, isRecursive, isOnTopLevel)

	for _, arr := range arrByRow {
		line := s.ConnectValues(arr)
		ret = append(ret, line)
	}

	return
}

func (s *CombineService) populateRowsFromTwoDimArr(arrOfArr [][]interface{}, isRecursive, isOnTopLevel bool) (
	values [][]interface{}) {
	count := vari.GlobalVars.Total

	if !isOnTopLevel {
		if isRecursive {
			count = s.getRecordCountForRecursive(arrOfArr)
		} else {
			count = s.getRecordCountForParallel(arrOfArr)
		}
	}

	indexArr := make([]int, 0)
	if isRecursive {
		indexArr = s.getModArrForChildrenRecursive(arrOfArr)
	}

	for i := 0; i < count; i++ {
		strArr := make([]interface{}, 0)
		for j := 0; j < len(arrOfArr); j++ {
			child := arrOfArr[j]
			if len(child) == 0 {
				continue
			}

			var index int
			if isRecursive {
				mod := indexArr[j]
				index = i / mod % len(child)
			} else {
				index = i % len(child)
			}

			val := child[index]
			strArr = append(strArr, val)
		}

		values = append(values, strArr)
	}

	return
}

func (s *CombineService) getRecordCountForParallel(arrOfArr [][]interface{}) int {
	// get max count of 2nd dim arr
	count := 1
	for _, arr := range arrOfArr {
		if count < len(arr) {
			count = len(arr)
		}
	}

	if count > vari.GlobalVars.Total {
		count = vari.GlobalVars.Total
	}

	return count
}

func (s *CombineService) getRecordCountForRecursive(arrOfArr [][]interface{}) int {
	count := 1
	for i := 0; i < len(arrOfArr); i++ {
		arr := arrOfArr[i]
		count = len(arr) * count
	}
	return count
}

func (s *CombineService) getModArrForChildrenRecursive(arrOfArr [][]interface{}) []int {
	indexArr := make([]int, 0)
	for range arrOfArr {
		indexArr = append(indexArr, 0)
	}

	for i := 0; i < len(arrOfArr); i++ {
		loop := 1
		for j := i + 1; j < len(arrOfArr); j++ {
			loop = loop * len(arrOfArr[j])
		}

		indexArr[i] = loop
	}

	return indexArr
}

func (s *CombineService) ConnectValues(values []interface{}) (ret string) {
	for i, item := range values {
		col := fmt.Sprintf("%v", item)

		if i > 0 && vari.GlobalVars.Human { // use a tab
			ret = strings.TrimRight(ret, "\t")
			col = strings.TrimLeft(col, "\t")

			ret += "\t" + col
		} else {
			ret += col
		}
	}

	return
}
