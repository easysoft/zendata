package serverService

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/service"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"strings"
)

type ResService struct {
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

func (s *ResService) LoadResField(refer model.ResFile) (ret []model.ResField) {
	resFile, resType, _ := fileUtils.GetResProp(refer.Name)

	if resType == "yaml" {
		typ, inst, ranges := service.ReadYamlData(resFile)
		if typ == "inst" {
			for i, item := range inst.Instances {
				field := model.ResField{Name: item.Instance, Index: i+1}
				ret = append(ret, field)
			}
		} else if typ == "range" {
			i := 0
			for name, _ := range ranges.Ranges {
				field := model.ResField{Name: name, Index: i+1}
				ret = append(ret, field)
				i++
			}
		}
	} else if resType == "text" {
		// no need to show fields in webpage, used as list

	} else if resType == "excel" {
		excel, _ := excelize.OpenFile(resFile)
		for _, sheet := range excel.GetSheetList() {
			if refer.Title != sheet { continue }

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
