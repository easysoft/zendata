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

func (s *DataViewCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("DataListCmd")
}
func (s *DataViewCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()
	t.AddSubSuite("DataViewCmd")
}
func (s *DataViewCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *DataViewCmdSuite) TestViewDataBuildinExcel(t provider.T) {
	t.ID("0")

	consts.Buf.Reset()
	helper.View("city.v1.city")
	out := consts.Buf.String()
	t.Require().Contains(out, "北京市", "check excel content")
}

func (s *DataViewCmdSuite) TestViewDataBuildinConfig(t provider.T) {
	t.ID("0")

	consts.Buf.Reset()
	helper.View("color/v1.yaml")
	out := consts.Buf.String()
	t.Require().Contains(out, "rgb", "check excel content")

	consts.Buf.Reset()
	helper.View("color.v1.yaml")
	out = consts.Buf.String()

	t.Require().Contains(out, "rgb", "check excel content")
}
