package main

import (
	"github.com/easysoft/zendata/cmd/test/gen"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestGenerateDemoCmd(t *testing.T) {
	suite.RunSuite(t, new(GenerateDemoCmdSuite))
}

type GenerateDemoCmdSuite struct {
	suite.Suite
}

func (s *GenerateDemoCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("GenerateDemoCmd")
}
func (s *GenerateDemoCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("GenerateDemoCmd")
}
func (s *GenerateDemoCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *GenerateDemoCmdSuite) TestGenerateDemoDefault(t provider.T) {
	t.ID("0")

	out := gen.New().
		SetConfigs([]string{"demo/default.yaml", "demo/config.yaml"}).
		SetFields("").
		Gen()

	t.Require().Contains(out, "123", "check generated data")
	t.Require().NotContains(out, "nil", "check not contains nil")
	t.Require().Contains(out, "{1}", "check overwrite pre-fix and post-fix")
}

//func (s *GenerateDemoCmdSuite) TestGenerateDemo01(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/01_range.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo02(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/02_fix.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo03(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/03_more_fields.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo04(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/04_rand.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo05(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/05_loop.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo06(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/06_from_file.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo07(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/07_nest.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo08(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/08_format.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo09(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/09_length.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo10(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/10_brace.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo11(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/11_loop_m_n.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo12(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/12_function.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo13(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/13_value.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo14(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/14_config.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo15(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/15_from_range.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo16(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/16_from_instance.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo17(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/17_from_results.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo18(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/18_from_excel.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo19(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/19_multi_from.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo20(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/20_nest_recursive.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo21(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/21_override.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo22(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/22_datetime.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo23(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/23_article.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo24(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/24_person_info.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo25(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/25_nested_json.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo26(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/26_assemble_json.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo27(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/config.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
//
//func (s *GenerateDemoCmdSuite) TestGenerateDemo28(t provider.T) {
//	t.ID("0")
//
//	out := gen.New().
//		SetConfigs([]string{"demo/default.yaml"}).
//		SetFields("").
//		Gen()
//
//	t.Require().Contains(out, "123", "check generated data")
//}
