package main

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestHelpCmd(t *testing.T) {
	suite.RunSuite(t, new(HelpCmdSuite))
}

type HelpCmdSuite struct {
	suite.Suite
}

func (s *HelpCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("HelpCmd")
}
func (s *HelpCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *HelpCmdSuite) TestPrintSample(t provider.T) {
	t.ID("0")

	logUtils.PrintExample()

	out := consts.Buf.String()

	t.Require().Contains(out, "语法说明", "check sample content")
}

func (s *HelpCmdSuite) TestPrintUsage(t provider.T) {
	t.ID("0")

	logUtils.PrintUsage()

	out := consts.Buf.String()

	t.Require().Contains(out, "数据生成工具", "check usage content")
}
