package serverService

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/easysoft/zendata/src/gen"
	"github.com/easysoft/zendata/src/model"
	serverRepo "github.com/easysoft/zendata/src/server/repo"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	"github.com/easysoft/zendata/src/service"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type DefService struct {
	DefRepo     *serverRepo.DefRepo     `inject:""`
	FieldRepo   *serverRepo.FieldRepo   `inject:""`
	ReferRepo   *serverRepo.ReferRepo   `inject:""`
	SectionRepo *serverRepo.SectionRepo `inject:""`

	ResService     *ResService     `inject:""`
	SectionService *SectionService `inject:""`
}

func (s *DefService) List(keywords string, page int) (list []*model.ZdDef, total int) {
	list, total, _ = s.DefRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *DefService) Get(id int) (def model.ZdDef, dirs []model.Dir) {
	if id > 0 {
		def, _ = s.DefRepo.Get(uint(id))
	} else {
		def.Folder = "users" + constant.PthSep
		def.Type = "text"
	}

	serverUtils.GetDirs(constant.ResDirUsers, &dirs)

	return
}

func (s *DefService) Save(def *model.ZdDef) (err error) {
	def.Folder = serverUtils.DealWithPathSepRight(def.Folder)
	def.Path = vari.ZdPath + def.Folder + serverUtils.AddExt(def.FileName, ".yaml")

	if def.ID == 0 {
		err = s.Create(def)
	} else {
		err = s.Update(def)
	}
	s.updateYaml(def.ID)

	return
}

func (s *DefService) Create(def *model.ZdDef) (err error) {
	def.ReferName = service.PathToName(def.Path, constant.ResDirUsers, def.Type)
	err = s.DefRepo.Create(def)

	// add root field node
	rootField, err := s.FieldRepo.CreateTreeNode(def.ID, 0, "字段", "root")
	s.ReferRepo.CreateDefault(rootField.ID, constant.ResTypeDef)
	err = s.DefRepo.Update(def)

	return
}

func (s *DefService) Update(def *model.ZdDef) (err error) {
	var old model.ZdDef
	old, err = s.DefRepo.Get(def.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if def.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	def.ReferName = service.PathToName(def.Path, constant.ResDirUsers, def.Type)
	err = s.DefRepo.Update(def)

	return
}

func (s *DefService) Remove(id int) (err error) {
	var old model.ZdDef
	old, err = s.DefRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}
	fileUtils.RemoveExist(old.Path)

	err = s.DefRepo.Remove(uint(id))
	return
}

func (s *DefService) updateYamlByField(fieldId uint) (err error) {
	field, _ := s.FieldRepo.Get(fieldId)
	return s.updateYaml(field.DefID)
}

func (s *DefService) updateYaml(id uint) (err error) {
	var po model.ZdDef
	po, _ = s.DefRepo.Get(id)

	s.genYaml(&po)
	err = s.DefRepo.UpdateYaml(po)
	fileUtils.WriteFile(po.Path, po.Yaml)

	return
}

func (s *DefService) genYaml(def *model.ZdDef) (str string) {
	root, err := s.FieldRepo.GetDefFieldTree(def.ID)
	if err != nil {
		return
	}

	yamlObj := model.DefData{}
	s.DefRepo.GenDef(*def, &yamlObj)

	for _, child := range root.Fields { // ignore the root
		defField := model.DefField{}

		refer, _ := s.ReferRepo.GetByOwnerId(child.ID)
		s.zdFieldToFieldForExport(*child, refer, &defField)

		yamlObj.Fields = append(yamlObj.Fields, defField)
	}

	bytes, err := yaml.Marshal(yamlObj)
	def.Yaml = stringUtils.ConvertYamlStringToMapFormat(bytes)

	return
}

func (s *DefService) zdFieldToFieldForExport(treeNode model.ZdField, refer model.ZdRefer, field *model.DefField) {
	genFieldFromZdField(treeNode, refer, field)

	for _, child := range treeNode.Fields {
		childField := model.DefField{}

		childRefer, _ := s.ReferRepo.GetByOwnerId(child.ID)
		s.zdFieldToFieldForExport(*child, childRefer, &childField)

		field.Fields = append(field.Fields, childField)
	}

	//for _, from := range treeNode.Froms { // only one level
	//	childField := model.DefField{}
	//	genFieldFromZdField(*from, &childField)
	//
	//	field.Froms = append(field.Froms, childField)
	//}

	if len(field.Fields) == 0 {
		field.Fields = nil
	}
	if len(field.Froms) == 0 {
		field.Froms = nil
	}

	return
}

