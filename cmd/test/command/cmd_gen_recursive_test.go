package main

import (
	"github.com/easysoft/zendata/cmd/command/action"
	"github.com/easysoft/zendata/cmd/test/consts"
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

func (s *GenerateRecursiveCmdSuite) TestGenerateRecursive(t provider.T) {
	t.ID("0")

	action.GenData([]string{"demo/17_nest_recursive.yaml"})

	out := consts.Buf.String()

	t.Require().Contains(out, "3_C\t1_C\t", "check generation")
}
