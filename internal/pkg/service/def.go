package service

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"time"
)

type DefService struct {
	ResService     *ResService     `inject:""`
	FieldService   *FieldService   `inject:""`
	CombineService *CombineService `inject:""`
	OutputService  *OutputService  `inject:""`
}

func (s *DefService) GenerateFromContent(files []string) {
	startTime := time.Now().Unix()
	count := 0

	if files[0] != "" {
		vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(files[0])
	} else {
		vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(files[1])
	}

	// get def and res data
	contents := gen.LoadFilesContents(files)
	vari.GlobalVars.DefData = gen.LoadDataContentDef(contents, &vari.GlobalVars.ExportFields)
	vari.GlobalVars.ResData = s.ResService.LoadResDef(vari.GlobalVars.ExportFields)

	if err := gen.CheckParams(); err != nil {
		return
	}
	gen.FixTotalNum()

	union := false
	if vari.GlobalVars.OutputFormat != consts.FormatJson {
		union = true
	}

	// gen for each field
	for i, field := range vari.GlobalVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, vari.GlobalVars.ExportFields) {
			continue
		}
		s.FieldService.Generate(&vari.GlobalVars.DefData.Fields[i], union)
	}

	// combine children fields
	for i, field := range vari.GlobalVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, vari.GlobalVars.ExportFields) {
			continue
		}
		s.CombineService.CombineChildrenIfNeeded(&vari.GlobalVars.DefData.Fields[i])
	}

	// get output
	if vari.GlobalVars.OutputFormat == consts.FormatJson {
		s.OutputService.GenJson(&vari.GlobalVars.DefData)
	} else {
		s.OutputService.GenRows(&vari.GlobalVars.DefData)
	}

	// print end msg
	entTime := time.Now().Unix()
	if vari.RunMode == consts.RunModeServerRequest {
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
	if vari.RunMode == consts.RunModeServerRequest {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", count, entTime-startTime))
	}
}
