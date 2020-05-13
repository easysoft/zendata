package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"regexp"
	"strconv"
	"strings"
)

func GenerateForDefinition(total int, fieldsToExport string, out string, table string) ([][]string, []bool) {
	def := constant.Definition

	fieldsToExportArr := strings.Split(fieldsToExport, ",")
	fieldNameToValues := map[string][]string{}

	colTypes := make([]bool, 0)

	// 为每个field生成值列表
	for index, field := range def.Fields {
		if !stringUtils.FindInArr(field.Name, fieldsToExportArr) {
			continue
		}

		values := GenerateForField(&field, total)
		def.Fields[index].Precision = field.Precision

		fieldNameToValues[field.Name] = values
		colTypes = append(colTypes, field.IsNumb)
	}

	// 生成指定数量行的数据
	rows := make([][]string, 0)
	for i := 0; i < total; i++ {
		for _, field := range def.Fields {
			if !stringUtils.FindInArr(field.Name, fieldsToExportArr) {
				continue
			}

			values := fieldNameToValues[field.Name]
			fieldVal := values[i % len(values)]
			if len(rows) == i { rows = append(rows, make([]string, 0)) }
			rows[i] = append(rows[i], fieldVal)
		}
	}

	return rows, colTypes
}

func GenerateForField(field *model.Field,  total int) []string {
	convertFieldReferToNestedIfNeeded(field)

	values := make([]string, 0)

	if len(field.Fields) > 0 { // nested definition
		arr := make([][]string, 0)
		for _, child := range field.Fields {
			childValues := GenerateForField(&child, total)
			arr = append(arr, childValues)
		}

		for i := 0; i < total; i++ {
			concat := ""
			for _, row := range arr {
				concat = concat + row[i]
			}

			concat = field.Prefix + concat + field.Postfix
			values = append(values, concat)
		}
	} else if field.Type == "list" { // list type
		values = GenerateFieldItemsFromDefinition(field, total)

	} else if field.Type == "custom" { // custom type
		if field.Range != "" { // specific custom file
			LoadDefinitionFromFile(constant.ResDir + field.Range)
		}

		referField := constant.LoadedFields[field.Name]
		values = GenerateFieldItemsFromDefinition(&referField, total)

	} else {// other type like address.city
		arr := strings.Split(field.Type, ".")
		referField := constant.LoadedFields[arr[0]]
		referField.Filter = field.Range

		values = GenerateFieldItemsFromDefinition(&referField, total)
	}

	return values
}

func GenerateFieldItemsFromDefinition(field *model.Field, total int) []string {
	if field.Loop == 0 {field.Loop = 1}

	values := make([]string, 0)

	// 整理出值的列表
	datatype := strings.TrimSpace(field.Type)
	if datatype == "" { datatype = "list" }

	fieldValue := model.FieldValue{}

	//switch datatype {
	//case constant.LIST.String():
		fieldValue = GenerateList(field, total)
	//
	//default:
	//}

	index := 0
	count := 0
	for {
		// 处理格式、前后缀、loop等
		str := GenerateFieldValWithLoop(*field, fieldValue, &index)
		values = append(values, str)

		count++
		if count >= total {
			break
		}
	}

	return values
}

func GenerateFieldValWithLoop(field model.Field, fieldValue model.FieldValue, indexOfRow *int) string {
	prefix := field.Prefix
	postfix := field.Postfix

	loopStr := ""
	for j := 0; j < field.Loop; j++ {
		if loopStr != "" {
			loopStr = loopStr + field.Loopfix
		}

		str := GenerateFieldVal(field, fieldValue, indexOfRow)
		loopStr = loopStr + str

		*indexOfRow++
	}

	return prefix + loopStr + postfix
}

func GenerateFieldVal(field model.Field, fieldValue model.FieldValue, index *int) string {
	str := ""

	// 叶节点
	idx := *index % len(fieldValue.Values)
	val := fieldValue.Values[idx]
	str = GetFieldValStr(field, val)

	return str
}

func GetFieldValStr(field model.Field, val interface{}) string {
	str := "n/a"
	success := false

	switch val.(type) {
		case int64:
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, val.(int64))
			}
			if !success {
				str = strconv.FormatInt(val.(int64), 10)
			}
		case float64:
			precision := 0
			if field.Precision > 0 {
				precision = field.Precision
			}
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, val.(float64))
			}
			if !success {
				str = strconv.FormatFloat(val.(float64), 'f', precision, 64)
			}
		case byte:
			str = string(val.(byte))
			if field.Format != "" {
				str, success = stringUtils.FormatStr(field.Format, str)
			}
			if !success {
				str = string(val.(byte))
			}
		case string:
			str = val.(string)
		default:
	}

	return str
}

func convertFieldReferToNestedIfNeeded(field *model.Field) {
	// ${user_name}_${numb}@${domain}
	regx := regexp.MustCompile(`\$\{([a-zA-z0-9_]+)\}`)
	arrOfName := regx.FindAllStringSubmatch(field.Range, -1)

	if len(arrOfName) > 0 {
		strLeft := field.Range
		for index, a := range arrOfName {
			found := a[0]
			name := a[1]

			arr := strings.Split(strLeft, found)

			// add string constant
			if arr[0] != "" {
				strChild := model.Field{Name: "child-" + strconv.Itoa(index), Type: "list", Range: arr[0]}
				field.Fields = append(field.Fields, strChild)
			}

			child := model.Field{}
			if constant.LoadedFields[name].Name != "" {
				child = constant.LoadedFields[name]
			} else {
				child.Name = a[1]
				child.Type = "custom"
			}

			field.Fields = append(field.Fields, child)

			arr = arr[1:]
			strLeft = strings.Join(arr, "")

			if index == len(arrOfName) - 1 && strLeft != "" {
				// add string constant
				strChild := model.Field{Name: "child-" + strconv.Itoa(index), Type: "list", Range: strLeft}
				field.Fields = append(field.Fields, strChild)
			}
		}
	}
}