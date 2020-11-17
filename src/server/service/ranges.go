package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	logUtils "github.com/easysoft/zendata/src/utils/log"
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

func (s *RangesService) GetItemTree(rangesId int) (root model.ZdRangesItem) {
	items, _ := s.rangesRepo.GetItems(rangesId)

	root.ID = 0
	root.Name = "序列"
	for _, item := range items {
		item.ParentID = root.ID
		root.Children = append(root.Children, item)
	}

	return
}
func (s *RangesService) GetItem(id int) (item model.ZdRangesItem) {
	item, _ = s.rangesRepo.GetItem(uint(id))
	return
}

func (s *RangesService) Save(ranges *model.ZdRanges) (err error) {
	err = s.rangesRepo.Save(ranges)

	return
}

func (s *RangesService) Remove(id int) (err error) {
	err = s.rangesRepo.Remove(uint(id))
	if err != nil {
		return
	}

	ranges, _ := s.rangesRepo.Get(uint(id))
	logUtils.PrintTo(ranges.Path)
	//fileUtils.RemoveExist(ranges.Path)

	return
}

func (s *RangesService) CreateItem(domainId, targetId int, mode string) (item *model.ZdRangesItem, err error) {
	item = &model.ZdRangesItem{Name: "ranges_", RangesID: uint(domainId)}
	item.Ord = s.rangesRepo.GetMaxOrder(domainId)

	err = s.rangesRepo.SaveItem(item)

	return
}
func (s *RangesService) SaveItem(item *model.ZdRangesItem) (err error) {
	err = s.rangesRepo.SaveItem(item)
	return
}

func (s *RangesService) RemoveItem(id int) (err error) {
	err = s.rangesRepo.RemoveItem(uint(id))
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
