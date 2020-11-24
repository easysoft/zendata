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
	"path"
	"strings"
)

type DefService struct {
	defRepo *serverRepo.DefRepo
	fieldRepo *serverRepo.FieldRepo
	referRepo *serverRepo.ReferRepo
	resService *ResService
}

func (s *DefService) List(keywords string, page int) (list []*model.ZdDef, total int) {
	list, total, _ = s.defRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *DefService) Get(id int) (def model.ZdDef, dirs []model.Dir) {
	def, _ = s.defRepo.Get(uint(id))

	serverUtils.GetDirs(constant.ResDirUsers, &dirs)

	return
}

func (s *DefService) Save(def *model.ZdDef) (err error) {
	def.Folder = serverUtils.DealWithPathSepRight(def.Folder)
	def.Path = vari.WorkDir + def.Folder + serverUtils.AddExt(def.FileName, ".yaml")

	if def.ID == 0 {
		err = s.Create(def)
	} else {
		err = s.Update(def)
	}
	s.updateYaml(def.ID)

	return
}

func (s *DefService) Create(def *model.ZdDef) (err error) {
	err = s.defRepo.Create(def)

	// add root field node
	rootField, err := s.fieldRepo.CreateTreeNode(def.ID, 0, "字段", "root")
	s.referRepo.CreateDefault(rootField.ID, constant.ResTypeDef)
	err = s.defRepo.Update(def)

	return
}

func (s *DefService) Update(def *model.ZdDef) (err error) {
	var old model.ZdDef
	old, err = s.defRepo.Get(def.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if def.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.defRepo.Update(def)

	return
}

func (s *DefService) Remove(id int) (err error) {
	var old model.ZdDef
	old, err = s.defRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}
	fileUtils.RemoveExist(old.Path)

	err = s.defRepo.Remove(uint(id))
	return
}

func (s *DefService) updateYaml(id uint) (err error) {
	var po model.ZdDef
	po, _ = s.defRepo.Get(id)

	s.genYaml(&po)
	err = s.defRepo.UpdateYaml(po)
	fileUtils.WriteFile(po.Path, po.Yaml)

	return
}

func (s *DefService) genYaml(def *model.ZdDef) (str string) {
	root, err := s.fieldRepo.GetDefFieldTree(def.ID)
	if err != nil {
		return
	}

	yamlObj := model.DefData{}
	s.defRepo.GenDef(*def, &yamlObj)

	for _, child := range root.Fields { // ignore the root
		defField := model.DefField{}
		zdFieldToFieldForExport(*child, &defField)

		yamlObj.Fields = append(yamlObj.Fields, defField)
	}

	bytes, err := yaml.Marshal(yamlObj)
	def.Yaml = stringUtils.ConvertYamlStringToMapFormat(bytes)

	return
}

func (s *DefService) Sync(files []model.ResFile) (err error) {
	list := s.defRepo.ListAll()

	mp := map[string]*model.ZdDef{}
	for _, item := range list {
		mp[item.Path] = item
	}

	for _, fi := range files {
		// for yaml "res", "data" type should be default value text
		if fi.ResType == "" || fi.ResType == constant.ResTypeYaml {
			fi.ResType = constant.ResTypeText
		}

		_, found := mp[fi.Path]
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.defRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		}
	}

	return
}
func (s *DefService) SyncToDB(fi model.ResFile) (err error) {
	content, _ := ioutil.ReadFile(fi.Path)
	yamlContent := stringUtils.ReplaceSpecialChars(content)
	po := model.ZdDef{}
	err = yaml.Unmarshal(yamlContent, &po)
	po.Title = fi.Title
	po.Type = fi.ResType
	po.Desc = fi.Desc
	po.Path = fi.Path
	po.Folder = serverUtils.GetRelativePath(po.Path)

	po.ReferName = service.PathToName(po.Path, constant.ResDirUsers, po.Type)
	po.FileName = fileUtils.GetFileName(po.Path)

	po.Yaml = string(content)

	s.defRepo.Create(&po)

	rootField, _ := s.fieldRepo.CreateTreeNode(po.ID, 0, "字段", "root")
	s.referRepo.CreateDefault(rootField.ID, constant.ResTypeDef)
	for i, field := range po.Fields {
		field.Ord = i + 1
		s.saveFieldToDB(&field, fi.Path, rootField.ID, po.ID)
	}

	return
}
func (s *DefService) saveFieldToDB(item *model.ZdField, currPath string, parentID, defID uint) {
	refer := model.ZdRefer{OwnerType: "def", OwnerID: item.ID}

	if item.Config != "" { // refer to config
		refer.Type = constant.ResTypeConfig

		item.Config = strings.TrimSpace(item.Config)
		rangeSections := gen.ParseRangeProperty(item.Config)
		if len(rangeSections) == 1 {
			rangeSection := rangeSections[0]
			desc, _, count := gen.ParseRangeSection(rangeSection)
			refer.Count = count

			path := FileToPath(desc, currPath)
			refer.File = strings.Replace(path, vari.WorkDir, "", 1)
		}

	} else if item.Select != "" { // refer to excel
		refer.Type = constant.ResTypeExcel

		refer.ColName = item.Select
		refer.Condition = item.Where
		refer.Rand = item.Rand

		_, sheet := fileUtils.ConvertResExcelPath(item.From)
		refer.File = item.From
		refer.Sheet = sheet

	} else if item.Use != "" { // refer to ranges or instances
		refer.File = item.From
		refer.ColName, _, refer.Count = gen.ParseRangeSection(item.Use)

		path := FileToPath(item.From, currPath)
		_, _, refer.Type = service.ReadYamlInfo(path)

	} else if item.Range != "" { // deal with yaml and text refer using range prop
		item.Range = strings.TrimSpace(item.Range)
		rangeSections := gen.ParseRangeProperty(item.Range)
		if len(rangeSections) == 1 {
			rangeSection := rangeSections[0]
			desc, _, count := gen.ParseRangeSection(rangeSection)

			if path.Ext(desc) == ".txt" {
				refer.Type = constant.ResTypeText
			} else if path.Ext(desc) == ".yaml" {
				refer.Type = constant.ResTypeYaml
			}
			if path.Ext(desc) == ".txt" || path.Ext(desc) == ".yaml" {
				refer.Count = count

				path := FileToPath(desc, currPath)
				refer.File = strings.Replace(path, vari.WorkDir, "", 1)
			}
		}
	}

	item.DefID = defID
	item.ParentID = parentID
	if item.Type == "" {
		item.Type = constant.FieldTypeList
	}
	if item.Mode == "" {
		item.Mode = constant.ModeParallel
	}
	s.fieldRepo.Save(item)
	refer.OwnerID = item.ID
	s.referRepo.Save(&refer)

	for _, child := range item.Fields {
		s.saveFieldToDB(child, currPath, item.ID, defID)
	}
}

func NewDefService(defRepo *serverRepo.DefRepo, fieldRepo *serverRepo.FieldRepo,
	referRepo *serverRepo.ReferRepo) *DefService {
	return &DefService{defRepo: defRepo, fieldRepo: fieldRepo, referRepo: referRepo}
}
