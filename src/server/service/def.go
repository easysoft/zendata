package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	serverUtils "github.com/easysoft/zendata/src/server/utils"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v3"
	"io/ioutil"
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

func (s *DefService) Get(id int) (def model.ZdDef, dirTree model.Dir) {
	def, _ = s.defRepo.Get(uint(id))

	dirTree = model.Dir{Name: fileUtils.AddSepIfNeeded(constant.ResDirUsers)}
	serverUtils.GetDirTree(&dirTree)

	return
}

func (s *DefService) Save(def *model.ZdDef) (err error) {
	def.Folder = serverUtils.DealWithPathSepRight(def.Folder)
	def.Path = vari.WorkDir + def.Folder + serverUtils.AddExt(def.Title, ".yaml")

	if def.ID == 0 {
		err = s.Create(def)
	} else {
		err = s.Update(def)
	}
	return
}

func (s *DefService) Create(def *model.ZdDef) (err error) {
	err = s.defRepo.Create(def)

	// add root field node
	rootField, err := s.fieldRepo.CreateTreeNode(def.ID, 0, "字段", "root")
	s.referRepo.CreateDefault(rootField.ID, constant.ResTypeDef)
	err = s.defRepo.Update(def)

	// update  yaml
	s.updateYaml(def.ID)

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
	s.updateYaml(def.ID)

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

	defData := model.DefData{}
	s.defRepo.GenDef(*def, &defData)

	for _, child := range root.Fields { // ignore the root
		defField := model.DefField{}
		convertToConfModel(*child, &defField)

		defData.Fields = append(defData.Fields, defField)
	}

	bytes, err := yaml.Marshal(defData)
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
		_, found := mp[fi.Path]
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.defRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		} else { // db is new

		}
	}

	return
}
func (s *DefService) saveFieldToDB(item *model.ZdField, parentID, defID uint) {
	item.DefID = defID
	item.ParentID = parentID
	s.fieldRepo.Save(item)

	for _, child := range item.Fields {
		s.saveFieldToDB(child, item.ID, defID)
	}
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
	po.Yaml = string(content)

	s.defRepo.Create(&po)

	rootField, _ := s.fieldRepo.CreateTreeNode(po.ID, 0, "字段", "root")
	s.referRepo.CreateDefault(rootField.ID, constant.ResTypeDef)
	for i, field := range po.Fields {
		field.Ord = i + 1
		s.saveFieldToDB(&field, rootField.ID, po.ID)
	}

	return
}

func NewDefService(defRepo *serverRepo.DefRepo, fieldRepo *serverRepo.FieldRepo,
	referRepo *serverRepo.ReferRepo) *DefService {
	return &DefService{defRepo: defRepo, fieldRepo: fieldRepo, referRepo: referRepo}
}
