package gen

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func GenerateFieldValuesFromText(field *model.DefField, fieldValue *model.FieldValue) {
	// get file and step string
	rang := strings.TrimSpace(field.Range)
	sectionArr := strings.Split(rang, ":")
	file := sectionArr[0]
	stepStr := "1"
	if len(sectionArr) == 2 { stepStr = sectionArr[1] }

	// read frome
	list := make([]string, 0)
	relaPath := vari.InputDir + file
	content, err := ioutil.ReadFile(relaPath)
	if err != nil {
		logUtils.Screen("fail to read " + relaPath + ", try to use global config")

		relaPath = "def" + string(os.PathSeparator) + file
		content, err = ioutil.ReadFile(relaPath)
		if err != nil {
			logUtils.Screen("fail to read " + relaPath + ", will return")

			fieldValue.Values = append(fieldValue.Values, "N/A")
			return
		} else {
			logUtils.Screen("success to read " + relaPath)
		}
	} else {
		logUtils.Screen("success to read " + relaPath)
	}
	str := string(content)
	str = strings.Replace(str, "\\r\\n", "\\n", -1)
	list = strings.Split(str, "\n")

	// get step and rand
	rand := false
	step := 1
	if strings.ToLower(strings.TrimSpace(stepStr)) != "r" {
		stepInt, err := strconv.Atoi(stepStr)
		if err == nil {
			step = stepInt
		}
	} else {
		rand = true
	}

	// get index for data retrieve
	numbs := GenerateIntItems(0, (int64)(len(list) - 1), step, rand)
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