package main

import (
	"testing"

	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestGenerateDemo1Cmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateDemo1CmdSuite))
}

type GenerateDemo1CmdSuite struct {
	suite.Suite
}

func (s *GenerateDemo1CmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateDemo1Cmd")
}
func (s *GenerateDemo1CmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateDemo1Cmd")
}
func (s *GenerateDemo1CmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemoDefault(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/default.yaml", "demo/config.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "123", "check generated data")
	t.Require().NotContains(out, "nil", "check not contains nil")
	t.Require().Contains(out, "{1}", "check overwrite pre-fix and post-fix")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo01(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/01_range.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "456\t1.70", "check generated data")
	t.Require().Contains(out, "[7]", "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo02(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/02_fix.yaml"}).
		SetFields("").
		Gen()

	expect := `int_2 | 2B___2B___2B___3B___`
	t.Require().Contains(out, expect, "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo04(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/04_rand.yaml"}).
		SetFields("").
		Gen()

	t.Require().NotContains(out, "1\t\n2\t\n3\t\n4\t\n5\t\n6\t\n7\t\n8\t\n9\t", "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo05(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/05_loop.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "c_d_e\t\n", "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo06(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/06_from_file.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "Tom", "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo07(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/07_nest.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "part1_B | part2_1", "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo08(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/08_format.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "passwd01\t", "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo09(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/09_length.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "006", "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo10(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/10_brace.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "a\t\na\t\na\t\nb\t\nb\t\nc\t\nc\t", "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo12(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/12_function.yaml"}).
		SetFields("").
		Gen()

	expect := "10 壹\thttp%3A%2F%2Fzendata.cn%3F%26%3Dword%2B\t0cc175b9c0f1b6a831c399e269772661\t5baa61e4c9b93f3f0682250b6cf8331b7ee68fd8\taHR0cDovL3plbmRhdGEuY24/Jj13b3JkKw=="
	t.Require().Contains(out, expect, "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo12_2(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/12_function2.yaml"}).
		SetFields("").
		Gen()

	t.Require().NotContains(out, "nil", "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo12_3(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/12_function3.yaml"}).
		SetFields("").
		Gen()

	expect := "==\tjson: {\"key\": \"value\"}"
	t.Require().Contains(out, expect, "check generated data")
}

func (s *GenerateDemo1CmdSuite) TestGenerateDemo13(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/13_value.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "1 x 6 = 6 three 1_three 6 x 6 = 36", "check generated data")
}
