package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/jinzhu/gorm"
)

type TextService struct {
	textRepo *serverRepo.TextRepo
	resService *ResService
}

func (s *TextService) List() (list []*model.ZdText) {
	texts := s.resService.LoadRes("text")
	list, _ = s.textRepo.List()

	s.saveResToDB(texts, list)
	list, _ = s.textRepo.List()

	return
}

func (s *TextService) Get(id int) (text model.ZdText) {
	text, _ = s.textRepo.Get(uint(id))

	return
}

func (s *TextService) Save(text *model.ZdText) (err error) {
	err = s.textRepo.Save(text)

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

func (s *TextService) saveResToDB(texts []model.ResFile, list []*model.ZdText) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range texts {
		if !stringUtils.FindInArrBool(item.Path, names) {
			text := model.ZdText{Title: item.Title, Name: item.Name,Path: item.Path}

			content := fileUtils.ReadFile(item.Path)
			text.Content = content

			s.textRepo.Save(&text)
		}
	}

	return
}

func NewTextService(textRepo *serverRepo.TextRepo) *TextService {
	return &TextService{textRepo: textRepo}
}
