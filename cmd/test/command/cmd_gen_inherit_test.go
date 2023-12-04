package main

import (
	"testing"

	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestGenerateInheritCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateInheritCmdSuite))
}

type GenerateInheritCmdSuite struct {
	suite.Suite
}

func (s *GenerateInheritCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateInheritCmd")
}
func (s *GenerateInheritCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateInheritCmd")
}
func (s *GenerateInheritCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateInheritCmdSuite) TestGenerateInherit(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile, consts.CommandTestFile2}).
		SetFields("f0,f1").
		Gen()

	t.Require().Contains(out, "!3!\t!3!\t", "check generated data")
}
