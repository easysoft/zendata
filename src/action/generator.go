package action

import (
	"fmt"
	"github.com/easysoft/zendata/src/gen"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"strconv"
)

func Generate(file string, total int, fieldsToExport string, out string, table string) {
	gen.LoadDefinitionFromFile(file)

	rows := gen.GenerateForDefinition(total, fieldsToExport, out, table)
	Print(rows)
}

func Print(rows [][]string) {
	width := stringUtils.GetNumbWidth(len(rows))

	for i, cols := range rows {
		line := ""
		for j, col := range cols {
			if j >0 {
				line = line + ", "
			}
			line = line + col
		}

		idStr := fmt.Sprintf("%" + strconv.Itoa(width) + "d", i+1)
		logUtils.Screen(fmt.Sprintf("%s: %s", idStr, line))
	}
}