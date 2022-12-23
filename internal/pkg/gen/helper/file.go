package genHelper

import (
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
)

func ComputerReferFilePath(file string, field *model.DefField) (resPath string) {
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
