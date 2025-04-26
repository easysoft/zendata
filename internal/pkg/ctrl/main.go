package ctrl

import (
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/service"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"path/filepath"
	"strings"
)

type MainCtrl struct {
	MainService         *service.MainService         `inject:""`
	FileService         *service.FileService         `inject:""`
	TableParseService   *service.TableParseService   `inject:""`
	MsTableParseService *service.MsTableParseService `inject:""`
	MockService         *service.MockService         `inject:""`
	ArticleService      *service.ArticleService      `inject:""`

	SqlParseService *service.SqlParseService `inject:""`
}

func (c *MainCtrl) Generate(files []string) {
	if len(files) == 0 {
		return
	}

	files = c.FileService.HandleFiles(files)

	if !helper.IsFromProtobuf(files[0]) { // default gen from yaml
		c.MainService.GenerateFromContents(files)

	} else { // gen from protobuf
		c.MainService.GenerateFromProtobuf(files)

	}
}

func (c *MainCtrl) GenYaml(input string) {
	if vari.GlobalVars.DBDsn != "" { // from db table
		if strings.Contains(vari.GlobalVars.DBDsn, "sqlserver") {
			c.MsTableParseService.GenYamlFromTable()
		} else {
			c.TableParseService.GenYamlFromTable()
		}
		return
	}

	ext := filepath.Ext(input)
	if ext == ".sql" { // from sql
		c.SqlParseService.GenYamlFromSql(input)
	} else if ext == ".txt" { // from article
		c.ArticleService.GenYamlFromArticle(input)
	}
}

func (c *MainCtrl) GenMock(input string) {
	if vari.GlobalVars.Output == "" {
		vari.GlobalVars.Output = fileUtils.GetFileOrFolderDir(input)
	}

	c.MockService.GenMockDef(input)
}
