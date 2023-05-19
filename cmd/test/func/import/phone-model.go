package main

import (
	"fmt"
	"strings"

	"github.com/easysoft/zendata/cmd/test/func/comm"
	"github.com/easysoft/zendata/cmd/test/func/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

func main() {
	filePath := "/Users/aaron/work/zentao/product/zd/行业数据/手机型号.txt"

	db := comm.GetDB()
	err := db.Exec(fmt.Sprintf(comm.TruncateTable, model.PhoneModel{}.TableName())).Error
	if err != nil {
		panic(err)
	}

	content := fileUtils.ReadFile(filePath)

	for _, line := range strings.Split(content, "\n") {
		arr := strings.Split(strings.TrimSpace(line), ",")

		brand := strings.TrimSpace(arr[0])
		model1 := strings.TrimSpace(arr[1])
		area := strings.TrimSpace(arr[2])
		brandName := strings.TrimSpace(arr[3])
		modelName := strings.TrimSpace(arr[4])

		po := model.PhoneModel{
			Brand:     brand,
			Model:     model1,
			Area:      area,
			BrandName: brandName,
			ModelName: modelName,
		}

		db.Save(&po)
	}
}
