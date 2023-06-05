package main

import (
	"bytes"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"log"
	"os"
	"testing"
)

var (
	buf bytes.Buffer
)

func TestHelpCmd(t *testing.T) {
	suite.RunSuite(t, new(HelpCmdSuite))
}

type HelpCmdSuite struct {
	suite.Suite
}

func (s *HelpCmdSuite) BeforeEach(t provider.T) {
	t.AddSubSuite("HelpCmd")
	vari.Config.Language = "zh"

	log.SetOutput(&buf)
}
func (s *HelpCmdSuite) AfterEach(t provider.T) {
	buf.Reset()
	log.SetOutput(os.Stdout)
}

func (s *HelpCmdSuite) TestPrintSample(t provider.T) {
	t.ID("0")

	logUtils.PrintExample()

	out := buf.String()

	t.Require().Contains(out, "语法说明", "check sample content")
}

func (s *HelpCmdSuite) TestPrintUsage(t provider.T) {
	t.ID("0")

	logUtils.PrintUsage()

	out := buf.String()

	t.Require().Contains(out, "数据生成工具", "check usage content")
}
