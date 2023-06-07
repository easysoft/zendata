package main

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestGenerateRecursiveCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateRecursiveCmdSuite))
}

type GenerateRecursiveCmdSuite struct {
	suite.Suite
}

func (s *GenerateRecursiveCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()

	t.AddSubSuite("GenerateRecursiveCmd")
}
func (s *GenerateRecursiveCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()

	t.AddSubSuite("GenerateRecursiveCmd")
}
func (s *GenerateRecursiveCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateRecursiveCmdSuite) TestGenerateRecursiveChildren(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f3").
		Gen()

	t.Require().Contains(out, "3_C\t1_C\t", "check generation")
}

func (s *GenerateRecursiveCmdSuite) TestGenerateRecursiveRow(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f1,f2").
		Gen()

	t.Require().Contains(out, "[1]\t456", "check generation")
}
