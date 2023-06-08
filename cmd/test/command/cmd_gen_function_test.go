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

func (s *GenerateFunctionCmdSuite) TestGenerateFuncMd5(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f51").
		Gen()

	arr := strings.Split(out, "\t\n")

	t.Require().Equal(len(arr[0]), 32, "check generated data")
}

func (s *GenerateFunctionCmdSuite) TestGenerateFuncSha1(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f52").
		Gen()

	arr := strings.Split(out, "\t\n")

	t.Require().Equal(len(arr[0]), 40, "check generated data")
}

func (s *GenerateFunctionCmdSuite) TestGenerateFuncBase64(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f53").
		Gen()

	arr := strings.Split(out, "\t\n")

	t.Require().Equal(len(arr[0]), 36, "check generated data")
}

func (s *GenerateFunctionCmdSuite) TestGenerateFuncUrlencode(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTotal(5).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f54").
		Gen()

	arr := strings.Split(out, "\t\n")

	t.Require().Equal(len(arr[0]), 39, "check generated data")
}
