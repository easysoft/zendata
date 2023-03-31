package command

import (
	"errors"
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/ctrl"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/facebookgo/inject"
	"os"
	"path/filepath"
	"strings"
)

func InitCtrl() (defCtrl *ctrl.MainCtrl, err error) {
	defCtrl = &ctrl.MainCtrl{}

	var g inject.Graph

	if err := g.Provide(
		&inject.Object{Value: vari.DB},
		&inject.Object{Value: defCtrl},
	); err != nil {
		logUtils.PrintErrMsg(fmt.Sprintf("provide usecase objects to the Graph: %v", err))
	}
	err = g.Populate()
	if err != nil {
		logUtils.PrintErrMsg(fmt.Sprintf("populate the incomplete Objects: %v", err))
	}

	return
}

func SetOutFormat() (err error) {
	if vari.GlobalVars.Output != "" {
		fileUtils.MkDirIfNeeded(filepath.Dir(vari.GlobalVars.Output))
		fileUtils.RemoveExist(vari.GlobalVars.Output)

		ext := strings.ToLower(filepath.Ext(vari.GlobalVars.Output))
		ext = strings.TrimLeft(ext, ".")

		if stringUtils.InArray(ext, consts.Formats) {
			vari.GlobalVars.OutputFormat = ext
		}

		// create file writer
		if vari.GlobalVars.OutputFormat == consts.FormatExcel {
			logUtils.OutputFilePath = vari.GlobalVars.Output
		} else {
			logUtils.OutputFileWriter, _ = os.OpenFile(vari.GlobalVars.Output, os.O_RDWR|os.O_CREATE, 0777)
		}
	}

	if vari.GlobalVars.DBDsn != "" {
		vari.GlobalVars.OutputFormat = consts.FormatSql
	}

	if vari.GlobalVars.OutputFormat == consts.FormatSql && vari.GlobalVars.Table == "" {
		msg := i118Utils.I118Prt.Sprintf("miss_table_name")
		logUtils.PrintErrMsg(msg)
		err = errors.New(msg)
	}

	return
}