func (s *DefService) Sync(files []model.ResFile) (err error) {
	list := s.DefRepo.ListAll()

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
			s.DefRepo.Remove(mp[fi.Path].ID)
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

	s.DefRepo.Create(&po)

	rootField, _ := s.FieldRepo.CreateTreeNode(po.ID, 0, "字段", "root")
	s.ReferRepo.CreateDefault(rootField.ID, constant.ResTypeDef)
	for i, field := range po.Fields {
		field.Ord = i + 1
		s.saveFieldToDB(&field, po, fi.Path, rootField.ID, po.ID)
	}

	return
}
func (s *DefService) saveFieldToDB(field *model.ZdField, def model.ZdDef, currPath string, parentID, defID uint) {
	if field.Froms != nil && len(field.Froms) > 0 {
		for idx, from := range field.Froms {
			if from.Field == "" {
				from.Field = "from" + strconv.Itoa(idx+1)
			}
			s.saveFieldToDB(from, def, currPath, parentID, defID)
		}

		return
	}

	// update field
	field.DefID = defID
	field.ParentID = parentID
	if field.From == "" && def.From != "" {
		field.From = def.From
	}
	if field.Type == "" {
		field.Type = constant.FieldTypeList
	}
	if field.Mode == "" {
		field.Mode = constant.ModeParallel
	}

	field.Range = strings.TrimSpace(field.Range)

	// save field
	s.FieldRepo.Save(field)

	// create refer
	refer := model.ZdRefer{OwnerType: "def", OwnerID: field.ID}

	needToCreateSections := false
	if field.Select != "" { // refer to excel
		refer.Type = constant.ResTypeExcel

		refer.ColName = field.Select
		refer.Condition = field.Where
		refer.Rand = field.Rand

		_, sheet := fileUtils.ConvertResExcelPath(field.From, currPath)
		refer.File = field.From
		refer.Sheet = sheet

	} else if field.Use != "" { // refer to ranges or instances, need to read yaml to get the type
		rangeSections := gen.ParseRangeProperty(field.Use)
		if len(rangeSections) > 0 { // only get the first one
			rangeSection := rangeSections[0]
			desc, _, count, countTag := gen.ParseRangeSection(rangeSection) // medium{2!}
			refer.ColName = desc
			refer.Count = count
			refer.CountTag = countTag
		}

		path := fileUtils.ConvertReferRangeToPath(field.From, currPath)
		_, _, refer.Type = service.ReadYamlInfo(path)
		refer.File = field.From

	} else if field.Config != "" { // refer to config
		refer.Type = constant.ResTypeConfig

		rangeSections := gen.ParseRangeProperty(field.Config) // dir/config.yaml
		if len(rangeSections) > 0 {                           // only get the first one
			rangeSection := rangeSections[0]
			desc, _, count, countTag := gen.ParseRangeSection(rangeSection)
			refer.Count = count
			refer.CountTag = countTag

			path := fileUtils.ConvertReferRangeToPath(desc, currPath)
			refer.File = GetRelatedPathWithResDir(path)
		}

	} else if field.Range != "" {
		rangeSections := gen.ParseRangeProperty(field.Range)
		if len(rangeSections) > 0 {
			rangeSection := rangeSections[0]                                   // deal with yaml and text refer using range prop
			desc, step, count, countTag := gen.ParseRangeSection(rangeSection) // dir/users.txt:R{3}
			if filepath.Ext(desc) == ".txt" || filepath.Ext(desc) == ".yaml" {
				if filepath.Ext(desc) == ".txt" { // dir/users.txt:2
					refer.Type = constant.ResTypeText

					if strings.ToLower(step) == "r" {
						refer.Rand = true
					} else {
						refer.Step, _ = strconv.Atoi(step)
					}

				} else if filepath.Ext(desc) == ".yaml" { // dir/content.yaml{3}
					refer.Type = constant.ResTypeYaml

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
	refer.OwnerID = field.ID
	s.ReferRepo.Save(&refer)

	// gen sections if needed
	if needToCreateSections {
		rangeSections := gen.ParseRangeProperty(field.Range)

		for i, rangeSection := range rangeSections {
			s.SectionRepo.SaveFieldSectionToDB(rangeSection, i, field.ID, "def")
		}
	}

	// deal with field's children
	for _, child := range field.Fields {
		s.saveFieldToDB(child, def, currPath, field.ID, defID)
	}
}
