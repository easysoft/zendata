package service

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"strings"
)

type CombineService struct {
	ExpressionService  *ExpressionService  `inject:""`
	LoopService        *LoopService        `inject:""`
	OutputService      *OutputService      `inject:""`
	PlaceholderService *PlaceholderService `inject:""`
}

func (s *CombineService) CombineChildrenIfNeeded(field *model.DefField, isOnTopLevel bool) {
	if len(field.Fields) == 0 || !field.Join {
		return
	}

	fieldNameToValuesMap := map[string][]interface{}{}
	fieldNameToFieldMap := map[string]model.DefField{}

	// 1. get values for child fields
	if len(field.Values) == 0 {
		for index, child := range field.Fields {
			if len(child.Fields) > 0 && len(child.Values) == 0 { // no need to do if already generated
				s.CombineChildrenIfNeeded(&(field.Fields[index]), false)
			}

			fieldNameToValuesMap[field.Fields[index].Field] = field.Fields[index].Values
			fieldNameToFieldMap[field.Fields[index].Field] = field.Fields[index]
		}
	}

	// 2. deal with expression
	arrByField := make([][]interface{}, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
	for i, child := range field.Fields {
		childValues := fieldNameToValuesMap[child.Field]

		if child.Value != "" {
			childValues = s.ExpressionService.GenExpressionValues(child, fieldNameToValuesMap, fieldNameToFieldMap)
		}

		arrByField = append(arrByField, childValues)

		// clear child values after combined
		field.Fields[i].Values = nil
	}

	// 3. get combined values for parent field
	isRecursive := vari.GlobalVars.Recursive
	if stringUtils.InArray(field.Mode, constant.Modes) { // set on field level
		isRecursive = field.Mode == constant.ModeRecursive || field.Mode == constant.ModeRecursiveShort
	}

	field.Values = s.combineChildrenValues(arrByField, isRecursive, isOnTopLevel)
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
		indexArr = s.getModArr(arrOfArr)
	}

	for i := 0; i < count; i++ {
		strArr := make([]interface{}, 0)
		for j := 0; j < len(arrOfArr); j++ {
			child := arrOfArr[j]

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

func (s *CombineService) getModArr(arrOfArr [][]interface{}) []int {
	indexArr := make([]int, 0)
	for _, _ = range arrOfArr {
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
		col := item.(string)

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
