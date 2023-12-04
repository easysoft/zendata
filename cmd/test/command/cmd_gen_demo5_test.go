package main

import (
	"testing"

	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestGenerateDemo5Cmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateDemo5CmdSuite))

}

type GenerateDemo5CmdSuite struct {
	suite.Suite
}

func (s *GenerateDemo5CmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateDemo2Cmd")
}
func (s *GenerateDemo5CmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateDemo2Cmd")
}
func (s *GenerateDemo5CmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateDemo5CmdSuite) TestGenerateDemo24(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/24_person_info.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "", "check generated data")
	t.Require().NotContains(out, "nil", "check generated data")
}
