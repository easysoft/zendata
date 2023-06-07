package main

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestGenerateFormatCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateFormatCmdSuite))
}

type GenerateFormatCmdSuite struct {
	suite.Suite
}

func (s *GenerateFormatCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()

	t.AddSubSuite("GenerateFormatCmd")
}
func (s *GenerateFormatCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()

	t.AddSubSuite("GenerateFormatCmd")
}
func (s *GenerateFormatCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateFormatCmdSuite) TestGenerateFormat(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f4").
		Gen()

	t.Require().Contains(out, "passwd02", "check generation")
}

func (s *GenerateFormatCmdSuite) TestGenerateTrim(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetTrim(true).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f1").
		Gen()

	t.Require().Contains(out, "\n2\n", "check generation")
}

func (s *GenerateFormatCmdSuite) TestGenerateHuman(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetHuman(true).
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f1").
		Gen()

	t.Require().Contains(out, "[2]", "check generation")
}
