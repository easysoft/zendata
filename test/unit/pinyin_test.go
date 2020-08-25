package main

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/Chain-Zhang/pinyin"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func TestPinYin(t *testing.T) {
	path := "../../data/name/cn.family.v1.xlsx"
	//path := "../../data/name/cn.given.v1.xlsx"
	excel, err := excelize.OpenFile(path)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}

	for _, sheet := range excel.GetSheetList() {
		rows, _ := excel.GetRows(sheet)

		for index, row := range rows {
			if index == 0 {
				continue
			}

			if row[0] == "" { // stop when finding a blank in first column
				break
			}

			name := strings.TrimSpace(row[1])
			pinyin, err := pinyin.New(name).Split(" ").Mode(pinyin.WithoutTone).Convert()
			if err == nil {
				t.Log(pinyin)
			} else {
				t.Error("fail to convert " + name)
			}

			pinyin = strings.Replace(pinyin," ", "", -1)
			excel.SetCellValue(sheet, "C" + strconv.Itoa(index + 1), pinyin)

			doub := "false"
			lent := ChineseCount(name)
			if lent > 1 {
				doub = "true"
			}
			excel.SetCellValue(sheet, "E" + strconv.Itoa(index + 1), doub)

		}
	}

	if err := excel.SaveAs(path); err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_write_file", path))
	}
}

func ChineseCount(str1 string) (count int) {
	for _, char := range str1{
		if unicode.Is(unicode.Han, char){
			count++
		}
	}

	return
}