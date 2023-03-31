package service

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/helper"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type TextService struct {
	RangeService *RangeService `inject:""`
	FileService  *FileService  `inject:""`
}

func (s *TextService) CreateFieldValuesFromText(field *domain.DefField) {
	ranges := strings.Split(strings.TrimSpace(field.Range), ",")

	for _, rang := range ranges {
		rang = strings.TrimSpace(rang)
		repeat, repeatTag, rangWithoutRepeat := s.RangeService.ParseRepeat(rang)

		// get file and step string
		sectionArr := strings.Split(rangWithoutRepeat, ":")
		file := sectionArr[0]
		stepStr := "1"
		if len(sectionArr) == 2 {
			stepStr = sectionArr[1]
		}

		// read from file
		list := make([]string, 0)
		realPath := s.FileService.ComputerReferFilePath(file, field)
		content, err := os.ReadFile(realPath)
		if err != nil {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", file+" - "+realPath))
			field.Values = append(field.Values, fmt.Sprintf("FILE_NOT_FOUND"))
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
		numbs := helper.GenerateItems(int64(0), int64(len(list)-1), int64(step), 0, rand, 1, "", 0)

		// gen data by index
		count := 0
		if repeatTag == "" {
			for _, numb := range numbs {
				item := list[numb.(int64)]
				if strings.TrimSpace(item) == "" { // ignore empty line
					continue
				}

				for i := 0; i < repeat; i++ {
					field.Values = append(field.Values, item)

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

					field.Values = append(field.Values, item)

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

	if len(field.Values) == 0 {
		field.Values = append(field.Values, "N/A")
	}
}
