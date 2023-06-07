package main

import (
	"github.com/easysoft/zendata/cmd/command/action"
	"github.com/easysoft/zendata/cmd/test/consts"
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

	testHelper.SetFields("f2")
	action.GenData([]string{consts.CommandTestFile})

	out := consts.Buf.String()

	t.Require().Contains(out, "123", "check generation")
}
