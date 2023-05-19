package main

import (
	"fmt"
	"strings"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/城市电话区号.xlsx"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.AreaCodeCity{}.TableName()))

	_, records := comm.GetExcelFirstSheet(filePath)

	for colIndex := 0; colIndex < len(records[0]); colIndex++ {
		state := ""

		for rowIndex, row := range records {
			if colIndex > len(row)-1 {
				break
			}

			content := strings.TrimSpace(row[colIndex])

			if rowIndex == 0 {
				state = content
				continue
			}

			if content == "" {
				break
			}

			arr := strings.Split(content, " ")

			if len(arr) < 2 {
				if len(arr) < 2 {
					continue
				}
			}

			po := model.AreaCodeCity{
				State: state,
				City:  arr[0],
				Code:  arr[1],
			}

			db.Save(&po)
		}
	}
}
