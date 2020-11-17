package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
)

type RangesService struct {
	rangesRepo *serverRepo.RangesRepo
	resService *ResService
}

func (s *RangesService) List() (list []*model.ZdRanges) {
	ranges := s.resService.LoadRes("ranges")
	list, _ = s.rangesRepo.List()

	s.saveResToDB(ranges, list)
	list, _ = s.rangesRepo.List()

	return
}

func (s *RangesService) Get(id int) (ranges model.ZdRanges) {
	ranges, _ = s.rangesRepo.Get(uint(id))

	return
}

func (s *RangesService) Save(ranges *model.ZdRanges) (err error) {
	err = s.rangesRepo.Save(ranges)

	return
}

func (s *RangesService) Remove(id int) (err error) {
	err = s.rangesRepo.Remove(uint(id))
	return
}

func (s *RangesService) saveResToDB(ranges []model.ResFile, list []*model.ZdRanges) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range ranges {
		if !stringUtils.FindInArrBool(item.Path, names) {
			ranges := model.ZdRanges{Path: item.Path, Name: item.Name, Title: item.Title, Desc: item.Desc, Field: item.Title, Note: item.Desc}
			s.rangesRepo.Save(&ranges)
		}
	}

	return
}

func NewRangesService(rangesRepo *serverRepo.RangesRepo) *RangesService {
	return &RangesService{rangesRepo: rangesRepo}
}
