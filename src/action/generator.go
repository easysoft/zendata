package action

import (
	"github.com/easysoft/zendata/src/gen"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"strings"
	"time"
)

func Generate(defaultFile string, configFile string, fieldsToExportStr, format, table string) {
	startTime := time.Now().Unix()

	if defaultFile != "" && configFile == "" {
		configFile = defaultFile
		defaultFile = ""
	}

	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	rows, colIsNumArr, err := gen.GenerateForOnTop(defaultFile, configFile, &fieldsToExport)
	if err != nil {
		return
	}
	gen.Print(rows, format, table, colIsNumArr, fieldsToExport)

	entTime := time.Now().Unix()
	if vari.RunMode == constant.RunModeServerRequest {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", len(rows), entTime - startTime))
	}
}