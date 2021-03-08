package gen

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func CreateFieldValuesFromText(field *model.DefField, fieldValue *model.FieldWithValues) {
	// get file and step string
	rang := strings.TrimSpace(field.Range)
	sectionArr := strings.Split(rang, ":")
	file := sectionArr[0]
	stepStr := "1"
	if len(sectionArr) == 2 {
		stepStr = sectionArr[1]
	}

	// read from file
	list := make([]string, 0)
	realPath := fileUtils.ComputerReferFilePath(file)
	content, err := ioutil.ReadFile(realPath)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", file+" - "+realPath))
		fieldValue.Values = append(fieldValue.Values, fmt.Sprintf("FILE_NOT_FOUND"))
		return
	}

	str := string(content)
	re := regexp.MustCompile(`\r?\t?\n`)
	str = re.ReplaceAllString(str, "\n")
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
	numbs := GenerateIntItems(0, (int64)(len(list)-1), step, rand, 1)
	// get data by index
	index := 0
	for _, numb := range numbs {
		item := list[numb.(int64)]

		if index >= constant.MaxNumb {
			break
		}
		if strings.TrimSpace(item) == "" {
			continue
		}

		fieldValue.Values = append(fieldValue.Values, item)
		index = index + 1
	}

	if len(fieldValue.Values) == 0 {
		fieldValue.Values = append(fieldValue.Values, "N/A")
	}
}
