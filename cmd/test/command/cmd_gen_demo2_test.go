package main

import (
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestGenerateDemo2Cmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateDemo2CmdSuite))
}

type GenerateDemo2CmdSuite struct {
	suite.Suite
}

func (s *GenerateDemo2CmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateDemo2Cmd")
}
func (s *GenerateDemo2CmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateDemo2Cmd")
}
func (s *GenerateDemo2CmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo14(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/14_from_config.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "'106,107'", "check generated data")
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo15(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/15_from_range.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "103\t101\t", "check generated data")
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo16(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/16_from_instance.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "172.18.2.3\t192.168.0.1", "check generated data")
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo17(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/17_from_results.yaml"}).
		SetFields("").
		Gen()

	expect := "'106,107'\t ~~~ 101\t1\t"
	t.Require().Contains(out, expect, "check generated data")
}

//func (s *GenerateDemo2CmdSuite) TestGenerateDemo18(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/18_from_excel.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "", "check generated data")
//}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo19(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/19_multi_from.yaml"}).
		SetFields("").
		Gen()

	expect := "[172.18.2.3]\t\n[1.0.0.1]"
	t.Require().Contains(out, expect, "check generated data")
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo20(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/20_children_recursive.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "2_B\t1_B\t", "check generated data")
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo22(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/22_datetime.yaml"}).
		SetFields("").
		Gen()

	regx := `201\d/`
	t.Require().Regexp(regx, out, "check a time that's 10 years before")
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo23(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/23_article.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "家人", "check generated data")
	t.Require().NotContains(out, "{人称代词}", "check generated data")
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo24(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/24_person_info.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "139", "check generated data")
	t.Require().NotContains(out, "nil", "check generated data")
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo25(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/25_json_nested.yaml"}).
		SetFields("").
		SetOutput("test/out/result.json").
		Gen()

	expect := `"join_true": "PART1_A | part2_a | CHILD1_A | CHILD2_a"`
	t.Require().Contains(out, expect, "check generated data")
}

func (s *GenerateDemo2CmdSuite) TestGenerateDemo26(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/26_json_join_children.yaml"}).
		SetFields("").
		SetOutput("test/out/result.json").
		Gen()

	expect := `"field2": "PART1_A | part2_U | part3_son1_b_son2_v"`
	t.Require().Contains(out, expect, "check generated data")
}
