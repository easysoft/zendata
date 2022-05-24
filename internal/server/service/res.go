package serverService

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverRepo "github.com/easysoft/zendata/internal/server/repo"
	"strings"
)

type ResService struct {
	RangesRepo    *serverRepo.RangesRepo    `inject:""`
	InstancesRepo *serverRepo.InstancesRepo `inject:""`
	ConfigRepo    *serverRepo.ConfigRepo    `inject:""`
	ExcelRepo     *serverRepo.ExcelRepo     `inject:""`
	TextRepo      *serverRepo.TextRepo      `inject:""`
	DefRepo       *serverRepo.DefRepo       `inject:""`
}

func (s *ResService) LoadRes(resType string) (ret []model.ResFile) {
	res, _, _ := helper.LoadRes(resType)

	for _, key := range constant.ResKeys {
		for _, res := range res[key] {
			if res.ResType == constant.ResTypeExcel && strings.Index(res.Title, "|") > -1 {
				// more than 1 sheet
				arr := strings.Split(res.Title, "|")

				res.Title = arr[0]
				ret = append(ret, res)

				var temp interface{} = res // clone
				res1 := temp.(model.ResFile)
				res1.Title = arr[1]
				ret = append(ret, res1)
			} else {
				ret = append(ret, res)
			}
		}
	}

	return
}

func (s *ResService) ListReferFileForSelection(resType string) (ret interface{}) {
	if resType == "ranges" {
		ret = s.RangesRepo.ListAll()
	} else if resType == "instances" {
		ret = s.InstancesRepo.ListAll()
	} else if resType == "config" {
		ret = s.ConfigRepo.ListAll()
	} else if resType == "yaml" {
		ret = s.DefRepo.ListAll()
	} else if resType == "excel" {
		ret = s.ExcelRepo.ListFiles()
	} else if resType == "text" {
		ret = s.TextRepo.ListAll()
	}

	return
}

func (s *ResService) ListReferSheetForSelection(referName string) (ret []*model.ZdExcel) {
	index := strings.LastIndex(referName, ".")
	file := referName[:index]

	ret = s.ExcelRepo.ListSheets(file)

	return
}
func (s *ResService) ListReferExcelColForSelection(referName string) (ret []model.ResField) {
	index := strings.LastIndex(referName, ".")
	file := referName[:index]
	sheet := referName[index+1:]

	res, _ := s.ExcelRepo.GetBySheet(file, sheet)
	excel, _ := excelize.OpenFile(res.Path)

	for _, sheet := range excel.GetSheetList() {
		if res.Sheet != sheet {
			continue
		}

		rows, _ := excel.GetRows(sheet)

		for index, row := range rows {
			if index > 0 {
				break
			}
			for i, col := range row {
				val := strings.TrimSpace(col)
				if val == "" {
					break
				}

				field := model.ResField{Name: val, Index: i + 1}
				ret = append(ret, field)
			}
		}
	}

	return
}

func (s *ResService) ListReferResFieldForSelection(resId int, resType string) (ret []model.ResField) {
	if resType == "instances" {
		items, _ := s.InstancesRepo.GetItems(uint(resId))
		for i, item := range items {
			if item.ParentID != 0 {
				return
			} // ignore sub nodes
			field := model.ResField{Name: item.Instance, Index: i + 1}
			ret = append(ret, field)
		}
	} else if resType == "ranges" {
		items, _ := s.RangesRepo.GetItems(resId)
		for i, item := range items {
			if item.ParentID != 0 {
				return
			} // ignore sub nodes
			field := model.ResField{Name: item.Field, Index: i + 1}
			ret = append(ret, field)
		}
	}

	return
}

func NewResService(rangesRepo *serverRepo.RangesRepo,
	instancesRepo *serverRepo.InstancesRepo,
	configRepo *serverRepo.ConfigRepo,
	excelRepo *serverRepo.ExcelRepo,
	textRepo *serverRepo.TextRepo,
	defRepo *serverRepo.DefRepo) *ResService {
	return &ResService{RangesRepo: rangesRepo, InstancesRepo: instancesRepo,
		ConfigRepo: configRepo, ExcelRepo: excelRepo,
		TextRepo: textRepo, DefRepo: defRepo}
}
