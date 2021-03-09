package gen

import (
	"encoding/json"
	"github.com/easysoft/zendata/src/model"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"os"
	"path/filepath"
	"strings"
)

const ()

func Decode(defaultFile, configFile, fieldsToExportStr, input, output string) {
	if output != "" {
		fileUtils.MkDirIfNeeded(filepath.Dir(output))
		fileUtils.RemoveExist(output)
		logUtils.FileWriter, _ = os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0777)
		defer logUtils.FileWriter.Close()
	}

	vari.DefaultFileDir = fileUtils.GetAbsDir(defaultFile)
	vari.ConfigFileDir = fileUtils.GetAbsDir(configFile)

	vari.Total = 10

	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	vari.Def = LoadDataDef(defaultFile, configFile, &fieldsToExport)
	vari.Res = LoadResDef(fieldsToExport)

	data := fileUtils.ReadFile(input)

	var ret []map[string]interface{}
	LinesToMap(data, fieldsToExport, &ret)
	jsonObj, _ := json.Marshal(ret)
	vari.JsonResp = string(jsonObj)

	//logUtils.PrintTo(i118Utils.I118Prt.Sprintf("analyse_success", output))
	logUtils.PrintLine(vari.JsonResp)
}

func LinesToMap(str string, fieldsToExport []string, ret *[]map[string]interface{}) {
	start := 0
	if vari.WithHead {
		start = 1
	}

	for index, line := range strings.Split(str, "\n") {
		if index < start {
			continue
		}

		rowMap := map[string]interface{}{}
		decodeOneLevel(line, vari.Def.Fields, &rowMap)
		*ret = append(*ret, rowMap)
	}
	return
}

func decodeOneLevel(line string, fields []model.DefField, rowMap *map[string]interface{}) {

	left := []rune(line)

	for j, field := range fields {
		col := ""

		if field.Length > 0 {
			col = string(left[:field.Length])
			left = left[field.Length:]
		} else {
			sepStr := ""
			if j < len(fields)-1 {
				sepStr = field.Postfix + fields[j+1].Prefix
			} else {
				sepStr = field.Postfix
			}
			sep := []rune(sepStr)

			if len(sep) > 0 {
				index := searchRune(left, sep)
				if index > -1 {
					col = string(left[:index+len(field.Postfix)])
					left = left[index+len(field.Postfix):]
				}
			} else if j == len(fields)-1 {
				col = string(left)
				left = []rune{}
			}
		}

		(*rowMap)[field.Field] = col

		children := field.Fields
		if len(children) > 0 {
			//colWithoutFix := col[runewidth.StringWidth(field.Postfix):
			//	runewidth.StringWidth(col) - runewidth.StringWidth(field.Postfix)]

			rowMapChild := map[string]interface{}{}
			decodeOneLevel(col, children, &rowMapChild)

			(*rowMap)[field.Field+".fields"] = rowMapChild
		}
	}

	return
}

func searchRune(text []rune, what []rune) int {
	for i := range text {
		found := true
		for j := range what {
			if i+j < len(text) && text[i+j] != what[j] {
				found = false
				break
			}
		}
		if found {
			return i
		}
	}
	return -1
}
