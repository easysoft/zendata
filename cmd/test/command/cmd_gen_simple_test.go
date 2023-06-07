package main

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestGenerateCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateCmdSuite))
}

type GenerateCmdSuite struct {
	suite.Suite
}

func (s *GenerateCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateCmd")
}
func (s *GenerateCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateCmd")
}
func (s *GenerateCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateCmdSuite) TestGenerate(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f2").
		Gen()

	t.Require().Contains(out, "123", "check generated data")
}
