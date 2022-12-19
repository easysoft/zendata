package service

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"io/ioutil"
)

type FileService struct {
}

func (s *FileService) ComputerReferFilePath(file string, field *model.DefField) (resPath string) {
	resPath = file
	if fileUtils.IsAbsPath(resPath) && fileUtils.FileExist(resPath) {
		return
	}

	resPath = field.FileDir + file
	if fileUtils.FileExist(resPath) {
		return
	}

	resPath = vari.GlobalVars.ConfigFileDir + file
	if fileUtils.FileExist(resPath) {
		return
	}

	resPath = vari.ZdPath + constant.ResDirUsers + constant.PthSep + file
	if fileUtils.FileExist(resPath) {
		return
	}
	resPath = vari.ZdPath + constant.ResDirYaml + constant.PthSep + file
	if fileUtils.FileExist(resPath) {
		return
	}

	resPath = vari.ZdPath + file
	if fileUtils.FileExist(resPath) {
		return
	}

	return
}

func (s *FileService) LoadFilesContents(files []string) (contents [][]byte) {
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
