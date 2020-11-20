package serverService

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
)

func (s *InstancesService) GetItemTree(instancesId int) (root model.ZdInstancesItem) {
	root = s.instancesRepo.GetItemTree(uint(instancesId))
	return
}
func (s *InstancesService) GetItem(id int) (item model.ZdInstancesItem) {
	item, _ = s.instancesRepo.GetItem(uint(id))
	return
}

func (s *InstancesService) CreateItem(domainId, targetId int, mode string) (item *model.ZdInstancesItem, err error) {
	item = &model.ZdInstancesItem{Field: "instances_", Note: "", InstancesID: uint(domainId)}
	item.Ord = s.instancesRepo.GetMaxOrder(domainId)

	err = s.instancesRepo.SaveItem(item)
	s.referRepo.CreateDefault(item.ID, constant.ResTypeInstances)

	return
}
func (s *InstancesService) SaveItem(item *model.ZdInstancesItem) (err error) {
	err = s.instancesRepo.SaveItem(item)
	return
}

func (s *InstancesService) RemoveItem(id int) (err error) {
	err = s.instancesRepo.RemoveItem(uint(id))
	return
}
