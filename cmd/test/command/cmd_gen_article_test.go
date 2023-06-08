package main

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestGenerateArticleCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateArticleCmdSuite))
}

type GenerateArticleCmdSuite struct {
	suite.Suite
}

func (s *GenerateArticleCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateArticleCmd")
}
func (s *GenerateArticleCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateArticleCmd")
}
func (s *GenerateArticleCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateArticleCmdSuite) TestGenerateArticle(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFileArticleConfig}).
		SetOutput(consts.CommandTestFileArticleOut).
		Gen()

	t.Require().NotContains(out, "{人称代词}", "check generated data")
}
