package service

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"os"
	"strings"
)

type FileService struct {
}

func (s *FileService) ComputerReferFilePath(file string, field *domain.DefField) (resPath string) {
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

	resPath = vari.WorkDir + consts.ResDirUsers + consts.PthSep + file
	if fileUtils.FileExist(resPath) {
		return
	}
	resPath = vari.WorkDir + consts.ResDirYaml + consts.PthSep + file
	if fileUtils.FileExist(resPath) {
		return
	}

	resPath = vari.WorkDir + file
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
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_read_file", pathDefaultFile), color.FgCyan)
			return
		}
		content, err := os.ReadFile(pathDefaultFile)
		if err != nil {
			logUtils.PrintToWithColor(i118Utils.I118Prt.Sprintf("fail_to_parse_file", pathDefaultFile), color.FgCyan)
			return
		}
		contents = append(contents, content)
	}

	return
}

func (s *FileService) GetFilesFromParams(args []string) (files []string, count int) {
	for _, arg := range args {
		if strings.Index(arg, "-") != 0 {
			files = append(files, arg)
			count++
		} else {
			break
		}
	}

	return
}

func (s *FileService) HandleFiles(files []string) []string {
	if len(files) != 2 {
		return files
	}

	if files[0] == "" && files[1] != "" { // no defaultFile
		files[0] = files[1]
		files[1] = ""
	} else if files[1] == "" && files[0] != "" { // no configFile
		files[1] = files[0]
		files[0] = ""
	}

	return files
}

func (s *FileService) HandleFileBuffers(files [][]byte) [][]byte {
	if len(files) != 2 {
		return files
	}

	if files[0] == nil && files[1] != nil { // no defaultFile
		return [][]byte{files[1]}

	} else if files[1] == nil && files[0] != nil { // no configFile
		return [][]byte{files[0]}
	}

	return files
}
