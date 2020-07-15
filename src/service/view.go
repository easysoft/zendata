package service

import (
	"fmt"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func ViewRes(res string) {
	msg := ""

	resType, resPath := fileUtils.ConvertResPath(res)

	if resType == "yaml" {
		typ, inst, ranges := readYamlData(resPath)
		if typ == "inst" {
			printInst(inst)
		} else if typ == "range" {
			printRanges(ranges)
		}
	} else if resType == "excel" {
	}

	logUtils.Screen(msg)
}

func readYamlData(path string) (typ string, insts model.ResInsts, ranges model.ResRanges) {
	if strings.Index(path, "system") > -1 {
		path = vari.ExeDir + "data" + constant.PthSep + path
	} else {
		path = vari.ExeDir + constant.PthSep + path
	}

	yamlContent, err := ioutil.ReadFile(path)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
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

func readExcelSheet(path string) (str string) {
	//excel, err := excelize.OpenFile(path)
	//if err != nil {
	//	logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
	//	return
	//}
	//
	//for index, sheet := range excel.GetSheetList() {
	//	if index > 0 {
	//		title = title + "|"
	//	}
	//	title = title + sheet
	//}
	//
	//desc = i118Utils.I118Prt.Sprintf("excel_data")
	return
}

func readExcelData(path string) (str string) {
	//excel, err := excelize.OpenFile(path)
	//if err != nil {
	//	logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
	//	return
	//}
	//
	//for index, sheet := range excel.GetSheetList() {
	//	if index > 0 {
	//		title = title + "|"
	//	}
	//	title = title + sheet
	//}
	//
	//desc = i118Utils.I118Prt.Sprintf("excel_data")
	return
}

func printInst(inst model.ResInsts) {
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
		if idx > 0 { msg = msg + "\n" }
		msg = msg + fmt.Sprintf("%d. %s - %s",
			idx+1, item.Instance + strings.Repeat(" ", width -len(item.Instance)), item.Note)
	}

	logUtils.Screen(msg)
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
		if i > 0 { msg = msg + "\n" }
		msg = msg + fmt.Sprintf("%d. %s - %s", i+1, name + strings.Repeat(" ", width -len(name)), item)

		i++
	}

	logUtils.Screen(msg)

}