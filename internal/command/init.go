package command

import (
	"errors"
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/ctrl"
	"github.com/easysoft/zendata/internal/pkg/gen"
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

func InitCtrl() (defCtrl *ctrl.DefCtrl, err error) {
	defCtrl = &ctrl.DefCtrl{}

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

func ClearCache() (ret bool) {
	cacheKey, cacheOpt, _, hasCache, isBatch := gen.ParseCache()
	if cacheOpt == "clear" {
		if isBatch {
			gen.ClearBatchCache(cacheKey)
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("success_to_clear_cache", cacheKey))

		} else if cacheKey == "all" {
			gen.ClearAllCache()
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("success_to_clear_all_cache"))

		} else {
			if hasCache {
				gen.ClearCache(cacheKey)
			}
			logUtils.PrintTo(i118Utils.I118Prt.Sprintf("success_to_clear_cache", cacheKey))

		}

		ret = true
	}

	return
}

func SetOutFormat() (err error) {
	if vari.GlobalVars.OutputFile != "" {
		fileUtils.MkDirIfNeeded(filepath.Dir(vari.GlobalVars.OutputFile))
		fileUtils.RemoveExist(vari.GlobalVars.OutputFile)

		ext := strings.ToLower(filepath.Ext(vari.GlobalVars.OutputFile))
		ext = strings.TrimLeft(ext, ".")

		if stringUtils.InArray(ext, constant.Formats) {
			vari.GlobalVars.OutputFormat = ext
		}

		// create file writer
		if vari.GlobalVars.OutputFormat == constant.FormatExcel {
			logUtils.OutputFilePath = vari.GlobalVars.OutputFile
		} else {
			logUtils.OutputFileWriter, _ = os.OpenFile(vari.GlobalVars.OutputFile, os.O_RDWR|os.O_CREATE, 0777)
		}
	}

	if vari.GlobalVars.DBDsn != "" {
		vari.GlobalVars.OutputFormat = constant.FormatSql
	}

	if vari.GlobalVars.OutputFormat == constant.FormatSql && vari.GlobalVars.Table == "" {
		msg := i118Utils.I118Prt.Sprintf("miss_table_name")
		logUtils.PrintErrMsg(msg)
		err = errors.New(msg)
	}

	return
}
