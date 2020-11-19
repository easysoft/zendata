package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/jinzhu/gorm"
)

type ExcelService struct {
	excelRepo *serverRepo.ExcelRepo
	resService *ResService
}

func (s *ExcelService) List() (list []*model.ZdExcel) {
	excel := s.resService.LoadRes("excel")
	list, _ = s.excelRepo.List()

	s.saveResToDB(excel, list)
	list, _ = s.excelRepo.List()

	return
}

func (s *ExcelService) Get(id int) (excel model.ZdExcel) {
	excel, _ = s.excelRepo.Get(uint(id))

	return
}

func (s *ExcelService) Save(excel *model.ZdExcel) (err error) {
	err = s.excelRepo.Save(excel)

	return
}

func (s *ExcelService) Remove(id int) (err error) {
	var old model.ZdExcel
	old, err = s.excelRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}
	fileUtils.RemoveExist(old.Path)

	err = s.excelRepo.Remove(uint(id))
	return
}

func (s *ExcelService) saveResToDB(excel []model.ResFile, list []*model.ZdExcel) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range excel {
		if !stringUtils.FindInArrBool(item.Path, names) {
			excel := model.ZdExcel{Title: item.Title, Name: item.Name, Path: item.Path, Sheet: item.Title}
			s.excelRepo.Save(&excel)
		}
	}

	return
}

func NewExcelService(excelRepo *serverRepo.ExcelRepo) *ExcelService {
	return &ExcelService{excelRepo: excelRepo}
}
