package genHelper

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func ComputerReferFilePath(file string, field *domain.DefField) (resPath string) {
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

	resPath = vari.ZdDir + consts.ResDirUsers + consts.PthSep + file
	if fileUtils.FileExist(resPath) {
		return
	}
	resPath = vari.ZdDir + consts.ResDirYaml + consts.PthSep + file
	if fileUtils.FileExist(resPath) {
		return
	}

	resPath = vari.ZdDir + file
	if fileUtils.FileExist(resPath) {
		return
	}

	return
}
