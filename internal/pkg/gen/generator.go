package gen

import (
	"errors"
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func GenerateFromContent(fileContents [][]byte, fieldsToExport *[]string) (
	rows [][]string, colIsNumArr []bool, err error) {

	vari.Def = LoadDataContentDef(fileContents, fieldsToExport)

	if len(vari.Def.Fields) == 0 {
		err = errors.New("")
		return
	} else if vari.Def.Type == constant.DefTypeArticle && vari.Out == "" {
		errMsg := i118Utils.I118Prt.Sprintf("gen_article_must_has_out_param")
		logUtils.PrintErrMsg(errMsg)
		err = errors.New(errMsg)
		return
	}

	if vari.Total < 0 {
		if vari.Def.Type == constant.DefTypeArticle {
			vari.Total = 1
		} else {
			vari.Total = constant.DefaultNumber
		}
	}

	// 为被引用的资源生成数据
	vari.ResLoading = true // not to use placeholder when loading res
	vari.Res = LoadResDef(*fieldsToExport)
	vari.ResLoading = false

	// 迭代fields生成值列表
	topLevelFieldNameToValuesMap := map[string][]string{}
	for index, field := range vari.Def.Fields {
		if !stringUtils.StrInArr(field.Field, *fieldsToExport) {
			continue
		}

		if field.Use != "" && field.From == "" {
			field.From = vari.Def.From
		}
		values := GenerateForFieldRecursive(&field, true)

		vari.Def.Fields[index].Precision = field.Precision

		topLevelFieldNameToValuesMap[field.Field] = values
		colIsNumArr = append(colIsNumArr, field.IsNumb)
	}

	// 处理TOP级别数据
	arrOfArr := make([][]string, 0) // 2 dimension arr for child, [ [a,b,c], [1,2,3] ]
	for _, child := range vari.Def.Fields {
		if !stringUtils.StrInArr(child.Field, *fieldsToExport) {
			continue
		}

		childValues := topLevelFieldNameToValuesMap[child.Field]

		// is value expression
		if child.Value != "" {
			childValues = helper.GenExpressionValues(child, topLevelFieldNameToValuesMap, vari.TopFieldMap)
		}

		// select from excel with expr
		if helper.SelectExcelWithExpr(child) {
			selects := helper.ReplaceVariableValues(child.Select, topLevelFieldNameToValuesMap)
			wheres := helper.ReplaceVariableValues(child.Where, topLevelFieldNameToValuesMap)

			childValues = make([]string, 0)
			childMapValues := make([][]string, 0)
			for index, slct := range selects {
				temp := child
				temp.Select = slct
				temp.Where = wheres[index%len(wheres)]

				resFile, _, sheet := fileUtils.GetResProp(temp.From, temp.FileDir)

				//	问题描述：
				//	原代码为：`selectCount := vari.Toal / len(selects)`
				//	因为整除的向下取整，如果`len(selects)`为3，`vari.Total`为8，则`selectCount`为2
				//	对于每一个`selects`的元素来说，都只会查两个元素，这样加起来一共只有6个结果，
				//	导致另外两个结果只能通过重复查到的数据的方式补充。
				//	解决方案：
				//  将代码改为: `selectCount := vari.Total / len(selects) + 1`,以达到使用人员的真正想要的
				//	即查到足够的数量，而不是通过重复补齐
				selectCount := vari.Total/len(selects) + 1
				mp := generateFieldValuesFromExcel(resFile, sheet, &temp, selectCount) // re-generate values
				for _, items := range mp {
					childMapValues = append(childMapValues, items)
				}
			}
			for index := 0; len(childValues) < vari.Total; {
				for i, _ := range selects {
					childValues = append(childValues, childMapValues[i][index%len(childMapValues[i])])
				}
				index++
			}
		}

		arrOfArr = append(arrOfArr, childValues)
	}
	rows = putChildrenToArr(arrOfArr, vari.Recursive)

	return
}

func GenerateFromYaml(files []string, fieldsToExport *[]string) (
	rows [][]string, colIsNumArr []bool, err error) {

	vari.ConfigFileDir = fileUtils.GetAbsDir(files[0])

	contents := LoadFilesContents(files)
	rows, colIsNumArr, err = GenerateFromContent(contents, fieldsToExport)

	return
}

func GenerateForFieldRecursive(field *model.DefField, withFix bool) (values []string) {
	dealwithFixRange(field)

	if len(field.Fields) > 0 { // has child fields
		values = genValuesForChildFields(field, withFix)

	} else if len(field.Froms) > 0 { // refer to multi res
		values = genValuesForMultiRes(field, withFix)

	} else if field.From != "" && field.Type != constant.FieldTypeArticle { // refer to res
		values = genValuesForSingleRes(field)

	} else if field.Config != "" { // refer to config
		values = genValuesForConfig(field)

	} else { // leaf field
		values = GenerateValuesForField(field)
	}

	if field.Rand && field.Type != constant.FieldTypeArticle {
		values = randomValues(values)
	}

	return values
}

func GenerateValuesForField(field *model.DefField) []string {
	values := make([]string, 0)

	fieldWithValues := CreateField(field)

	computerLoop(field)
	indexOfRow := 0
	count := 0

	uniqueTotal := len(fieldWithValues.Values)

	if l := len(field.PostfixRange.Values); l > 0 {
		uniqueTotal *= l
	}
	if l := len(field.PrefixRange.Values); l > 0 {
		uniqueTotal *= l
	}

	for {
		// 2. random replacement
		isRandomAndLoopEnd := !vari.ResLoading && //  ignore rand in resource
			!(*field).ReferToAnotherYaml &&
			(*field).IsRand && (*field).LoopIndex > (*field).LoopEnd
		// isNotRandomAndValOver := !(*field).IsRand && indexOfRow >= len(fieldWithValues.Values)
		if count >= vari.Total || count >= uniqueTotal || isRandomAndLoopEnd {
			for _, v := range fieldWithValues.Values {
				str := fmt.Sprintf("%v", v)
				str = addFix(str, field, count, true)
				values = append(values, str)
			}
			break
		}

		// 处理格式、前后缀、loop等
		val := loopFieldValWithFix(field, fieldWithValues, &indexOfRow, count, true)
		values = append(values, val)

		count++

		if count >= vari.Total || count >= uniqueTotal {
			break
		}

		(*field).LoopIndex = (*field).LoopIndex + 1
		if (*field).LoopIndex > (*field).LoopEnd {
			(*field).LoopIndex = (*field).LoopStart
		}
	}

	return values
}
