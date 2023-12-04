package comm

import (
	"fmt"
	"reflect"

	"github.com/easysoft/zendata/cmd/test/others/func/model"

	"github.com/easysoft/zendata/cmd/test/others/func/model"
)

func GetStructFields(interf interface{}) (fieldNames []string) {
	retH := reflect.TypeOf(interf)
	for i := 0; i < retH.NumField(); i++ {
		f := retH.Field(i)

		name := f.Name

		if name == "DataComm" {
			name = "Id"
		}

		if name != "CreatedAt" && name != "UpdatedAt" && name != "Deleted" && name != "Disabled" {
			fieldNames = append(fieldNames, name)
		}
	}

	return
}

func GetExcelColsByTableDef(infos []model.TableInfo) (excelColNameArr, excelColNameHeader []string) {
	excelColIndex := 'A'

	for _, info := range infos {
		if info.Field == "created_at" || info.Field == "updated_at" || info.Field == "deleted" || info.Field == "disabled" {
			continue
		}

		colName := fmt.Sprintf("%c", excelColIndex)
		excelColNameArr = append(excelColNameArr, colName)

		excelColNameHeader = append(excelColNameHeader, info.Field)

		excelColIndex += 1
	}

	return
}
