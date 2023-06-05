package main

import (
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
}

func (s *HelpCmdSuite) TestProductListApi(t provider.T) {
	t.ID("0")

	//firstProductId := gjson.Get(string(bodyBytes), "products.0.id").Int()
	//t.Require().Greater(firstProductId, int64(0), "list product")
}
