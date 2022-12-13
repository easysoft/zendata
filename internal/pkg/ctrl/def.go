package ctrl

import (
	"github.com/easysoft/zendata/internal/pkg/action"
	"github.com/easysoft/zendata/internal/pkg/service"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

type DefCtrl struct {
	DefService *service.DefService
}

func NewDefCtrl() *DefCtrl {
	return &DefCtrl{
		DefService: service.NewDefService(),
	}
}

func (c *DefCtrl) Generate(files []string, fieldsToExportStr, format, table string) {
	if len(files) == 0 {
		return
	}

	files = fileUtils.HandleFiles(files)

	if !action.IsFromProtobuf(files[0]) { // default gen from yaml
		c.DefService.GenerateFromContent(files, fieldsToExportStr, format, table)

	} else { // gen from protobuf
		c.DefService.GenerateFromProtobuf(files)

	}
}
