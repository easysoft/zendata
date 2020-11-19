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
)

type InstancesService struct {
	instancesRepo *serverRepo.InstancesRepo
	referRepo  *serverRepo.ReferRepo
	resService *ResService
}

func (s *InstancesService) List() (list []*model.ZdInstances) {
	instances := s.resService.LoadRes("instances")
	list, _ = s.instancesRepo.List()

	s.importResToDB(instances, list)
	list, _ = s.instancesRepo.List()

	return
}

func (s *InstancesService) Get(id int) (instances model.ZdInstances) {
	instances, _ = s.instancesRepo.Get(uint(id))

	return
}

func (s *InstancesService) Save(instances *model.ZdInstances) (err error) {
	instances.Folder = serverUtils.DealWithPathSepRight(instances.Folder)
	instances.Path = vari.WorkDir + instances.Folder + serverUtils.AddExt(instances.Title, ".yaml")
	instances.Name = service.PathToName(instances.Path, constant.ResDirYaml)

	if instances.ID == 0 {
		err = s.Create(instances)
	} else {
		err = s.Update(instances)
	}

	return
}
func (s *InstancesService) Create(instances *model.ZdInstances) (err error) {
	s.dataToYaml(instances)
	err = s.instancesRepo.Create(instances)

	return
}
func (s *InstancesService) Update(instances *model.ZdInstances) (err error) {
	var old model.ZdInstances
	old, err = s.instancesRepo.Get(instances.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if instances.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	s.dataToYaml(instances)
	err = s.instancesRepo.Update(instances)

	return
}

func (s *InstancesService) Remove(id int) (err error) {
	var old model.ZdInstances
	old, err = s.instancesRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}
	fileUtils.RemoveExist(old.Path)

	err = s.instancesRepo.Remove(uint(id))
	return
}

func (s *InstancesService) GetItemTree(instancesId int) (root model.ZdInstancesItem) {
	root = s.instancesRepo.GetItemTree(instancesId)
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

func (s *InstancesService) dataToYaml(inst *model.ZdInstances) (str string) {

	return
}

func (s *InstancesService) importResToDB(instances []model.ResFile, list []*model.ZdInstances) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, inst := range instances {
		if !stringUtils.FindInArrBool(inst.Path, names) {
			//if strings.Contains(inst.Path, "_test") {
				content, _ := ioutil.ReadFile(inst.Path)
				yamlContent := stringUtils.ReplaceSpecialChars(content)
				instPo := model.ZdInstances{}
				err = yaml.Unmarshal(yamlContent, &instPo)
				instPo.Title = inst.Title
				instPo.Name = inst.Name
				instPo.Desc = inst.Desc
				instPo.Path = inst.Path
				instPo.Yaml = string(content)

				s.instancesRepo.Create(&instPo)

				for i, item := range instPo.Instances {
					item.Ord = i + 1
					s.saveItemToDB(&item, 0, instPo.ID)
				}
			//}
		}
	}

	return
}
func (s *InstancesService) saveItemToDB(item *model.ZdInstancesItem, parentID, instancesID uint) {
	if item.Instance != "" { // instance node
		item.Field = item.Instance
	}

	item.InstancesID = instancesID
	item.ParentID = parentID
	s.instancesRepo.SaveItem(item)

	for i, child := range item.Fields {
		child.Ord = i + 1
		s.saveItemToDB(child, item.ID, instancesID)
	}
}

func NewInstancesService(instancesRepo *serverRepo.InstancesRepo, referRepo *serverRepo.ReferRepo) *InstancesService {
	return &InstancesService{instancesRepo: instancesRepo, referRepo: referRepo}
}
