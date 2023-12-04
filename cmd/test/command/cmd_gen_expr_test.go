package main

import (
	"strings"
	"testing"

	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestGenerateExpressionCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateExpressionCmdSuite))
}

type GenerateExpressionCmdSuite struct {
	suite.Suite
}

func (s *GenerateExpressionCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateExpressionCmd")
}
func (s *GenerateExpressionCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateExpressionCmd")
}
func (s *GenerateExpressionCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateExpressionCmdSuite) TestGenerateExpression(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f5").
		Gen()

	arr := strings.Split(out, "\t\n")

	t.Require().Equal(len(arr[0]), 32, "check generated data")
}

func (s *GenerateExpressionCmdSuite) TestGenerateExpr(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile, consts.CommandTestFile2}).
		SetFields("f0,f6").
		Gen()

	t.Require().Contains(out, "3 * 2 = 6", "check generated data")
}

func (s *GenerateExpressionCmdSuite) TestGenerateDatetime(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f7").
		Gen()

	t.Require().Contains(out, "2013/", "check generated data")
}
