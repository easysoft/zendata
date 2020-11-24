package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	"github.com/easysoft/zendata/src/service"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/gorm"
	"strings"
)

type TextService struct {
	textRepo *serverRepo.TextRepo
	resService *ResService
}

func (s *TextService) List(keywords string, page int) (list []*model.ZdText, total int) {
	list, total, _ = s.textRepo.List(strings.TrimSpace(keywords), page)

	return
}

func (s *TextService) Get(id int) (text model.ZdText, dirs []model.Dir) {
	text, _ = s.textRepo.Get(uint(id))

	serverUtils.GetDirs(constant.ResDirYaml, &dirs)

	return
}

func (s *TextService) Save(text *model.ZdText) (err error) {
	text.Folder = serverUtils.DealWithPathSepRight(text.Folder)
	text.Path = vari.WorkDir + text.Folder + serverUtils.AddExt(text.FileName, ".txt")
	text.ReferName = service.PathToName(text.Path, constant.ResDirYaml, constant.ResTypeText)

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

func (s *TextService) Sync(files []model.ResFile) (err error) {
	list := s.textRepo.ListAll()

	mp := map[string]*model.ZdText{}
	for _, item := range list {
		mp[item.Path] = item
	}

	for _, fi := range files {
		_, found := mp[fi.Path]
		//logUtils.PrintTo(fi.UpdatedAt.String() + ", " + mp[fi.Path].UpdatedAt.String())
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.textRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		} else { // db is new

		}
	}

	return
}

func (s *TextService) SyncToDB(file model.ResFile) (err error) {
	text := model.ZdText{
		Title: file.Title,
		Path: file.Path,
		Folder: serverUtils.GetRelativePath(file.Path),
		FileName: fileUtils.GetFileName(file.Path),
	}
	if strings.Index(text.Path, constant.ResDirYaml) > -1 {
		text.ReferName = service.PathToName(text.Path, constant.ResDirYaml, constant.ResTypeText)
	} else {
		text.ReferName = service.PathToName(text.Path, constant.ResDirUsers, constant.ResTypeText)
	}
	text.Content = fileUtils.ReadFile(file.Path)

	s.textRepo.Create(&text)

	return
}

func NewTextService(textRepo *serverRepo.TextRepo) *TextService {
	return &TextService{textRepo: textRepo}
}
