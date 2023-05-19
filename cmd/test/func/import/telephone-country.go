package main

import (
	"encoding/json"
	"fmt"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/国家区号.json"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.AreaCodeCountry{}.TableName()))

	content := fileUtils.ReadFileBuf(filePath)

	arr := make([]model.AreaCodeCountry, 0)

	json.Unmarshal(content, &arr)

	for _, item := range arr {
		db.Save(&item)
	}
}
