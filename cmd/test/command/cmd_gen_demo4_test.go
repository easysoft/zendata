package main

import (
	"testing"

	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestGenerateDemo4Cmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateDemo4CmdSuite))

}

type GenerateDemo4CmdSuite struct {
	suite.Suite
}

func (s *GenerateDemo4CmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateDemo4Cmd")
}
func (s *GenerateDemo4CmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateDemo4Cmd")
}
func (s *GenerateDemo4CmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateDemo4CmdSuite) TestGenerateDemo22(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/22_datetime.yaml"}).
		SetFields("").
		Gen()

	regx := `201\d/`
	t.Require().Regexp(regx, out, "check a time that's 10 years before")
}
