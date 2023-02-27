package gen

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	genHelper "github.com/easysoft/zendata/internal/pkg/gen/helper"
	valueGen "github.com/easysoft/zendata/internal/pkg/gen/value"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func CreateFieldValuesFromText(field *domain.DefField, fieldValue *domain.FieldWithValues) {
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
		realPath := genHelper.ComputerReferFilePath(file, field)
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

		// get index list for data retrieve
		numbs := valueGen.GenerateItems(int64(0), int64(len(list)-1), int64(step), 0, rand, 1, "", 0)

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
					if count >= consts.MaxNumb {
						break
					}
				}

				count++
				if count >= consts.MaxNumb {
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
					if count >= consts.MaxNumb {
						break
					}
				}

				count++
				if count >= consts.MaxNumb {
					break
				}
			}
		}
	}

	if len(fieldValue.Values) == 0 {
		fieldValue.Values = append(fieldValue.Values, "N/A")
	}
}