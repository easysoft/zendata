package helper

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func View(res string) {
	resPath, resType, sheet := fileUtils.GetResProp(res, vari.ZdDir)

	if resType == "yaml" {
		typ, inst, ranges := ReadYamlData(resPath)
		if typ == "inst" && inst.Instances != nil {
			printInst(inst)
		} else if typ == "range" && ranges.Ranges != nil {
			printRanges(ranges)
		} else { // common yaml file, print tips
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("print_common_yaml_file", res))
		}
	} else if resType == "excel" {
		printExcelSheet(resPath, sheet)
	}
}

func ReadYamlData(path string) (typ string, insts model.ResInstances, ranges model.ResRanges) {
	yamlContent, err := ioutil.ReadFile(path)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}

	err = yaml.Unmarshal(yamlContent, &insts)
	if err == nil && insts.Instances != nil && len(insts.Instances) > 0 {
		typ = "inst"
	} else {
		err = yaml.Unmarshal(yamlContent, &ranges)
		typ = "range"
	}

	return
}

func printInst(inst model.ResInstances) {
	msg := ""
	msg = msg + inst.Title + " " + inst.Desc + "\n"

	width := 0
	for _, item := range inst.Instances {
		lent := runewidth.StringWidth(item.Instance)
		if lent > width {
			width = lent
		}
	}

	for idx, item := range inst.Instances {
		if idx > 0 {
			msg = msg + "\n"
		}
		msg = msg + fmt.Sprintf("%d. %s - %s",
			idx+1, item.Instance+strings.Repeat(" ", width-runewidth.StringWidth(item.Instance)), item.Note)
	}

	logUtils.PrintTo(msg)
}

func printRanges(ranges model.ResRanges) {
	msg := ""
	msg = msg + ranges.Title + " " + ranges.Desc + "\n"

	width := 0
	for name, _ := range ranges.Ranges {
		lent := runewidth.StringWidth(name)
		if lent > width {
			width = lent
		}
	}

	i := 0
	for name, item := range ranges.Ranges {
		if i > 0 {
			msg = msg + "\n"
		}
		msg = msg + fmt.Sprintf("%d. %s - %s", i+1, name+strings.Repeat(" ", width-runewidth.StringWidth(name)), item)

		i++
	}

	logUtils.PrintTo(msg)
}

func printExcelSheet(path, sheetName string) {
	excel, err := excelize.OpenFile(path)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}

	msg := i118Utils.I118Prt.Sprintf("excel_data_1") + "\n"

	if sheetName == "" {
		for index, sheet := range excel.GetSheetList() {
			msg = msg + fmt.Sprintf("%d. %s", index+1, sheet) + "\n"
		}
	}

	logUtils.PrintTo(msg)

	if sheetName != "" {
		for _, sheet := range excel.GetSheetList() {
			if sheet != sheetName {
				continue
			}

			widthArr := make([]int, 0)
			dataArr := make([][]string, 0)
			dataArr = append(dataArr, make([]string, 0))
			rows, _ := excel.GetRows(sheet)

			colCount := 0
			for index, row := range rows {
				if index >= 10 {
					break
				}

				if index == 0 { // deal with the title
					for _, col := range rows[index] {
						val := strings.TrimSpace(col)
						if val == "" {
							break
						}

						widthArr = append(widthArr, runewidth.StringWidth(val))

						dataArr[0] = append(dataArr[0], val)
						colCount++
					}
				} else {
					colArr := make([]string, 0)
					for idx, col := range row {
						if idx >= colCount {
							break
						}

						val := strings.TrimSpace(col)

						lent := runewidth.StringWidth(val)
						if widthArr[idx] < lent {
							widthArr[idx] = lent
						}

						colArr = append(colArr, val)

					}
					dataArr = append(dataArr, colArr)
				}
			}

			for _, row := range dataArr {
				line := ""
				for colIdx, col := range row {
					if colIdx >= colCount {
						break
					}
					space := widthArr[colIdx] - runewidth.StringWidth(col)
					line = line + col + strings.Repeat(" ", space) + " "
				}

				logUtils.PrintTo(line)
			}
		}
	}
}
