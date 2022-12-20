package service

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"time"
)

type MainService struct {
	ResService      *ResService      `inject:""`
	FieldService    *FieldService    `inject:""`
	CombineService  *CombineService  `inject:""`
	OutputService   *OutputService   `inject:""`
	ProtobufService *ProtobufService `inject:""`
	FileService     *FileService     `inject:""`
	DefService      *DefService      `inject:""`
	ParamService    *ParamService    `inject:""`
}

func (s *MainService) GenerateFromContents(files []string) {
	startTime := time.Now().Unix()

	count := s.GenerateData(files)

	// get output
	if vari.GlobalVars.OutputFormat == consts.FormatText { // text
		s.OutputService.GenText(false)
	} else if vari.GlobalVars.OutputFormat == consts.FormatJson { // json
		s.OutputService.GenJson()
	} else if vari.GlobalVars.OutputFormat == consts.FormatXml { // xml
		s.OutputService.GenXml()
	} else if vari.GlobalVars.OutputFormat == consts.FormatExcel || vari.GlobalVars.OutputFormat == consts.FormatExcel { // excel
		s.OutputService.GenExcel()
	} else if vari.GlobalVars.OutputFormat == consts.FormatSql { // excel
		s.OutputService.GenSql()
	}

	// print end msg
	entTime := time.Now().Unix()
	if vari.RunMode == consts.RunModeServerRequest {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", count, entTime-startTime))
	}
}

func (s *MainService) GenerateData(files []string) (count int) {
	if files[0] != "" {
		vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(files[0])
	} else {
		vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(files[1])
	}

	s.ParamService.FixTotalNum()

	// get def and res data
	contents := s.FileService.LoadFilesContents(files)
	vari.GlobalVars.DefData = s.DefService.LoadDataContentDef(contents, &vari.GlobalVars.ExportFields)
	s.ResService.LoadResDef(vari.GlobalVars.ExportFields)

	if err := s.ParamService.CheckParams(); err != nil {
		return
	}

	join := false
	if vari.GlobalVars.OutputFormat != consts.FormatJson {
		join = true
	}

	// gen for each field
	for i, field := range vari.GlobalVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, vari.GlobalVars.ExportFields) {
			continue
		}

		s.FieldService.Generate(&vari.GlobalVars.DefData.Fields[i], join)

		vari.GlobalVars.ColIsNumArr = append(vari.GlobalVars.ColIsNumArr, field.IsNumb)
	}

	// combine children fields
	for i, field := range vari.GlobalVars.DefData.Fields {
		if !stringUtils.StrInArr(field.Field, vari.GlobalVars.ExportFields) {
			continue
		}

		// combine child
		s.CombineService.CombineChildrenIfNeeded(&vari.GlobalVars.DefData.Fields[i], false)
	}

	return
}

func (s *MainService) GenerateFromProtobuf(files []string) {
	startTime := time.Now().Unix()
	count := 0

	buf, pth := s.ProtobufService.GenerateProtobuf(files[0])

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
