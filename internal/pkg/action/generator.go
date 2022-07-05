package action

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/gen/helper"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"path/filepath"
	"strings"
	"time"
)

func Generate(files []string, fieldsToExportStr, format, table string) {
	startTime := time.Now().Unix()
	if len(files) == 0 {
		return
	}

	count := 0
	files = fileUtils.HandleFiles(files)
	if !isFromProtobuf(files[0]) { // default gen from yaml
		if files[0] != "" {
			vari.ConfigFileDir = fileUtils.GetAbsDir(files[0])
		} else {
			vari.ConfigFileDir = fileUtils.GetAbsDir(files[1])
		}
		contents := gen.LoadFilesContents(files)
		GenerateByContent(contents, fieldsToExportStr, format, table)

	} else { // gen from protobuf
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

	var rows [][]string
	var colIsNumArr []bool
	var err error

	cacheKey, cacheOpt, hasCache := gen.ParseCache()
	if hasCache && cacheOpt != "new" { // retrieve from cache
		rows, colIsNumArr, err = gen.RetrieveCache(cacheKey, &fieldsToExport)
		vari.DefType = constant.DefTypeText

	} else { // gen new data and save
		rows, colIsNumArr, err = gen.GenerateFromContent(contents, &fieldsToExport)
		gen.CreateCache(cacheKey, fieldsToExport, rows, colIsNumArr)
	}

	if err != nil {
		return
	}

	if !isFromExcel(format) { // returned is for preview, sql exec and article writing
		lines = gen.Print(rows, format, table, colIsNumArr, fieldsToExport)
	} else { // for Excel and cvs
		gen.Write(rows, table, colIsNumArr, fieldsToExport)
	}

	// exec insert sql
	if vari.DBDsn != "" {
		helper.ExecSqlInUserDB(lines)
	}

	// article need to write to more than one files
	if isGenArticle(format) {
		gen.GenArticle(lines)
	}

	count = len(rows)

	entTime := time.Now().Unix()
	if vari.RunMode == constant.RunModeServerRequest {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", count, entTime-startTime))
	}

	return
}

func isFromProtobuf(file string) bool {
	return strings.ToLower(filepath.Ext(file)) == "."+constant.FormatProto
}

func isFromExcel(format string) bool {
	return format == constant.FormatExcel || format == constant.FormatCsv
}

func isGenArticle(format string) bool {
	return format == constant.FormatText && vari.Def.Type == constant.DefTypeArticle
}
