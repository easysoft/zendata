package controller

import (
	"bytes"
	"io"
	"path/filepath"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/service"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/kataras/iris/v12"
)

type DataCtrl struct {
	DecodeService *service.DecodeService `inject:""`
	MainService   *service.MainService   `inject:""`
	FileService   *service.FileService   `inject:""`
	BaseCtrl
}

func (c *DataCtrl) GenerateByFile(ctx iris.Context) {
	c.DealwithParams(ctx)

	defaultFile := ctx.URLParam("default")
	configFile := ctx.URLParam("config")

	vari.GlobalVars.ExportChildField = ""
	vari.GlobalVars.OutputFormat = ctx.URLParamDefault("format", "json")
	if vari.GlobalVars.OutputFormat == "text" {
		vari.GlobalVars.OutputFormat = "txt"
	}

	vari.GlobalVars.ExportFields = nil
	field := ctx.URLParamDefault("field", "")
	if field != "" {
		if strings.Contains(field, "~~") {
			vari.GlobalVars.ExportChildField = field
		} else {
			vari.GlobalVars.ExportFields = strings.Split(field, ",")
		}
	}

	//root := ctx.URLParam("root")
	//if root != "" {
	//	configUtils.UpdateRootDir(root)
	//
	//	if defaultFile != "" {
	//		defaultFile = filepath.Join(root, defaultFile)
	//	}
	//
	//	if configFile != "" {
	//		configFile = filepath.Join(root, configFile)
	//	}
	//}

	if defaultFile != "" {
		defaultFile = filepath.Join(vari.WorkDir, defaultFile)
	}
	if configFile != "" {
		configFile = filepath.Join(vari.WorkDir, configFile)
	}

	vari.GlobalVars.DefData = domain.DefData{}

	if defaultFile != "" {
		vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(defaultFile)
	} else {
		vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(configFile)
	}

	defaultContent := c.GetDistFileContent(defaultFile)
	configContent := c.GetDistFileContent(configFile)

	contents := [][]byte{defaultContent, configContent}
	contents = c.FileService.HandleFileBuffers(contents)

	c.MainService.GenerateDataByContents(contents)
	c.MainService.PrintOutput()
}

func (c *DataCtrl) DealwithParams(ctx iris.Context) {
	vari.GlobalVars.RunMode = consts.RunModeServerRequest
	logUtils.OutputHttpWriter = ctx.ResponseWriter()

	vari.GlobalVars.Total, _ = ctx.URLParamInt("lines")
	vari.GlobalVars.Trim, _ = ctx.URLParamBool("trim")
	vari.GlobalVars.Table = ctx.URLParam("table")
	vari.GlobalVars.Output = ctx.URLParam("outputFile")
	vari.GlobalVars.OutputFormat = ctx.URLParam("format")

	fields := strings.TrimSpace(ctx.URLParam("fields"))

	if fields != "" {
		vari.GlobalVars.ExportFields = strings.Split(fields, ",")
	}

	if vari.GlobalVars.OutputFormat == "" {
		vari.GlobalVars.OutputFormat = consts.FormatJson
	}

	return
}

func (c *DataCtrl) GetDistFileContent(file string) (ret []byte) {
	if fileUtils.FileExist(file) {
		ret = fileUtils.ReadFileBuf(file)
	}

	return
}

func (c *DataCtrl) GetFormFileContent(ctx iris.Context, name string) (ret []byte) {
	postFile, _, _ := ctx.FormFile(name)
	if postFile != nil {
		defer postFile.Close()

		buf := bytes.NewBuffer(nil)
		io.Copy(buf, postFile)

		ret = buf.Bytes()
	}

	return
}
