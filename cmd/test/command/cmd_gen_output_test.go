package main

import (
	"testing"

	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestGenerateOutputCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateOutputCmdSuite))
}

type GenerateOutputCmdSuite struct {
	suite.Suite
}

func (s *GenerateOutputCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateOutputCmd")
}
func (s *GenerateOutputCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateOutputCmd")
}
func (s *GenerateOutputCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateOutputCmdSuite) TestGenerateText(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f1,f2").
		SetOutput("test/out/result.txt").
		Gen()

	t.Require().Contains(out, "[1]\t123\t\n", "check generated data")
}

func (s *GenerateOutputCmdSuite) TestGenerateJson(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f1,f2").
		SetOutput("test/out/result.json").
		SetTrim(true).
		Gen()

	t.Require().Contains(out, `"f2": "123"`, "check generated data")
}

func (s *GenerateOutputCmdSuite) TestGenerateXml(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f1,f2").
		SetOutput("test/out/result.xml").
		SetTrim(true).
		Gen()

	t.Require().Contains(out, "<f1>1</f1>", "check generated data")
}

func (s *GenerateOutputCmdSuite) TestGenerateSql(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f1,f2").
		SetDBTable("mysql").
		SetTable("user").
		SetOutput("test/out/result.sql").
		SetTrim(true).
		Gen()

	t.Require().Contains(out, "('1','123')", "check generated data")
}

func (s *GenerateOutputCmdSuite) TestGenerateCsv(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f1,f2").
		SetOutput("test/out/result.csv").
		SetTrim(true).
		Gen()

	t.Require().Contains(out, "1,123\n", "check generated data")
}

func (s *GenerateOutputCmdSuite) TestGenerateExcel(t provider.T) {
	t.ID("0")

	pth := "test/out/result.xlsx"

	gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f1,f2").
		SetOutput(pth).
		SetTrim(true).
		Gen()

	value := testHelper.GetExcelData(pth, 0, 1, 1)

	t.Require().Contains(value, "123", "check generated data")
}

func (s *GenerateOutputCmdSuite) TestGenerateProtobuf(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFileProto}).
		SetProtoCls("Person").
		Gen()

	t.Require().Contains(out, "class Person", "check generated data")
}
