package serverService

import (
	"github.com/easysoft/zendata/src/model"
)

func (s *RangesService) GetItemTree(rangesId int) (root model.ZdRangesItem) {
	items, _ := s.RangesRepo.GetItems(rangesId)

	root.ID = 0
	root.Field = "序列"
	for _, item := range items {
		item.ParentID = root.ID
		root.Fields = append(root.Fields, item)
	}

	return
}
func (s *RangesService) GetItem(id int) (item model.ZdRangesItem) {
	item, _ = s.RangesRepo.GetItem(uint(id))
	return
}

func (s *RangesService) CreateItem(domainId, targetId int, mode string) (item *model.ZdRangesItem, err error) {
	item = &model.ZdRangesItem{Field: "ranges_", RangesID: uint(domainId)}
	item.Ord = s.RangesRepo.GetMaxOrder(domainId)

	err = s.RangesRepo.SaveItem(item)

	return
}
func (s *RangesService) SaveItem(item *model.ZdRangesItem) (err error) {
	err = s.RangesRepo.SaveItem(item)
	s.updateYaml(item.RangesID)
	return
}

func (s *RangesService) RemoveItem(id, domainId int) (err error) {
	err = s.RangesRepo.RemoveItem(uint(id))
	s.updateYaml(uint(domainId))
	return
}
