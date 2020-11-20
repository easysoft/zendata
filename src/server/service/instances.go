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

type InstancesService struct {
	instancesRepo *serverRepo.InstancesRepo
	referRepo  *serverRepo.ReferRepo
	resService *ResService
}

func (s *InstancesService) List(keywords string, page int) (list []*model.ZdInstances, total int) {
	instances := s.resService.LoadRes("instances")
	list, _, _ = s.instancesRepo.List("", -1)

	s.importResToDB(instances, list)
	list, total, _ = s.instancesRepo.List(strings.TrimSpace(keywords), page)

	return
}

func (s *InstancesService) Get(id int) (instances model.ZdInstances, dirTree model.Dir) {
	instances, _ = s.instancesRepo.Get(uint(id))

	dirTree = model.Dir{Name: fileUtils.AddSepIfNeeded(constant.ResDirYaml)}
	serverUtils.GetDirTree(&dirTree)

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
			content, _ := ioutil.ReadFile(inst.Path)
			yamlContent := stringUtils.ReplaceSpecialChars(content)
			instPo := model.ZdInstances{}
			err = yaml.Unmarshal(yamlContent, &instPo)
			instPo.Title = inst.Title
			instPo.Name = inst.Name
			instPo.Desc = inst.Desc
			instPo.Path = inst.Path
			instPo.Folder = serverUtils.GetRelativePath(instPo.Path)
			instPo.Yaml = string(content)

			s.instancesRepo.Create(&instPo)

			for i, item := range instPo.Instances {
				item.Ord = i + 1
				s.saveItemToDB(&item, 0, instPo.ID)
			}
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
