package comm

import (
	"fmt"
	"github.com/easysoft/zendata/cmd/test/model"
	"reflect"
)

func GetStructFields(interf interface{}) (fieldNames []string) {
	retH := reflect.TypeOf(interf)
	for i := 0; i < retH.NumField(); i++ {
		f := retH.Field(i)

		if f.Name != "CreatedAt" && f.Name != "UpdatedAt" && f.Name != "Deleted" && f.Name != "Disabled" {
			fieldNames = append(fieldNames, f.Name)
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
