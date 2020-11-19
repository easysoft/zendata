package serverService

import (
	"github.com/easysoft/zendata/src/model"
)

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
