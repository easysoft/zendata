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
	list, total, _ = s.instancesRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *InstancesService) Get(id int) (instances model.ZdInstances, dirs []model.Dir) {
	instances, _ = s.instancesRepo.Get(uint(id))

	serverUtils.GetDirs(constant.ResDirYaml, &dirs)

	return
}

func (s *InstancesService) Save(instances *model.ZdInstances) (err error) {
	instances.Folder = serverUtils.DealWithPathSepRight(instances.Folder)
	instances.Path = vari.WorkDir + instances.Folder + serverUtils.AddExt(instances.FileName, ".yaml")
	instances.ReferName = service.PathToName(instances.Path, constant.ResDirYaml)

	if instances.ID == 0 {
		err = s.Create(instances)
	} else {
		err = s.Update(instances)
	}
	s.updateYaml(instances.ID)

	return
}
func (s *InstancesService) Create(instances *model.ZdInstances) (err error) {
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

func (s *InstancesService) updateYaml(id uint) (err error) {
	var po model.ZdInstances
	po, _ = s.instancesRepo.Get(id)

	s.genYaml(&po)
	err = s.instancesRepo.UpdateYaml(po)
	fileUtils.WriteFile(po.Path, po.Yaml)

	return
}
func (s *InstancesService) genYaml(instances *model.ZdInstances) (str string) {
	//items, err := s.instancesRepo.GetItems(instances.ID)
	//if err != nil {
	//	return
	//}
	//
	//yamlObj := model.ResInstances{}
	//s.instancesRepo.GenInst(*instances, &yamlObj)
	//
	//for _, item := range items {
	//	inst := instancesItemToResInstForExport(*item)
	//
	//	yamlObj.Instances = append(yamlObj.Instances, inst)
	//}

	root := s.GetItemTree(instances.ID)
	for _, item := range root.Fields {
		instances.Instances = append(instances.Instances, *item)
	}
	bytes, _ := yaml.Marshal(instances)
	instances.Yaml = stringUtils.ConvertYamlStringToMapFormat(bytes)

	return
}

func (s *InstancesService) Sync(files []model.ResFile) {
	list := s.instancesRepo.ListAll()

	mp := map[string]*model.ZdInstances{}
	for _, item := range list {
		mp[item.Path] = item
	}

	for _, fi := range files {
		_, found := mp[fi.Path]
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.instancesRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		}
	}

	return
}

func (s *InstancesService) SyncToDB(file model.ResFile) (err error) {
	content, _ := ioutil.ReadFile(file.Path)
	yamlContent := stringUtils.ReplaceSpecialChars(content)
	po := model.ZdInstances{}
	err = yaml.Unmarshal(yamlContent, &po)

	po.Title = file.Title
	po.Desc = file.Desc
	po.Path = file.Path
	po.Folder = serverUtils.GetRelativePath(po.Path)
	if strings.Index(po.Path, constant.ResDirYaml) > -1 {
		po.ReferName = service.PathToName(po.Path, constant.ResDirYaml)
	} else {
		po.ReferName = service.PathToName(po.Path, constant.ResDirUsers)
	}
	po.FileName = fileUtils.GetFileName(po.Path)
	po.Yaml = string(content)

	s.instancesRepo.Create(&po)

	for i, item := range po.Instances {
		item.Ord = i + 1
		s.saveItemToDB(&item, 0, po.ID)
	}

	return
}
func (s *InstancesService) saveItemToDB(item *model.ZdInstancesItem, parentID, instancesID uint) {
	if item.Instance != "" { // instance node
		item.Field = item.Instance
	}

	item.InstancesID = instancesID
	item.ParentID = parentID
	if item.Type == "" {
		item.Type = constant.FieldTypeList
	}
	if item.Mode == "" {
		item.Mode = constant.ModeParallel
	}

	s.instancesRepo.SaveItem(item)

	for i, child := range item.Fields {
		child.Ord = i + 1
		s.saveItemToDB(child, item.ID, instancesID)
	}
}

func NewInstancesService(instancesRepo *serverRepo.InstancesRepo, referRepo *serverRepo.ReferRepo) *InstancesService {
	return &InstancesService{instancesRepo: instancesRepo, referRepo: referRepo}
}
