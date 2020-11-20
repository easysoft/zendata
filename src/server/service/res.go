package serverService

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/model"
	serverRepo "github.com/easysoft/zendata/src/server/repo"
	"github.com/easysoft/zendata/src/service"
	constant "github.com/easysoft/zendata/src/utils/const"
	"strings"
)

type ResService struct {
	rangesRepo *serverRepo.RangesRepo
	instancesRepo *serverRepo.InstancesRepo
	configRepo *serverRepo.ConfigRepo
	excelRepo *serverRepo.ExcelRepo
	textRepo *serverRepo.TextRepo
	defRepo *serverRepo.DefRepo
}

func (s *ResService) LoadRes(resType string) (ret []model.ResFile) {
	res, _, _ := service.LoadRes(resType)

	for _, key := range constant.ResKeys {
		for _, res := range res[key] {
			if res.ResType == constant.ResTypeExcel && strings.Index(res.Title, "|") > -1 {
				// more than 1 sheet
				arr := strings.Split(res.Title, "|")

				res.Title = arr[0]
				ret = append(ret,  res)

				var temp interface{} = res // clone
				res1 := temp.(model.ResFile)
				res1.Title = arr[1]
				ret = append(ret,  res1)
			} else {
				ret = append(ret,  res)
			}
		}
	}

	return
}

func (s *ResService) ListReferForSelection(resType string) (ret interface{}) {
	if resType == "ranges" {
		ret, _, _ = s.rangesRepo.List("",-1)
	} else if resType == "instances" {
		ret, _, _ = s.instancesRepo.List("",-1)
	} else if resType == "config" {
		ret, _, _ = s.configRepo.List("",-1)
	} else if resType == "yaml" {
		ret, _, _ = s.defRepo.List("",-1)
	} else if resType == "excel" {
		ret, _, _ = s.excelRepo.List("",-1)
	} else if resType == "text" {
		ret, _, _ = s.textRepo.List("",-1)
	}

	return
}
func (s *ResService) ListReferFieldForSelection(resId int, resType string) (ret []model.ResField) {
	if resType == "instances" {
		items, _ := s.instancesRepo.GetItems(resId)
		for i, item := range items {
			if item.ParentID != 0 { return } // ignore sub nodes
			field := model.ResField{Name: item.Instance, Index: i+1}
			ret = append(ret, field)
		}
	} else if resType == "ranges" {
		items, _ := s.rangesRepo.GetItems(resId)
		for i, item := range items {
			if item.ParentID != 0 { return } // ignore sub nodes
			field := model.ResField{Name: item.Name, Index: i+1}
			ret = append(ret, field)
		}
	} else if resType == "excel" {
		res, _ := s.excelRepo.Get(uint(resId))
		excel, _ := excelize.OpenFile(res.Path)

		for _, sheet := range excel.GetSheetList() {
			if res.Sheet != sheet { continue }

			rows, _ := excel.GetRows(sheet)

			for index, row := range rows {
				if index > 0 { break }
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
	}

	return
}

func NewResService(rangesRepo *serverRepo.RangesRepo,
	instancesRepo *serverRepo.InstancesRepo,
	configRepo *serverRepo.ConfigRepo,
	excelRepo *serverRepo.ExcelRepo,
	textRepo *serverRepo.TextRepo,
	defRepo *serverRepo.DefRepo) *ResService {
	return &ResService{rangesRepo: rangesRepo, instancesRepo: instancesRepo,
		configRepo: configRepo, excelRepo: excelRepo,
		textRepo: textRepo, defRepo: defRepo}
}
