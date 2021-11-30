package serverService

import (
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
)

func (s *InstancesService) GetItemTree(instancesId uint) (root model.ZdInstancesItem) {
	root = s.InstancesRepo.GetItemTree(uint(instancesId))
	return
}
func (s *InstancesService) GetItem(id int) (item model.ZdInstancesItem) {
	item, _ = s.InstancesRepo.GetItem(uint(id))
	return
}

func (s *InstancesService) CreateItem(domainId, targetId int, mode string) (item *model.ZdInstancesItem, err error) {
	item = &model.ZdInstancesItem{InstancesID: uint(domainId)}
	item.Field = "instances_"

	item.Ord = s.InstancesRepo.GetMaxOrder(domainId)

	err = s.InstancesRepo.SaveItem(item)
	s.ReferRepo.CreateDefault(item.ID, constant.ResTypeInstances)

	return
}
func (s *InstancesService) SaveItem(item *model.ZdInstancesItem) (err error) {
	err = s.InstancesRepo.SaveItem(item)
	return
}

func (s *InstancesService) RemoveItem(id int) (err error) {
	err = s.InstancesRepo.RemoveItem(uint(id))
	return
}
