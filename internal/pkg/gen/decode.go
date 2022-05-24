package gen

import (
	"encoding/json"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/mattn/go-runewidth"
	"os"
	"path/filepath"
	"strings"
)

const ()

func Decode(files []string, fieldsToExportStr, input string) {
	if vari.Out != "" {
		fileUtils.MkDirIfNeeded(filepath.Dir(vari.Out))
		fileUtils.RemoveExist(vari.Out)
		logUtils.FileWriter, _ = os.OpenFile(vari.Out, os.O_RDWR|os.O_CREATE, 0777)
		defer logUtils.FileWriter.Close()
	}

	vari.ConfigFileDir = fileUtils.GetAbsDir(files[0])

	vari.Total = 10

	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	vari.Def = LoadDataDef(files, &fieldsToExport)
	vari.Res = LoadResDef(fieldsToExport)

	data := fileUtils.ReadFile(input)

	var ret []map[string]interface{}
	linesToMap(data, fieldsToExport, &ret)
	jsonObj, _ := json.Marshal(ret)
	vari.JsonResp = string(jsonObj)

	//logUtils.PrintTo(i118Utils.I118Prt.Sprintf("analyse_success", output))
	logUtils.PrintLine(vari.JsonResp)
}

func linesToMap(str string, fieldsToExport []string, ret *[]map[string]interface{}) {
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
			len := field.Length + runewidth.StringWidth(field.Prefix) + runewidth.StringWidth(field.Postfix)

			col = string(left[:len])
			left = left[len:]
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

		if vari.Trim {
			col = strings.TrimLeft(col, field.Prefix)
			col = strings.TrimRight(col, field.Postfix)
		}
		(*rowMap)[field.Field] = col

		children := field.Fields
		if len(children) > 0 {
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
