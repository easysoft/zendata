package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

const (
	size = 4
)

func ListRes() {
	res := map[string][][size]string{}
	path := vari.ExeDir + "data"
	GetFilesAndDirs(path, &res)

	names := make([]string, 0)
	nameWidth := 0
	titleWidth := 0
	for key, arrOfArr := range res {
		for index, arr := range arrOfArr {
			path := arr[0]
			if key == "yaml" {
				arr[2], arr[3] = readYamlInfo(path)
			} else if key == "excel" {
				arr[2], arr[3] = readExcelInfo(path)
			}

			res[key][index] = arr
			name := pathToName(arr[1])
			names = append(names, name)
			lent := runewidth.StringWidth(name)
			if lent > nameWidth {
				nameWidth = lent
			}

			if key == "excel" {
				sheets := strings.Split(arr[2], "|")
				for _, sheet := range sheets {
					lent2 := runewidth.StringWidth(sheet)
					if lent2 > titleWidth {
						titleWidth = lent2
					}
				}
			} else {
				lent2 := runewidth.StringWidth(arr[2])
				if lent2 > titleWidth {
					titleWidth = lent2
				}
			}
		}
	}

	sysMsg := ""
	customMsg := ""
	idx := 0
	for _, arrOfArr := range res {
		for _, arr := range arrOfArr {
			name := names[idx]

			titleStr := arr[2]
			titles := strings.Split(titleStr, "|")

			idx2 := 0
			for _, title := range titles {
				if idx2 > 0 {
					name = ""
				}
				name = name + strings.Repeat(" ", nameWidth - runewidth.StringWidth(name))

				title = title  + strings.Repeat(" ", titleWidth - runewidth.StringWidth(title))
				msg := fmt.Sprintf("%s  %s  %s\n", name, title, arr[3])
				if strings.Index(name, "system") > -1 {
					sysMsg = sysMsg + msg
				} else {
					customMsg = sysMsg + msg
				}

				idx2++
			}

			idx++
		}
	}

	logUtils.PrintTo(sysMsg + "\n" + customMsg)
}

func GetFilesAndDirs(path string, res *map[string][][size]string)  {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, fi := range dir {
		if fi.IsDir() {
			GetFilesAndDirs(path + constant.PthSep + fi.Name(), res)
		} else {
			name := fi.Name()
			arr := [size]string{}
			if strings.HasSuffix(name, ".yaml") {
				arr[0] = path + constant.PthSep + name
				arr[1] = path[strings.LastIndex(path, "data"):] + constant.PthSep + name
				arr[1] = strings.Trim(arr[1], "data"+constant.PthSep)

				(*res)["yaml"] = append((*res)["yaml"], arr)
			} else if strings.HasSuffix(name, ".xlsx") {
				arr[0] = path + constant.PthSep + name
				arr[1] = path[strings.LastIndex(path, "data"):] + constant.PthSep + name
				arr[1] = strings.Trim(arr[1], "data"+constant.PthSep)

				(*res)["excel"] = append((*res)["excel"], arr)
			}
		}
	}
}

func readYamlInfo(path string) (title string, desc string) {
	configDef := model.DefData{}

	yamlContent, err := ioutil.ReadFile(path)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}
	err = yaml.Unmarshal(yamlContent, &configDef)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_file", path))
		return
	}

	title = configDef.Title
	desc = configDef.Desc
	return
}

func readExcelInfo(path string) (title string, desc string) {
	excel, err := excelize.OpenFile(path)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}

	for index, sheet := range excel.GetSheetList() {
		if index > 0 {
			title = title + "|"
		}
		title = title + sheet
	}

	desc = i118Utils.I118Prt.Sprintf("excel_data")
	return
}

func pathToName(path string) string {
	name := strings.ReplaceAll(path, constant.PthSep,".")
	name = name[:strings.LastIndex(name, ".")]

	return name
}