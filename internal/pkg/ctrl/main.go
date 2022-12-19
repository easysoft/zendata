package ctrl

import (
	"github.com/easysoft/zendata/internal/pkg/action"
	"github.com/easysoft/zendata/internal/pkg/service"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
)

type DefCtrl struct {
	MainService *service.MainService `inject:""`
}

func (c *DefCtrl) Generate(files []string) {
	if len(files) == 0 {
		return
	}

	files = fileUtils.HandleFiles(files)

	if !action.IsFromProtobuf(files[0]) { // default gen from yaml
		c.MainService.GenerateFromContents(files)

	} else { // gen from protobuf
		c.MainService.GenerateFromProtobuf(files)

	}
}
