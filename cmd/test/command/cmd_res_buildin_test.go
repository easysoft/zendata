package main

import (
	"strings"
	"testing"

	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestResBuildinCmd(t *testing.T) {
	suite.RunSuite(t, new(ResBuildinCmdSuite))
}

type ResBuildinCmdSuite struct {
	suite.Suite
}

func (s *ResBuildinCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("ResBuildinCmd")
}
func (s *ResBuildinCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("ResBuildinCmd")
}
func (s *ResBuildinCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *ResBuildinCmdSuite) TestResBuildin(t provider.T) {
	t.ID("0")

	pth := "test/unittest/out/res-buildin.txt"

	out := gen.New().
		SetConfigs([]string{"test/unittest/res-buildin.yaml"}).
		SetFields("").
		SetOutput(pth).
		SetTotal(1).
		Gen()

	length := len(strings.Split(strings.TrimSpace(out), "\n"))

	t.Require().Contains(out, "框框套套", "check generated data")
	t.Require().NotContains(out, "nil", "check not contains nil")
	t.Require().Equal(length, 35, "check generated data")
}
