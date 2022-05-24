package gen

import (
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func CreateFieldValuesFromText(field *model.DefField, fieldValue *model.FieldWithValues) {
	ranges := strings.Split(strings.TrimSpace(field.Range), ",")

	for _, rang := range ranges {
		rang = strings.TrimSpace(rang)
		repeat, repeatTag, rangWithoutRepeat := ParseRepeat(rang)

		// get file and step string
		sectionArr := strings.Split(rangWithoutRepeat, ":")
		file := sectionArr[0]
		stepStr := "1"
		if len(sectionArr) == 2 {
			stepStr = sectionArr[1]
		}

		// read from file
		list := make([]string, 0)
		realPath := fileUtils.ComputerReferFilePath(file, field)
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
		numbs := GenerateIntItems(0, (int64)(len(list)-1), step, rand, 1, "")
		// gen data by index
		count := 0
		if repeatTag == "" {
			for _, numb := range numbs {
				item := list[numb.(int64)]
				if strings.TrimSpace(item) == "" { // ignore empty line
					continue
				}

				for i := 0; i < repeat; i++ {
					fieldValue.Values = append(fieldValue.Values, item)

					count++
					if count >= constant.MaxNumb {
						break
					}
				}

				count++
				if count >= constant.MaxNumb {
					break
				}
			}
		} else if repeatTag == "!" {
			for i := 0; i < repeat; i++ {
				for _, numb := range numbs {
					item := list[numb.(int64)]
					if strings.TrimSpace(item) == "" { // ignore empty line
						continue
					}

					fieldValue.Values = append(fieldValue.Values, item)

					count++
					if count >= constant.MaxNumb {
						break
					}
				}

				count++
				if count >= constant.MaxNumb {
					break
				}
			}
		}
	}

	if len(fieldValue.Values) == 0 {
		fieldValue.Values = append(fieldValue.Values, "N/A")
	}
}
