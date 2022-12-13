package main

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/comm"
	"github.com/easysoft/zendata/cmd/test/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"strings"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/电脑文件扩展名.txt"

	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, model.PcFileExt{}.TableName())).Error
	if err != nil {
		panic(err)
	}

	content := fileUtils.ReadFile(filePath)

	for _, line := range strings.Split(content, "\n") {
		arr := strings.Split(strings.TrimSpace(line), " ")
		name := strings.TrimSpace(arr[0])
		desc := strings.TrimSpace(arr[1])

		po := model.PcFileExt{
			DataComm: model.DataComm{
				Name: name,
			},
			Desc: desc,
		}

		db.Save(&po)
	}
}
