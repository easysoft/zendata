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

type TextService struct {
	textRepo *serverRepo.TextRepo
	resService *ResService
}

func (s *TextService) List() (list []*model.ZdText) {
	texts := s.resService.LoadRes("text")
	list, _ = s.textRepo.List()

	s.importResToDB(texts, list)
	list, _ = s.textRepo.List()

	return
}

func (s *TextService) Get(id int) (text model.ZdText) {
	text, _ = s.textRepo.Get(uint(id))

	return
}

func (s *TextService) Save(text *model.ZdText) (err error) {
	text.Folder = serverUtils.DealWithPathSepRight(text.Folder)
	text.Path = vari.WorkDir + text.Folder + serverUtils.AddExt(text.Title, ".txt")
	text.Name = service.PathToName(text.Path, constant.ResDirYaml)

	if text.ID == 0 {
		err = s.Create(text)
	} else {
		err = s.Update(text)
	}

	return
}
func (s *TextService) Create(text *model.ZdText) (err error) {
	err = s.textRepo.Create(text)

	return
}
func (s *TextService) Update(text *model.ZdText) (err error) {
	var old model.ZdText
	old, err = s.textRepo.Get(text.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if text.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.textRepo.Update(text)

	return
}

func (s *TextService) Remove(id int) (err error) {
	var old model.ZdText
	old, err = s.textRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}
	fileUtils.RemoveExist(old.Path)

	err = s.textRepo.Remove(uint(id))

	return
}

func (s *TextService) dataToYaml(text *model.ZdText) (str string) {

	return
}

func (s *TextService) importResToDB(texts []model.ResFile, list []*model.ZdText) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range texts {
		if !stringUtils.FindInArrBool(item.Path, names) {
			text := model.ZdText{Title: item.Title, Name: item.Name,Path: item.Path}
			text.Folder = serverUtils.GetRelativePath(text.Path)
			content := fileUtils.ReadFile(item.Path)
			text.Content = content

			s.textRepo.Create(&text)
		}
	}

	return
}

func NewTextService(textRepo *serverRepo.TextRepo) *TextService {
	return &TextService{textRepo: textRepo}
}
