package ctrl

import (
	"github.com/easysoft/zendata/internal/pkg/action"
	"github.com/easysoft/zendata/internal/pkg/service"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"path/filepath"
)

type MainCtrl struct {
	MainService       *service.MainService       `inject:""`
	FileService       *service.FileService       `inject:""`
	TableParseService *service.TableParseService `inject:""`
	MockService       *service.MockService       `inject:""`

	SqlParseService *service.SqlParseService `inject:""`
}

func (c *MainCtrl) Generate(files []string) {
	if len(files) == 0 {
		return
	}

	files = c.FileService.HandleFiles(files)

	if !action.IsFromProtobuf(files[0]) { // default gen from yaml
		c.MainService.GenerateFromContents(files)

	} else { // gen from protobuf
		c.MainService.GenerateFromProtobuf(files)

	}
}

func (c *MainCtrl) GenYaml(input string) {
	if vari.GlobalVars.DBDsn != "" { // from db table
		c.TableParseService.GenYamlFromTable()
		return
	}

	ext := filepath.Ext(input)
	if ext == ".sql" { // from sql
		c.SqlParseService.GenYamlFromSql(input)
	} else if ext == ".txt" { // from article
		action.GenYamlFromArticle(input)
	}
}

func (c *MainCtrl) GenMock(input string) {
	if vari.GlobalVars.Output == "" {
		vari.GlobalVars.Output = fileUtils.GetFileOrFolderDir(input)
	}

	c.MockService.GenMockDef(input)
}
