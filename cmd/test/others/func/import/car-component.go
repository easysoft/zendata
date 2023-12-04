package main

import (
	"fmt"
	"strings"

	"github.com/easysoft/zendata/cmd/test/others/func/comm"
	"github.com/easysoft/zendata/cmd/test/others/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/汽车零件.txt"

	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, model.CarComponent{}.TableName())).Error
	if err != nil {
		panic(err)
	}

	content := fileUtils.ReadFile(filePath)

	for _, line := range strings.Split(content, "\n") {
		arr := strings.Split(strings.TrimSpace(line), " ")
		col1 := arr[0]

		po := model.CarComponent{
			DataComm: model.DataComm{
				Name: col1,
			},
		}

		db.Save(&po)
	}
}
