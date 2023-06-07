package main

import (
	"github.com/easysoft/zendata/cmd/test/consts"
	testHelper "github.com/easysoft/zendata/cmd/test/helper"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"testing"
)

func TestDataListCmd(t *testing.T) {
	suite.RunSuite(t, new(DataListCmdSuite))
}

type DataListCmdSuite struct {
	suite.Suite
}

func (s *DataListCmdSuite) BeforeAll(t provider.T) {
	testHelper.BeforeAll()
	t.AddSubSuite("DataListCmd")
}
func (s *DataListCmdSuite) BeforeEach(t provider.T) {
	testHelper.PreCase()

	t.AddSubSuite("DataListCmd")
}
func (s *DataListCmdSuite) AfterEach(t provider.T) {
	testHelper.PostCase()
}

func (s *DataListCmdSuite) TestDataList(t provider.T) {
	t.ID("0")

	helper.ListData()

	out := consts.Buf.String()

	t.Require().Contains(out, "test.v1.yaml", "check list content")
}

func (s *DataListCmdSuite) TestResList(t provider.T) {
	t.ID("0")

	helper.ListRes()

	out := consts.Buf.String()

	t.Require().Contains(out, "city.v1.xlsx", "check list content")
}
