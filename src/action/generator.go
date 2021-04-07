package action

import (
	"fmt"
	"github.com/easysoft/zendata/src/gen"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"os"
	"strings"
	"time"
)

func Generate(defaultFile string, configFile string, fieldsToExportStr, format, table string) (lines []interface{}) {
	startTime := time.Now().Unix()

	if defaultFile != "" && configFile == "" {
		configFile = defaultFile
		defaultFile = ""
	}

	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	rows, colIsNumArr, err := gen.GenerateOnTopLevel(defaultFile, configFile, &fieldsToExport)
	if err != nil {
		return
	}

	if format == constant.FormatExcel || format == constant.FormatCsv { // for excel and cvs
		gen.Write(rows, format, table, colIsNumArr, fieldsToExport)
	} else { // for preview and article writing
		lines = gen.Print(rows, format, table, colIsNumArr, fieldsToExport)
	}

	// article need to write to more than one files
	if format == constant.FormatText && vari.Def.Type == constant.ConfigTypeArticle {
		var filePath = logUtils.FileWriter.Name()
		defer logUtils.FileWriter.Close()

		for index, line := range lines {
			articlePath := fileUtils.GenArticleFiles(filePath, index)
			logUtils.FileWriter, _ = os.OpenFile(articlePath, os.O_RDWR|os.O_CREATE, 0777)
			fmt.Fprint(logUtils.FileWriter, line)
		}
	}

	entTime := time.Now().Unix()
	if vari.RunMode == constant.RunModeServerRequest {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", len(rows), entTime-startTime))
	}

	return
}
