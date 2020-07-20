package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"strconv"
	"strings"
)

const (
	md5Col = "CW"
)

func AddMd5(path string) {
	excel, err := excelize.OpenFile(path)
	if err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}

	for _, sheet := range excel.GetSheetList() {
		rows, _ := excel.GetRows(sheet)

		colCount := 0
		for index, row := range rows {
			if index == 0 { // deal with the title
				for _, col := range rows[index] {
					val := strings.TrimSpace(col)
					if val == "" { break }
					colCount++
				}
				continue
			}

			if row[0] == "" { // stop when finding a blank in first column
				break
			}

			str := ""
			for idx, col := range row {
				if idx >= colCount { break }

				val := strings.TrimSpace(col)
				str = str + val
			}
			md5Str := md5V(str)
			excel.SetCellValue(sheet, md5Col + strconv.Itoa(index + 1), md5Str)
		}
	}

	if err := excel.SaveAs(path); err != nil {
		logUtils.Screen(i118Utils.I118Prt.Sprintf("fail_to_write_file", path))
	}
}

func md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}