package serverService

import (
	"github.com/easysoft/zendata/internal/pkg/domain"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverRepo "github.com/easysoft/zendata/internal/server/repo"
	serverUtils "github.com/easysoft/zendata/internal/server/utils"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gorm.io/gorm"
)

type TextService struct {
	TextRepo   *serverRepo.TextRepo `inject:""`
	ResService *ResService          `inject:""`
}

func (s *TextService) List(keywords string, page int) (list []*model.ZdText, total int) {
	list, total, _ = s.TextRepo.List(strings.TrimSpace(keywords), page)

	return
}

func (s *TextService) Get(id int) (text model.ZdText, dirs []domain.Dir) {
	text, _ = s.TextRepo.Get(uint(id))

	serverUtils.GetDirs(consts.ResDirYaml, &dirs)

	return
}

func (s *TextService) Save(text *model.ZdText) (err error) {
	text.Folder = serverUtils.DealWithPathSepRight(text.Folder)
	text.Path = vari.ZdDir + text.Folder + serverUtils.AddExt(text.FileName, ".txt")
	text.ReferName = helper.PathToName(text.Path, consts.ResDirYaml, consts.ResTypeText)

	if text.ID == 0 {
		err = s.Create(text)
	} else {
		err = s.Update(text)
	}

	return
}
func (s *TextService) Create(text *model.ZdText) (err error) {
	err = s.TextRepo.Create(text)

	return
}
func (s *TextService) Update(text *model.ZdText) (err error) {
	var old model.ZdText
	old, err = s.TextRepo.Get(text.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if text.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.TextRepo.Update(text)

	return
}

func (s *TextService) Remove(id int) (err error) {
	var old model.ZdText
	old, err = s.TextRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}

	fileUtils.RemoveExist(old.Path)
	err = s.TextRepo.Remove(uint(id))

	return
}

func (s *TextService) Sync(files []domain.ResFile) (err error) {
	list := s.TextRepo.ListAll()

	mp := map[string]*model.ZdText{}
	for _, item := range list {
		mp[item.Path] = item
	}

	for _, fi := range files {
		_, found := mp[fi.Path]
		//logUtils.PrintTo(fi.UpdatedAt.OpenApiDataTypeString() + ", " + mp[fi.Path].UpdatedAt.OpenApiDataTypeString())
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.TextRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		} else { // db is new

		}
	}

	return
}

func (s *TextService) SyncToDB(file domain.ResFile) (err error) {
	text := model.ZdText{
		Title:    file.Title,
		Path:     file.Path,
		Folder:   serverUtils.GetRelativePath(file.Path),
		FileName: fileUtils.GetFileName(file.Path),
	}
	if strings.Index(text.Path, consts.ResDirYaml) > -1 {
		text.ReferName = helper.PathToName(text.Path, consts.ResDirYaml, consts.ResTypeText)
	} else {
		text.ReferName = helper.PathToName(text.Path, consts.ResDirUsers, consts.ResTypeText)
	}
	text.Content = fileUtils.ReadFile(file.Path)

	s.TextRepo.Create(&text)

	return
}

func NewTextService(textRepo *serverRepo.TextRepo) *TextService {
	return &TextService{TextRepo: textRepo}
}
