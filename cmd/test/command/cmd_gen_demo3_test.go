package main

import (
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestGenerateDemo3Cmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateDemo3CmdSuite))
}

type GenerateDemo3CmdSuite struct {
	suite.Suite
}

func (s *GenerateDemo3CmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateDemo3Cmd")
}
func (s *GenerateDemo3CmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateDemo3Cmd")
}
func (s *GenerateDemo3CmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateDemo3CmdSuite) TestGenerateNestedRange(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/test/nested-range.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "2;\t1|2,", "check generated data")
}

func (s *GenerateDemo3CmdSuite) TestGenerateNestedRes(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/test/nested-res.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "10.0.0.1/8\tmask is 8", "check generated data")
	t.Require().NotContains(out, "nil", "check generated data")
}
