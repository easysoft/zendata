package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

func CreateFieldValuesFromText(field *model.DefField, fieldValue *model.FieldWithValues) {
	// get file and step string
	rang := strings.TrimSpace(field.Range)
	sectionArr := strings.Split(rang, ":")
	file := sectionArr[0]
	stepStr := "1"
	if len(sectionArr) == 2 { stepStr = sectionArr[1] }

	// read from file
	list := make([]string, 0)
	realPath := findFilePath(file)
	content, err := ioutil.ReadFile(realPath)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", realPath))
		fieldValue.Values = append(fieldValue.Values, fmt.Sprintf("FILE_NOT_FOUND: %s", realPath))
		return
	}

	str := string(content)
	str = strings.Replace(str, "\\r\\n", "\\n", -1)
	str = stringUtils.TrimAll(str)
	list = strings.Split(str, "\n")

	// get step and rand
	rand := false
	step := 1
	if strings.ToLower(strings.TrimSpace(stepStr)) == "r" {
		rand = true
	} else {
		stepInt, err := strconv.Atoi(stepStr)
		if err == nil {
			step = stepInt
			if step < 0 {
				step = step * -1
			}
		}
	}

	// get index for data retrieve
	numbs := GenerateIntItems(0, (int64)(len(list) - 1), step, rand, 1)
	// get data by index
	index := 0
	for _, numb := range numbs {
		item := list[numb.(int64)]

		if index >= constant.MaxNumb { break }
		if strings.TrimSpace(item) == "" { continue }

		fieldValue.Values = append(fieldValue.Values, item)
		index = index + 1
	}

	if len(fieldValue.Values) == 0 {
		fieldValue.Values = append(fieldValue.Values, "N/A")
	}
}

func findFilePath(file string) string {
	resPath := file
	if !filepath.IsAbs(resPath) {

		resPath = vari.ConfigDir + file
		if !fileUtils.FileExist(resPath) {

			resPath = vari.DefaultDir + file
			if !fileUtils.FileExist(resPath) {

				resPath = vari.WorkDir + file
				if !fileUtils.FileExist(resPath) {
					resPath = ""
				}
			}
		}
	} else {
		if !fileUtils.FileExist(resPath) {
			resPath = ""
		}
	}

	return resPath
}