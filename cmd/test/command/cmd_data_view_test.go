package main

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestDataViewCmd(t *testing.T) {
	suite.RunSuite(t, new(DataViewCmdSuite))
}

type DataViewCmdSuite struct {
	suite.Suite
}

func (s *DataViewCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()

	t.AddSubSuite("DataViewCmd")
}
func (s *DataViewCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

//func (s *DataViewCmdSuite) TestDataList(t provider.T) {
//	t.ID("0")
//
//	helper.ListData()
//
//	out := consts.Buf.String()
//
//	t.Require().Contains(out, "语法说明", "check list content")
//}
//
//func (s *DataViewCmdSuite) TestResList(t provider.T) {
//	t.ID("0")
//
//	helper.ListRes()
//
//	out := consts.Buf.String()
//
//	t.Require().Contains(out, "语法说明", "check list content")
//}

func (s *DataViewCmdSuite) TestViewDetail(t provider.T) {
	t.ID("0")

	helper.View("city.v1.city")
	out := consts.Buf.String()

	helper.View("color/v1.yaml")
	out = consts.Buf.String()
	t.Require().Contains(out, "北京市", "check excel content")

	helper.View("color.v1.yaml")
	out = consts.Buf.String()

	t.Require().Contains(out, "北京市", "check excel content")
}
