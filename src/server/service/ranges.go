package serverService

import (
	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	"github.com/easysoft/zendata/src/service"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

type RangesService struct {
	rangesRepo  *serverRepo.RangesRepo
	resService  *ResService
	sectionRepo *serverRepo.SectionRepo
}

func (s *RangesService) List(keywords string, page int) (list []*model.ZdRanges, total int) {
	list, total, _ = s.rangesRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *RangesService) Get(id int) (ranges model.ZdRanges, dirs []model.Dir) {
	ranges, _ = s.rangesRepo.Get(uint(id))

	serverUtils.GetDirs(constant.ResDirYaml, &dirs)

	return
}

func (s *RangesService) Save(ranges *model.ZdRanges) (err error) {
	ranges.Folder = serverUtils.DealWithPathSepRight(ranges.Folder)
	ranges.Path = vari.ZdPath + ranges.Folder + serverUtils.AddExt(ranges.FileName, ".yaml")
	ranges.ReferName = service.PathToName(ranges.Path, constant.ResDirYaml, constant.ResTypeRanges)

	if ranges.ID == 0 {
		err = s.Create(ranges)
	} else {
		err = s.Update(ranges)
	}
	s.updateYaml(ranges.ID)

	return
}
func (s *RangesService) Create(ranges *model.ZdRanges) (err error) {
	err = s.rangesRepo.Create(ranges)

	return
}
func (s *RangesService) Update(ranges *model.ZdRanges) (err error) {
	var old model.ZdRanges
	old, err = s.rangesRepo.Get(ranges.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if ranges.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.rangesRepo.Update(ranges)

	return
}

func (s *RangesService) Remove(id int) (err error) {
	var old model.ZdRanges
	old, err = s.rangesRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}

	fileUtils.RemoveExist(old.Path)
	err = s.rangesRepo.Remove(uint(id))
	return
}

func (s *RangesService) Sync(files []model.ResFile) (err error) {
	list := s.rangesRepo.ListAll()

	mp := map[string]*model.ZdRanges{}
	for _, item := range list {
		mp[item.Path] = item
	}

	for _, fi := range files {
		_, found := mp[fi.Path]
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.rangesRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		}
	}

	return
}

func (s *RangesService) SyncToDB(fi model.ResFile) (err error) {
	content, _ := ioutil.ReadFile(fi.Path)
	yamlContent := stringUtils.ReplaceSpecialChars(content)
	po := model.ZdRanges{}
	err = yaml.Unmarshal(yamlContent, &po)

	po.Title = fi.Title
	po.Desc = fi.Desc
	po.Path = fi.Path
	po.Folder = serverUtils.GetRelativePath(po.Path)
	if strings.Index(po.Path, vari.ZdPath+constant.ResDirYaml) > -1 {
		po.ReferName = service.PathToName(po.Path, constant.ResDirYaml, constant.ResTypeRanges)
	} else {
		po.ReferName = service.PathToName(po.Path, constant.ResDirUsers, constant.ResTypeRanges)
	}
	po.FileName = fileUtils.GetFileName(po.Path)
	po.Yaml = string(content)

	s.rangesRepo.Create(&po)

	i := 1
	for k, v := range po.RangeMap {
		item := model.ZdRangesItem{Field: k, Value: v}
		item.RangesID = po.ID
		item.Ord = i
		s.rangesRepo.SaveItem(&item)
		i += 1

		rangeSections := gen.ParseRangeProperty(item.Value)
		for i, rangeSection := range rangeSections {
			s.sectionRepo.SaveFieldSectionToDB(rangeSection, i, item.ID, "ranges")
		}
	}

	return
}

func (s *RangesService) updateYamlByItem(itemId uint) (err error) {
	item, _ := s.rangesRepo.GetItem(itemId)
	return s.updateYaml(item.RangesID)
}
func (s *RangesService) updateYaml(id uint) (err error) {
	var po model.ZdRanges
	po, _ = s.rangesRepo.Get(id)

	s.genYaml(&po)
	err = s.rangesRepo.UpdateYaml(po)
	fileUtils.WriteFile(po.Path, po.Yaml)

	return
}
func (s *RangesService) genYaml(ranges *model.ZdRanges) (str string) {
	items, err := s.rangesRepo.GetItems(int(ranges.ID))
	if err != nil {
		return
	}

	yamlObj := model.ResRanges{}
	yamlObj.Ranges = map[string]string{}
	s.rangesRepo.GenRangesRes(*ranges, &yamlObj)

	for _, item := range items {

		yamlObj.Ranges[item.Field] = item.Value
	}

	bytes, err := yaml.Marshal(yamlObj)
	ranges.Yaml = stringUtils.ConvertYamlStringToMapFormat(bytes)

	return
}

func NewRangesService(rangesRepo *serverRepo.RangesRepo, sectionRepo *serverRepo.SectionRepo) *RangesService {
	return &RangesService{rangesRepo: rangesRepo, sectionRepo: sectionRepo}
}
