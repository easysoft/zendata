package service

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"io/ioutil"
	"time"
)

type DefService struct {
	ResService      *ResService      `inject:""`
	FieldService    *FieldService    `inject:""`
	CombineService  *CombineService  `inject:""`
	OutputService   *OutputService   `inject:""`
	ProtobufService *ProtobufService `inject:""`
}

func (s *DefService) GenerateFromContents(files []string) {
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

func (s *DefService) GenerateData(files []string) (count int) {
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
		s.CombineService.CombineChildrenIfNeeded(&vari.GlobalVars.DefData.Fields[i])
	}

	return
}

func (s *DefService) GenerateFromProtobuf(files []string) {
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

func (s *DefService) LoadFilesContents(files []string) (contents [][]byte) {
	contents = make([][]byte, 0)
	for _, f := range files {
		if f == "" {
			continue
		}
		pathDefaultFile := fileUtils.GetAbsolutePath(f)
		if !fileUtils.FileExist(pathDefaultFile) {
			return
		}
		content, err := ioutil.ReadFile(pathDefaultFile)
		if err != nil {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_parse_file"), color.FgCyan)
			return
		}
		contents = append(contents, content)
	}

	return
}
