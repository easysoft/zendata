package serverService

import (
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/service"
	"os"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverRepo "github.com/easysoft/zendata/internal/server/repo"
	serverUtils "github.com/easysoft/zendata/internal/server/utils"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type RangesService struct {
	RangesRepo  *serverRepo.RangesRepo  `inject:""`
	ResService  *ResService             `inject:""`
	SectionRepo *serverRepo.SectionRepo `inject:""`

	SectionService *SectionService       `inject:""`
	RangeService   *service.RangeService `inject:""`
}

func (s *RangesService) List(keywords string, page int) (list []*model.ZdRanges, total int) {
	list, total, _ = s.RangesRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *RangesService) Get(id int) (ranges model.ZdRanges, dirs []domain.Dir) {
	ranges, _ = s.RangesRepo.Get(uint(id))

	serverUtils.GetDirs(consts.ResDirYaml, &dirs)

	return
}

func (s *RangesService) Save(ranges *model.ZdRanges) (err error) {
	ranges.Folder = serverUtils.DealWithPathSepRight(ranges.Folder)
	ranges.Path = vari.WorkDir + ranges.Folder + serverUtils.AddExt(ranges.FileName, ".yaml")
	ranges.ReferName = helper.PathToName(ranges.Path, consts.ResDirYaml, consts.ResTypeRanges)

	if ranges.ID == 0 {
		err = s.Create(ranges)
	} else {
		err = s.Update(ranges)
	}
	s.updateYaml(ranges.ID)

	return
}
func (s *RangesService) Create(ranges *model.ZdRanges) (err error) {
	err = s.RangesRepo.Create(ranges)

	return
}
func (s *RangesService) Update(ranges *model.ZdRanges) (err error) {
	var old model.ZdRanges
	old, err = s.RangesRepo.Get(ranges.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if ranges.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.RangesRepo.Update(ranges)

	return
}

func (s *RangesService) Remove(id int) (err error) {
	var old model.ZdRanges
	old, err = s.RangesRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}

	fileUtils.RemoveExist(old.Path)
	err = s.RangesRepo.Remove(uint(id))
	return
}

func (s *RangesService) Sync(files []domain.ResFile) (err error) {
	list := s.RangesRepo.ListAll()

	mp := map[string]*model.ZdRanges{}
	for _, item := range list {
		mp[item.Path] = item
	}

	for _, fi := range files {
		_, found := mp[fi.Path]
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.RangesRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		}
	}

	return
}

func (s *RangesService) SyncToDB(fi domain.ResFile) (err error) {
	content, _ := os.ReadFile(fi.Path)
	yamlContent := helper.ReplaceSpecialChars(content)
	po := model.ZdRanges{}
	err = yaml.Unmarshal(yamlContent, &po)

	po.Title = fi.Title
	po.Desc = fi.Desc
	po.Path = fi.Path
	po.Folder = serverUtils.GetRelativePath(po.Path)
	if strings.Index(po.Path, vari.WorkDir+consts.ResDirYaml) > -1 {
		po.ReferName = helper.PathToName(po.Path, consts.ResDirYaml, consts.ResTypeRanges)
	} else {
		po.ReferName = helper.PathToName(po.Path, consts.ResDirUsers, consts.ResTypeRanges)
	}
	po.FileName = fileUtils.GetFileName(po.Path)
	po.Yaml = string(content)

	s.RangesRepo.Create(&po)

	i := 1
	for k, v := range po.RangeMap {
		item := model.ZdRangesItem{Field: k, Value: v}
		item.RangesID = po.ID
		item.Ord = i
		s.RangesRepo.SaveItem(&item)
		i += 1

		rangeSections := s.RangeService.ParseRangeProperty(item.Value)
		for i, rangeSection := range rangeSections {
			s.SectionService.SaveFieldSectionToDB(rangeSection, i, item.ID, "ranges")
		}
	}

	return
}

func (s *RangesService) updateYamlByItem(itemId uint) (err error) {
	item, _ := s.RangesRepo.GetItem(itemId)
	return s.updateYaml(item.RangesID)
}
func (s *RangesService) updateYaml(id uint) (err error) {
	var po model.ZdRanges
	po, _ = s.RangesRepo.Get(id)

	s.genYaml(&po)
	err = s.RangesRepo.UpdateYaml(po)
	fileUtils.WriteFile(po.Path, po.Yaml)

	return
}
func (s *RangesService) genYaml(ranges *model.ZdRanges) (str string) {
	items, err := s.RangesRepo.GetItems(int(ranges.ID))
	if err != nil {
		return
	}

	yamlObj := domain.ResRanges{}
	yamlObj.Ranges = map[string]string{}
	s.RangesRepo.GenRangesRes(*ranges, &yamlObj)

	for _, item := range items {

		yamlObj.Ranges[item.Field] = item.Value
	}

	bytes, err := yaml.Marshal(yamlObj)
	ranges.Yaml = helper.ConvertYamlStringToMapFormat(bytes)

	return
}

func NewRangesService(rangesRepo *serverRepo.RangesRepo, sectionRepo *serverRepo.SectionRepo) *RangesService {
	return &RangesService{RangesRepo: rangesRepo, SectionRepo: sectionRepo}
}
