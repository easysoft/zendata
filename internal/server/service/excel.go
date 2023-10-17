package serverService

import (
	"strings"

	"github.com/easysoft/zendata/internal/pkg/domain"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverRepo "github.com/easysoft/zendata/internal/server/repo"
	serverUtils "github.com/easysoft/zendata/internal/server/utils"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/gorm"
)

type ExcelService struct {
	ExcelRepo  *serverRepo.ExcelRepo `inject:""`
	ResService *ResService           `inject:""`
}

func (s *ExcelService) List(keywords string, page int) (list []*model.ZdExcel, total int) {
	list, total, _ = s.ExcelRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *ExcelService) Get(id int) (excel model.ZdExcel, dirs []domain.Dir) {
	excel, _ = s.ExcelRepo.Get(uint(id))

	serverUtils.GetDirs(consts.ResDirData, &dirs)

	return
}

func (s *ExcelService) Save(excel *model.ZdExcel) (err error) {
	excel.Folder = serverUtils.DealWithPathSepRight(excel.Folder)
	excel.Path = vari.WorkDir + excel.Folder + serverUtils.AddExt(excel.FileName, ".xlsx")
	excel.ReferName = helper.PathToName(excel.Path, consts.ResDirData, consts.ResTypeExcel)

	if excel.ID == 0 {
		// excel should not be create on webpage
	} else {
		err = s.Update(excel)
	}

	return
}
func (s *ExcelService) Update(excel *model.ZdExcel) (err error) {
	var old model.ZdExcel
	old, err = s.ExcelRepo.Get(excel.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if excel.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.ExcelRepo.Update(excel)

	return
}

func (s *ExcelService) Remove(id int) (err error) {
	var old model.ZdExcel
	old, err = s.ExcelRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}

	fileUtils.RemoveExist(old.Path)
	err = s.ExcelRepo.Remove(uint(id))
	return
}

func (s *ExcelService) Sync(files []domain.ResFile) (err error) {
	list := s.ExcelRepo.ListAll()

	mp := map[string]*model.ZdExcel{}
	for _, item := range list {
		mp[item.Path] = item
	}

	for _, fi := range files {
		_, found := mp[fi.Path]
		//logUtils.PrintTo(fi.UpdatedAt.OpenApiDataTypeString() + ", " + mp[fi.Path].UpdatedAt.OpenApiDataTypeString())
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.ExcelRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		} else { // db is new

		}
	}

	return
}
func (s *ExcelService) SyncToDB(file domain.ResFile) (err error) {
	excel := model.ZdExcel{
		Title:     file.Title,
		Sheet:     file.Title,
		Path:      file.Path,
		Folder:    serverUtils.GetRelativePath(file.Path),
		ReferName: helper.PathToName(file.Path, consts.ResDirData, consts.ResTypeExcel),
		FileName:  fileUtils.GetFileName(file.Path),
	}
	s.ExcelRepo.Create(&excel)

	return
}

func NewExcelService(excelRepo *serverRepo.ExcelRepo) *ExcelService {
	return &ExcelService{ExcelRepo: excelRepo}
}
