package main

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestGenerateReferCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateReferCmdSuite))
}

type GenerateReferCmdSuite struct {
	suite.Suite
}

func (s *GenerateReferCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateReferCmd")
}
func (s *GenerateReferCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateReferCmd")
}
func (s *GenerateReferCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateReferCmdSuite) TestGenerateReferConfig(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile, consts.CommandTestFile2}).
		SetFields("f8").
		Gen()

	t.Require().Contains(out, "106", "check generated data")
}

func (s *GenerateReferCmdSuite) TestGenerateReferRange(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f9").
		Gen()

	t.Require().Contains(out, "\n2\t\n", "check generated data")
}

func (s *GenerateReferCmdSuite) TestGenerateReferInstance(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f10").
		Gen()

	t.Require().Contains(out, "192.168.0.1", "check generated data")
}

func (s *GenerateReferCmdSuite) TestGenerateReferExecl(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f11").
		Gen()

	t.Require().Contains(out, "青岛市", "check generated data")
}

func (s *GenerateReferCmdSuite) TestGenerateReferText(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f12").
		Gen()

	t.Require().Contains(out, "Jack", "check generated data")
}

func (s *GenerateReferCmdSuite) TestGenerateReferContent(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f13").
		Gen()

	t.Require().Contains(out, "c-3", "check generated data")
}

func (s *GenerateReferCmdSuite) TestGenerateReferMultiFrom(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{consts.CommandTestFile}).
		SetFields("f14").
		Gen()

	t.Require().Contains(out, "1.0.0.2", "check generated data")
}
