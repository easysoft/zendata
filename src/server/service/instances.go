package serverService

import (
	"github.com/easysoft/zendata/src/gen"
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
	"log"
	"path"
	"strconv"
	"strings"
)

type InstancesService struct {
	instancesRepo  *serverRepo.InstancesRepo
	referRepo      *serverRepo.ReferRepo
	resService     *ResService
	sectionService *SectionService
	sectionRepo    *serverRepo.SectionRepo
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
	instances.ReferName = service.PathToName(instances.Path, constant.ResDirYaml, constant.ResTypeInstances)

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

func (s *InstancesService) updateYamlByItem(itemId uint) (err error) {
	item, _ := s.instancesRepo.GetItem(itemId)
	return s.updateYaml(item.InstancesID)
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

func (s *InstancesService) SyncToDB(fi model.ResFile) (err error) {
	content, _ := ioutil.ReadFile(fi.Path)
	yamlContent := stringUtils.ReplaceSpecialChars(content)
	po := model.ZdInstances{}
	err = yaml.Unmarshal(yamlContent, &po)

	po.Title = fi.Title
	po.Desc = fi.Desc
	po.Path = fi.Path
	po.Folder = serverUtils.GetRelativePath(po.Path)
	if strings.Index(po.Path, constant.ResDirYaml) > -1 {
		po.ReferName = service.PathToName(po.Path, constant.ResDirYaml, constant.ResTypeInstances)
	} else {
		po.ReferName = service.PathToName(po.Path, constant.ResDirUsers, constant.ResTypeInstances)
	}
	po.FileName = fileUtils.GetFileName(po.Path)
	po.Yaml = string(content)

	s.instancesRepo.Create(&po)

	for i, item := range po.Instances {
		item.Ord = i + 1
		s.saveItemToDB(&item, fi.Path, 0, po.ID)
	}

	return
}
func (s *InstancesService) saveItemToDB(item *model.ZdInstancesItem, currPath string, parentID, instancesID uint) {
	if item.Froms != nil && len(item.Froms) > 0 {
		for _, from := range item.Froms {
			s.saveItemToDB(from, currPath, parentID, instancesID)
		}

		return
	}

	// update field
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
	item.Range = strings.TrimSpace(item.Range)

	// save item
	s.instancesRepo.SaveItem(item)

	// create refer
	refer := model.ZdRefer{OwnerType: "instances", OwnerID: item.ID}

	if strings.Index(currPath, "_test-instacnes.yaml") > -1 && item.Field == "field1-1" {
		log.Println("")
	}

	needToCreateSections := false
	if item.Select != "" { // refer to excel
		refer.Type = constant.ResTypeExcel

		refer.ColName = item.Select
		refer.Condition = item.Where
		refer.Rand = item.Rand

		_, sheet := fileUtils.ConvertResExcelPath(item.From)
		refer.File = item.From
		refer.Sheet = sheet

	} else if item.Use != "" { // refer to ranges or instances, need to read yaml to get the type
		rangeSections := gen.ParseRangeProperty(item.Use)
		if len(rangeSections) > 0 { // only get the first one
			rangeSection := rangeSections[0]
			desc, _, count := gen.ParseRangeSection(rangeSection) // medium{2}
			refer.ColName = desc
			refer.Count = count
		}

		path := fileUtils.ConvertReferRangeToPath(item.From, currPath)
		_, _, refer.Type = service.ReadYamlInfo(path)
		refer.File = item.From

	} else if item.Config != "" { // refer to config
		refer.Type = constant.ResTypeConfig

		rangeSections := gen.ParseRangeProperty(item.Config) // dir/config.yaml
		if len(rangeSections) > 0 {                          // only get the first one
			rangeSection := rangeSections[0]
			desc, _, count := gen.ParseRangeSection(rangeSection)
			refer.Count = count

			path := fileUtils.ConvertReferRangeToPath(desc, currPath)
			refer.File = GetRelatedPathWithResDir(path)
		}

	} else if item.Range != "" { // deal with yaml and text refer using range prop

		rangeSections := gen.ParseRangeProperty(item.Range)
		if len(rangeSections) > 0 { // only get the first one
			rangeSection := rangeSections[0]
			desc, step, count := gen.ParseRangeSection(rangeSection) // dir/users.txt:R{3}
			if path.Ext(desc) == ".txt" || path.Ext(desc) == ".yaml" {
				if path.Ext(desc) == ".txt" { // dir/users.txt:2
					refer.Type = constant.ResTypeText

					if strings.ToLower(step) == "r" {
						refer.Rand = true
					} else {
						refer.Step, _ = strconv.Atoi(step)
					}

				} else if path.Ext(desc) == ".yaml" { // dir/content.yaml{3}
					refer.Type = constant.ResTypeYaml

					refer.Count = count
				}

				path := fileUtils.ConvertReferRangeToPath(desc, currPath)
				refer.File = GetRelatedPathWithResDir(path)
			} else { // like 1-9,a-z
				needToCreateSections = true
			}
		}
	}

	// save refer
	refer.OwnerID = item.ID
	s.referRepo.Save(&refer)

	// gen sections if needed
	if needToCreateSections {
		rangeSections := gen.ParseRangeProperty(item.Range)

		for i, rangeSection := range rangeSections {
			s.sectionRepo.SaveFieldSectionToDB(rangeSection, i, item.ID, "instances")
		}
	}

	// deal with field's children
	for i, child := range item.Fields {
		child.Ord = i + 1
		s.saveItemToDB(child, currPath, item.ID, instancesID)
	}
}

func NewInstancesService(instancesRepo *serverRepo.InstancesRepo, referRepo *serverRepo.ReferRepo,
	sectionRepo *serverRepo.SectionRepo) *InstancesService {
	return &InstancesService{instancesRepo: instancesRepo, referRepo: referRepo, sectionRepo: sectionRepo}
}
