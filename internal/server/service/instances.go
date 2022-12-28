package serverService

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/gen"
	"github.com/easysoft/zendata/internal/pkg/helper"
	"github.com/easysoft/zendata/internal/pkg/model"
	serverRepo "github.com/easysoft/zendata/internal/server/repo"
	serverUtils "github.com/easysoft/zendata/internal/server/utils"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type InstancesService struct {
	InstancesRepo  *serverRepo.InstancesRepo `inject:""`
	ReferRepo      *serverRepo.ReferRepo     `inject:""`
	ResService     *ResService               `inject:""`
	SectionService *SectionService           `inject:""`
	SectionRepo    *serverRepo.SectionRepo   `inject:""`
}

func (s *InstancesService) List(keywords string, page int) (list []*model.ZdInstances, total int) {
	list, total, _ = s.InstancesRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *InstancesService) Get(id int) (instances model.ZdInstances, dirs []model.Dir) {
	instances, _ = s.InstancesRepo.Get(uint(id))

	serverUtils.GetDirs(consts.ResDirYaml, &dirs)

	return
}

func (s *InstancesService) Save(instances *model.ZdInstances) (err error) {
	instances.Folder = serverUtils.DealWithPathSepRight(instances.Folder)
	instances.Path = vari.ZdPath + instances.Folder + serverUtils.AddExt(instances.FileName, ".yaml")
	instances.ReferName = helper.PathToName(instances.Path, consts.ResDirYaml, consts.ResTypeInstances)

	if instances.ID == 0 {
		err = s.Create(instances)
	} else {
		err = s.Update(instances)
	}
	s.updateYaml(instances.ID)

	return
}
func (s *InstancesService) Create(instances *model.ZdInstances) (err error) {
	err = s.InstancesRepo.Create(instances)

	return
}
func (s *InstancesService) Update(instances *model.ZdInstances) (err error) {
	var old model.ZdInstances
	old, err = s.InstancesRepo.Get(instances.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if instances.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.InstancesRepo.Update(instances)

	return
}

func (s *InstancesService) Remove(id int) (err error) {
	var old model.ZdInstances
	old, err = s.InstancesRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}

	fileUtils.RemoveExist(old.Path)
	err = s.InstancesRepo.Remove(uint(id))
	return
}

func (s *InstancesService) updateYamlByItem(itemId uint) (err error) {
	item, _ := s.InstancesRepo.GetItem(itemId)
	return s.updateYaml(item.InstancesID)
}
func (s *InstancesService) updateYaml(id uint) (err error) {
	var po model.ZdInstances
	po, _ = s.InstancesRepo.Get(id)

	s.genYaml(&po)
	err = s.InstancesRepo.UpdateYaml(po)
	fileUtils.WriteFile(po.Path, po.Yaml)

	return
}

func (s *InstancesService) genYaml(instances *model.ZdInstances) (str string) {
	//items, err := s.InstancesRepo.GetItems(instances.ID)
	//if err != nil {
	//	return
	//}
	//
	//yamlObj := model.ResInstances{}
	//s.InstancesRepo.GenInst(*instances, &yamlObj)
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
	instances.Yaml = helper.ConvertYamlStringToMapFormat(bytes)

	return
}

func (s *InstancesService) Sync(files []model.ResFile) {
	list := s.InstancesRepo.ListAll()

	mp := map[string]*model.ZdInstances{}
	for _, item := range list {
		mp[item.Path] = item
	}

	for _, fi := range files {
		_, found := mp[fi.Path]
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.InstancesRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		}
	}

	return
}

