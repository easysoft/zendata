package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	"github.com/easysoft/zendata/src/service"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/gorm"
)

type ExcelService struct {
	excelRepo *serverRepo.ExcelRepo
	resService *ResService
}

func (s *ExcelService) List() (list []*model.ZdExcel) {
	excel := s.resService.LoadRes("excel")
	list, _ = s.excelRepo.List()

	s.importResToDB(excel, list)
	list, _ = s.excelRepo.List()

	return
}

func (s *ExcelService) Get(id int) (excel model.ZdExcel, dirTree model.Dir) {
	excel, _ = s.excelRepo.Get(uint(id))

	dirTree = model.Dir{Name: fileUtils.AddSepIfNeeded(constant.ResDirData)}
	serverUtils.GetDirTree(&dirTree)

	return
}

func (s *ExcelService) Save(excel *model.ZdExcel) (err error) {
	excel.Folder = serverUtils.DealWithPathSepRight(excel.Folder)
	excel.Path = vari.WorkDir + excel.Folder + serverUtils.AddExt(excel.Title, ".xlsx")
	excel.Name = service.PathToName(excel.Path, constant.ResDirData)

	if excel.ID == 0 {
		// excel should not be create on webpage
	} else {
		err = s.Update(excel)
	}

	return
}
func (s *ExcelService) Update(excel *model.ZdExcel) (err error) {
	var old model.ZdExcel
	old, err = s.excelRepo.Get(excel.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if excel.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.excelRepo.Update(excel)

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

func (s *ExcelService) dataToYaml(excel *model.ZdExcel) (str string) {

	return
}

func (s *ExcelService) importResToDB(excel []model.ResFile, list []*model.ZdExcel) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range excel {
		if !stringUtils.FindInArrBool(item.Path, names) {
			excel := model.ZdExcel{Title: item.Title, Name: item.Name,
				Path: item.Path,
				Sheet: item.Title}
			excel.Folder = serverUtils.GetRelativePath(excel.Path)
			s.excelRepo.Create(&excel)
		}
	}

	return
}

func NewExcelService(excelRepo *serverRepo.ExcelRepo) *ExcelService {
	return &ExcelService{excelRepo: excelRepo}
}
