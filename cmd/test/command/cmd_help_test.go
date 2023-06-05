package main

import (
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
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
	t.AddSubSuite("HelpCmd")
	vari.Config.Language = "zh"
}

func (s *HelpCmdSuite) TestPrintSample(t provider.T) {
	t.ID("0")

	logUtils.PrintExample()

	//firstProductId := gjson.Get(string(bodyBytes), "products.0.id").Int()
	//t.Require().Greater(firstProductId, int64(0), "list product")
}

func (s *HelpCmdSuite) TestPrintUsage(t provider.T) {
	t.ID("0")

	logUtils.PrintUsage()

	//firstProductId := gjson.Get(string(bodyBytes), "products.0.id").Int()
	//t.Require().Greater(firstProductId, int64(0), "list product")
}