func (s *InstancesService) SyncToDB(fi model.ResFile) (err error) {
	content, _ := os.ReadFile(fi.Path)
	yamlContent := helper.ReplaceSpecialChars(content)
	po := model.ZdInstances{}
	err = yaml.Unmarshal(yamlContent, &po)

	po.Title = fi.Title
	po.Desc = fi.Desc
	po.Path = fi.Path
	po.Folder = serverUtils.GetRelativePath(po.Path)
	if strings.Index(po.Path, consts.ResDirYaml) > -1 {
		po.ReferName = helper.PathToName(po.Path, consts.ResDirYaml, consts.ResTypeInstances)
	} else {
		po.ReferName = helper.PathToName(po.Path, consts.ResDirUsers, consts.ResTypeInstances)
	}
	po.FileName = fileUtils.GetFileName(po.Path)
	po.Yaml = string(content)

	s.InstancesRepo.Create(&po)

	for i, item := range po.Instances {
		item.Ord = i + 1
		s.saveItemToDB(&item, po, fi.Path, 0, po.ID)
	}

	return
}
func (s *InstancesService) saveItemToDB(item *model.ZdInstancesItem, instances model.ZdInstances, currPath string, parentID, instancesID uint) {
	if item.Froms != nil && len(item.Froms) > 0 {
		for idx, from := range item.Froms {
			if from.Field == "" {
				from.Field = "from" + strconv.Itoa(idx+1)
			}
			s.saveItemToDB(from, instances, currPath, parentID, instancesID)
		}

		return
	}

	// update field
	if item.Instance != "" { // instance node
		item.Field = item.Instance
	}

	item.InstancesID = instancesID
	item.ParentID = parentID
	if item.From == "" && instances.From != "" {
		item.From = instances.From
	}
	if item.Type == "" {
		item.Type = consts.FieldTypeList
	}
	if item.Mode == "" {
		item.Mode = consts.ModeParallel
	}
	item.Range = strings.TrimSpace(item.Range)

	// save item
	s.InstancesRepo.SaveItem(item)

	// create refer
	refer := model.ZdRefer{OwnerType: "instances", OwnerID: item.ID}

	if strings.Index(currPath, "_test-instacnes.yaml") > -1 && item.Field == "field1-1" {
		log.Println("")
	}

	needToCreateSections := false
	if item.Select != "" { // refer to excel
		refer.Type = consts.ResTypeExcel

		refer.ColName = item.Select
		refer.Condition = item.Where
		refer.Rand = item.Rand

		_, sheet := fileUtils.ConvertResExcelPath(item.From, currPath)
		refer.File = item.From
		refer.Sheet = sheet

	} else if item.Use != "" { // refer to ranges or instances, need to read yaml to get the type
		rangeSections := gen.ParseRangeProperty(item.Use)
		if len(rangeSections) > 0 { // only get the first one
			rangeSection := rangeSections[0]
			desc, _, count, countTag := gen.ParseRangeSection(rangeSection) // medium{2}
			refer.ColName = desc
			refer.Count = count
			refer.CountTag = countTag
		}

		path := fileUtils.ConvertReferRangeToPath(item.From, currPath)
		_, _, refer.Type = helper.ReadYamlInfo(path)
		refer.File = item.From

	} else if item.Config != "" { // refer to config
		refer.Type = consts.ResTypeConfig

		rangeSections := gen.ParseRangeProperty(item.Config) // dir/config.yaml
		if len(rangeSections) > 0 {                          // only get the first one
			rangeSection := rangeSections[0]
			desc, _, count, countTag := gen.ParseRangeSection(rangeSection)
			refer.Count = count
			refer.CountTag = countTag

			path := fileUtils.ConvertReferRangeToPath(desc, currPath)
			refer.File = GetRelatedPathWithResDir(path)
		}

	} else if item.Range != "" { // deal with yaml and text refer using range prop

		rangeSections := gen.ParseRangeProperty(item.Range)
		if len(rangeSections) > 0 { // only get the first one
			rangeSection := rangeSections[0]
			desc, step, count, countTag := gen.ParseRangeSection(rangeSection) // dir/users.txt:R{3}
			if filepath.Ext(desc) == ".txt" || filepath.Ext(desc) == ".yaml" {
				if filepath.Ext(desc) == ".txt" { // dir/users.txt:2
					refer.Type = consts.ResTypeText

					if strings.ToLower(step) == "r" {
						refer.Rand = true
					} else {
						refer.Step, _ = strconv.Atoi(step)
					}

				} else if filepath.Ext(desc) == ".yaml" { // dir/content.yaml{3}
					refer.Type = consts.ResTypeYaml

					refer.Count = count
					refer.CountTag = countTag
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
	s.ReferRepo.Save(&refer)

	// gen sections if needed
	if needToCreateSections {
		rangeSections := gen.ParseRangeProperty(item.Range)

		for i, rangeSection := range rangeSections {
			s.SectionRepo.SaveFieldSectionToDB(rangeSection, i, item.ID, "instances")
		}
	}

	// deal with field's children
	for i, child := range item.Fields {
		child.Ord = i + 1
		s.saveItemToDB(child, instances, currPath, item.ID, instancesID)
	}
}

func NewInstancesService(instancesRepo *serverRepo.InstancesRepo, referRepo *serverRepo.ReferRepo,
	sectionRepo *serverRepo.SectionRepo) *InstancesService {
	return &InstancesService{InstancesRepo: instancesRepo, ReferRepo: referRepo, SectionRepo: sectionRepo}
}
