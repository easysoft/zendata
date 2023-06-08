package main

import (
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestResUsersCmd(t *testing.T) {
	suite.RunSuite(t, new(ResUsersCmdSuite))
}

type ResUsersCmdSuite struct {
	suite.Suite
}

func (s *ResUsersCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("ResUsersCmd")
}
func (s *ResUsersCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("ResUsersCmd")
}
func (s *ResUsersCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *ResUsersCmdSuite) TestResUsers(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"test/unittest/res-yaml.yaml"}).
		SetFields("f1").
		Gen()

	t.Require().Contains(out, "102", "check generated data")
}
