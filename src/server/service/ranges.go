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
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

type RangesService struct {
	rangesRepo *serverRepo.RangesRepo
	resService *ResService
}

func (s *RangesService) List(keywords string, page int) (list []*model.ZdRanges, total int) {
	ranges := s.resService.LoadRes("ranges")
	list, _, _ = s.rangesRepo.List("", -1)

	s.importResToDB(ranges, list)
	list, total, _ = s.rangesRepo.List(strings.TrimSpace(keywords), page)

	return
}

func (s *RangesService) Get(id int) (ranges model.ZdRanges, dirTree model.Dir) {
	ranges, _ = s.rangesRepo.Get(uint(id))

	dirTree = model.Dir{Name: fileUtils.AddSepIfNeeded(constant.ResDirYaml)}
	serverUtils.GetDirTree(&dirTree)

	return
}

func (s *RangesService) Save(ranges *model.ZdRanges) (err error) {
	ranges.Folder = serverUtils.DealWithPathSepRight(ranges.Folder)
	ranges.Path = vari.WorkDir + ranges.Folder + serverUtils.AddExt(ranges.Title, ".yaml")
	ranges.Name = service.PathToName(ranges.Path, constant.ResDirYaml)

	if ranges.ID == 0 {
		err = s.Create(ranges)
	} else {
		err = s.Update(ranges)
	}

	return
}
func (s *RangesService) Create(ranges *model.ZdRanges) (err error) {
	s.dataToYaml(ranges)
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

	s.dataToYaml(ranges)
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

func (s *RangesService) dataToYaml(ranges *model.ZdRanges) (str string) {

	return
}

func (s *RangesService) importResToDB(ranges []model.ResFile, list []*model.ZdRanges) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range ranges {
		if !stringUtils.FindInArrBool(item.Path, names) {
			content, _ := ioutil.ReadFile(item.Path)
			yamlContent := stringUtils.ReplaceSpecialChars(content)
			ranges := model.ZdRanges{}
			err = yaml.Unmarshal(yamlContent, &ranges)
			ranges.Title = item.Title
			ranges.Name = item.Name
			ranges.Desc = item.Desc
			ranges.Path = item.Path
			ranges.Folder = serverUtils.GetRelativePath(ranges.Path)
			ranges.Field = item.Title
			ranges.Note = item.Desc
			ranges.Yaml = string(content)

			s.rangesRepo.Create(&ranges)

			i := 1
			for k, v := range ranges.RangeMap {
				item := model.ZdRangesItem{Name: k, Value: v}
				item.RangesID = ranges.ID
				item.Ord = i
				s.rangesRepo.SaveItem(&item)
				i += 1
			}
		}
	}

	return
}

func NewRangesService(rangesRepo *serverRepo.RangesRepo) *RangesService {
	return &RangesService{rangesRepo: rangesRepo}
}
