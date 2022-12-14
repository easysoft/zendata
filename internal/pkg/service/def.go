package service

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"strings"
	"time"
)

type DefService struct {
	ResService     *ResService     `inject:""`
	FieldService   *FieldService   `inject:""`
	CombineService *CombineService `inject:""`
}

func (s *DefService) GenerateFromContent(files []string, fieldsToExportStr, format, table string) {
	startTime := time.Now().Unix()
	count := 0

	if files[0] != "" {
		vari.GenVars.ConfigFileDir = fileUtils.GetAbsDir(files[0])
	} else {
		vari.GenVars.ConfigFileDir = fileUtils.GetAbsDir(files[1])
	}

	fieldsToExport := make([]string, 0)
	if fieldsToExportStr != "" {
		fieldsToExport = strings.Split(fieldsToExportStr, ",")
	}

	// get def and res data
	contents := gen.LoadFilesContents(files)
	vari.GenVars.DefData = gen.LoadDataContentDef(contents, &fieldsToExport)
	vari.GenVars.ResData = s.ResService.LoadResDef(fieldsToExport)

	if err := gen.CheckParams(); err != nil {
		return
	}
	gen.FixTotalNum()

	// gen for each field
	for i, field := range vari.GenVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, fieldsToExport) {
			continue
		}
		s.FieldService.Generate(&vari.GenVars.DefData.Fields[i], false)
	}

	// combine children fields
	for i, field := range vari.GenVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, fieldsToExport) {
			continue
		}
		s.CombineService.CombineChildrenIfNeeded(&vari.GenVars.DefData.Fields[i])
	}

	// print end msg
	entTime := time.Now().Unix()
	if vari.RunMode == constant.RunModeServerRequest {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", count, entTime-startTime))
	}
}

func (s *DefService) GenerateFromProtobuf(files []string) {
	startTime := time.Now().Unix()
	count := 0

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
