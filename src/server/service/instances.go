package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
)

type InstancesService struct {
	instancesRepo *serverRepo.InstancesRepo
	resService *ResService
}

func (s *InstancesService) List() (list []*model.ZdInstances) {
	instances := s.resService.LoadRes("instances")
	list, _ = s.instancesRepo.List()

	s.saveResToDB(instances, list)
	list, _ = s.instancesRepo.List()

	return
}

func (s *InstancesService) Get(id int) (instances model.ZdInstances) {
	instances, _ = s.instancesRepo.Get(uint(id))

	return
}

func (s *InstancesService) Save(instances *model.ZdInstances) (err error) {
	err = s.instancesRepo.Save(instances)

	return
}

func (s *InstancesService) Remove(id int) (err error) {
	err = s.instancesRepo.Remove(uint(id))
	if err != nil {
		return
	}

	instances, _ := s.instancesRepo.Get(uint(id))
	logUtils.PrintTo(instances.Path)
	//fileUtils.RemoveExist(instances.Path)

	return
}

func (s *InstancesService) GetItemTree(rangesId int) (root model.ZdInstancesItem) {
	items, _ := s.instancesRepo.GetItems(rangesId)

	root.ID = 0
	root.Note = "实例"
	for _, item := range items {
		item.ParentID = root.ID
		root.Children = append(root.Children, item)
	}

	return
}
func (s *InstancesService) GetItem(id int) (item model.ZdInstancesItem) {
	item, _ = s.instancesRepo.GetItem(uint(id))
	return
}

func (s *InstancesService) CreateItem(domainId, targetId int, mode string) (item *model.ZdInstancesItem, err error) {
	item = &model.ZdInstancesItem{Field: "instances_", InstancesID: uint(domainId)}
	item.Ord = s.instancesRepo.GetMaxOrder(domainId)

	err = s.instancesRepo.SaveItem(item)

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

func (s *InstancesService) saveResToDB(instances []model.ResFile, list []*model.ZdInstances) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range instances {
		if !stringUtils.FindInArrBool(item.Path, names) {
			instances := model.ZdInstances{Title: item.Title, Name: item.Name, Desc: item.Desc, Path: item.Path}
			s.instancesRepo.Save(&instances)
		}
	}

	return
}

func NewInstancesService(instancesRepo *serverRepo.InstancesRepo) *InstancesService {
	return &InstancesService{instancesRepo: instancesRepo}
}
