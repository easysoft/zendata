package main

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"strings"
	"testing"
)

func TestGenerateFunctionCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateFunctionCmdSuite))
}

type GenerateFunctionCmdSuite struct {
	suite.Suite
}

func (s *GenerateFunctionCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateFunctionCmd")
}
func (s *GenerateFunctionCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateFunctionCmd")
}
func (s *GenerateFunctionCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateFunctionCmdSuite) TestGenerateFunction(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f5").
		Gen()

	arr := strings.Split(out, "\t\n")

	t.Require().Equal(len(arr[0]), 32, "check generated data")
}

func (s *GenerateFunctionCmdSuite) TestGenerateExpr(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile, consts.CommandTestFile2}).
		SetFields("f0,f6").
		Gen()

	t.Require().Contains(out, "3 * 2 = 6", "check generated data")
}

func (s *GenerateFunctionCmdSuite) TestGenerateDatetime(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f7").
		Gen()

	t.Require().Contains(out, "2013/", "check generated data")
}
