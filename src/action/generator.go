package action

import (
	"fmt"
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/gen/helper"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Generate(files []string, fieldsToExportStr, format, table string) (lines []interface{}) {
	startTime := time.Now().Unix()

	if len(files) == 0 {
		return
	}

	count := 0
	if strings.ToLower(filepath.Ext(files[1])) == "."+constant.FormatProto { //gen from protobuf
		buf, pth := gen.GenerateFromProtobuf(files[0])

		if vari.Verbose {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("protobuf_path", pth))
		}
		logUtils.PrintLine(buf)

		count = 1
		entTime := time.Now().Unix()
		if vari.RunMode == constant.RunModeServerRequest {
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", count, entTime-startTime))
		}
	} else { // default gen from yaml
		contents := gen.LoadFilesContents(files)
		lines = GenerateByContent(contents, fieldsToExportStr, format, table)
	}
	return
}

func GenerateByContent(contents [][]byte, fieldsToExportStr, format, table string) (lines []interface{}) {
	startTime := time.Now().Unix()

	if contents == nil {
		return
	}

	count := 0

	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	rows, colIsNumArr, err := gen.GenerateFromContent(contents, &fieldsToExport)
	if err != nil {
		return
	}
	if format == constant.FormatExcel || format == constant.FormatCsv { // for excel and cvs
		gen.Write(rows, table, colIsNumArr, fieldsToExport)
	} else { // returned is for preview, sql exec and article writing
		lines = gen.Print(rows, format, table, colIsNumArr, fieldsToExport)
	}

	// exec insert sql
	if vari.DBDsn != "" {
		helper.ExecSqlInUserDB(lines)
	}

	// article need to write to more than one files
	if format == constant.FormatText && vari.Def.Type == constant.ConfigTypeArticle {
		var filePath = logUtils.FileWriter.Name()
		defer logUtils.FileWriter.Close()
		fileUtils.RmFile(filePath)

		for index, line := range lines {
			articlePath := fileUtils.GenArticleFiles(filePath, index)
			fileWriter, _ := os.OpenFile(articlePath, os.O_RDWR|os.O_CREATE, 0777)
			fmt.Fprint(fileWriter, line)
			fileWriter.Close()
		}
	}

	count = len(rows)

	entTime := time.Now().Unix()
	if vari.RunMode == constant.RunModeServerRequest {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", count, entTime-startTime))
	}

	return
}
