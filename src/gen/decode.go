package gen

import (
	"encoding/json"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"strings"
)

const (

)

func Decode(defaultFile, configFile, fieldsToExportStr, input, output string) {
	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	vari.Def = LoadConfigDef(defaultFile, configFile, &fieldsToExport)
	vari.Res = LoadResDef(fieldsToExport)

	data := fileUtils.ReadFile(input)

	ret := []map[string]interface{}{}
	LinesToMap(data, fieldsToExport, &ret)
	jsonObj, _ := json.Marshal(ret)
	vari.JsonResp = string(jsonObj)

	logUtils.Screen(i118Utils.I118Prt.Sprintf("analyse_success", output ))
	if vari.RunMode != constant.RunModeServerRequest {
		fileUtils.WriteFile(output, vari.JsonResp)
	}
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

		rowMap := decodeOnLevel(line, fieldsToExport, vari.Def.Fields)
		*ret = append(*ret, rowMap)
	}
	return
}

func decodeOnLevel(line string, fieldsToExport []string, fields []model.DefField) (ret map[string]interface{}) {
	ret = map[string]interface{}{}
	left := []rune(line)

	for j, field := range fields {
		col := ""
		if field.Width > 0 {
			col = string(left[:field.Width])
			left = left[field.Width:]
		} else {
			sepStr := ""
			if j < len(vari.Def.Fields) - 1 {
				sepStr = field.Postfix + vari.Def.Fields[j+1].Prefix
			} else {
				sepStr = field.Postfix
			}
			sep := []rune(sepStr)

			if len(sep) > 0 {
				index := searchRune(left, sep)
				if index > -1 {
					col = string(left[: index+runewidth.StringWidth(field.Postfix)])
					left = left[index+runewidth.StringWidth(field.Postfix) :]
				}
			}
		}

		ret[fieldsToExport[j]] = col
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