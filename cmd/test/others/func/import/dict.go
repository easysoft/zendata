package main

import (
	"encoding/json"
	"fmt"

	"github.com/easysoft/zendata/cmd/test/others/func/comm"
	"github.com/easysoft/zendata/cmd/test/others/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/新华字典/data/word.json"

	db := comm.GetDB()
	db.Exec(fmt.Sprintf(comm.TruncateTable, model.DataColor{}.TableName()))
	db.AutoMigrate(
		&model.DataDict{},
	)

	content := fileUtils.ReadFileBuf(filePath)
	records := make([]model.DataDict, 0)
	json.Unmarshal(content, &records)

	for _, record := range records {
		db.Save(&record)
	}
}
