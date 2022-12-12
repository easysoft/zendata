package main

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/comm"
	"github.com/easysoft/zendata/cmd/test/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"strings"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/电脑系统.txt"

	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, model.PcOs{}.TableName())).Error
	if err != nil {
		panic(err)
	}

	content := fileUtils.ReadFile(filePath)

	for _, line := range strings.Split(content, "\n") {
		arr := strings.Split(strings.TrimSpace(line), "|")
		shortName := strings.TrimSpace(arr[0])
		name := strings.TrimSpace(arr[1])
		version := strings.TrimSpace(arr[2])
		website := strings.TrimSpace(arr[3])

		po := model.PcOs{
			DataComm: model.DataComm{
				Name: name,
			},
			ShortName: shortName,
			Version:   version,
			Website:   website,
		}

		db.Save(&po)
	}
}
