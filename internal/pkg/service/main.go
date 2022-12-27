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

	count, err := s.GenerateDataByFile(files)
	if err != nil {
		return
	}

	s.PrintOutput()

	// print end msg
	entTime := time.Now().Unix()
	if vari.GlobalVars.RunMode == consts.RunModeServerRequest {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", count, entTime-startTime))
	}
}

func (s *MainService) GenerateDataByFile(files []string) (count int, err error) {
	if files[0] != "" {
		vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(files[0])
	} else {
		vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(files[1])
	}

	// get def and res data
	contents := s.FileService.LoadFilesContents(files)
	count, err = s.GenerateDataByContents(contents)

	return
}

func (s *MainService) GenerateDataByContents(contents [][]byte) (count int, err error) {
	vari.GlobalVars.DefData = s.DefService.LoadDataContentDef(contents, &vari.GlobalVars.ExportFields)

	s.ParamService.FixTotalNum()

	s.ResService.LoadResDef(vari.GlobalVars.ExportFields)

	err = s.ParamService.CheckParams()
	if err != nil {
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

func (s *MainService) PrintOutput() {
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
	if vari.GlobalVars.RunMode == consts.RunModeServerRequest {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("server_response", count, entTime-startTime))
	}
}
