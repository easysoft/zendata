package main

import (
	"github.com/easysoft/zendata/cmd/command/action"
	"github.com/easysoft/zendata/cmd/test/consts"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"path/filepath"
	"testing"
)

func TestParseTableCmd(t *testing.T) {
	suite.RunSuite(t, new(ParseTableCmdSuite))
}

type ParseTableCmdSuite struct {
	suite.Suite
}

func (s *ParseTableCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("ParseTableCmd")
}
func (s *ParseTableCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("ParseTableCmd")
}
func (s *ParseTableCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *ParseTableCmdSuite) TestParseSql(t provider.T) {
	t.ID("0")

	vari.GlobalVars.Output = filepath.Dir(consts.CommandTestFileTablesOut)

	action.GenYaml(consts.CommandTestFileTables)

	out := fileUtils.ReadFile(filepath.Join(consts.CommandTestFileTablesOut))

	t.Require().Contains(out, "from: keys.yaml", "check parse sql")
}

func (s *ParseTableCmdSuite) TestParseArticle(t provider.T) {
	t.ID("0")

	vari.GlobalVars.Output = filepath.Dir(consts.CommandTestFileArticleOut)

	action.GenYaml(consts.CommandTestFileArticle)

	out := fileUtils.ReadFile(filepath.Join(consts.CommandTestFileArticleOut))

	t.Require().Contains(out, "{时间}{人称代词}", "check parse article")
}
