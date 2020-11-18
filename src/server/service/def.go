package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v3"
	"strings"
)

type DefService struct {
	defRepo *serverRepo.DefRepo
	fieldRepo *serverRepo.FieldRepo
	referRepo *serverRepo.ReferRepo
}

func (s *DefService) List() (defs []*model.ZdDef) {
	defs, _ = s.defRepo.List()
	return
}

func (s *DefService) Get(id int) (def model.ZdDef, err error) {
	def, _ = s.defRepo.Get(uint(id))
	def.Folder = s.getFolder(def.Path)

	return
}

func (s *DefService) Create(def *model.ZdDef) (err error) {
	def.Folder = s.dealWithPathSepRight(def.Folder)

	def.Path = def.Folder + def.Title
	def.Path = s.addExt(def.Path)
	err = s.defRepo.Create(def)

	rootField, err := s.fieldRepo.CreateTreeNode(def.ID, 0, "字段", "root")
	s.referRepo.CreateDefault(rootField.ID, constant.ResTypeDef)

	s.dataToYaml(def)
	err = s.defRepo.Update(def)

	return
}

func (s *DefService) Update(def *model.ZdDef) (err error) {
	def.Folder = s.dealWithPathSepRight(def.Folder)

	def.Path = def.Folder + def.Title
	def.Path = s.addExt(def.Path)

	var oldDef model.ZdDef
	oldDef, err = s.defRepo.Get(def.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}

	if def.Path != oldDef.Path {
		fileUtils.RemoveExist(oldDef.Path)
	}

	s.dataToYaml(def)
	err = s.defRepo.Update(def)

	return
}

func (s *DefService) Remove(id int) (err error) {
	var oldDef model.ZdDef
	oldDef, err = s.defRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}

	fileUtils.RemoveExist(oldDef.Path)

	var def model.ZdDef
	def.ID = uint(id)
	err = s.defRepo.Remove(uint(id))
	return
}

func (s *DefService) UpdateYaml(defId uint) (err error) {
	var def model.ZdDef
	def, err = s.Get(int(defId))

	s.dataToYaml(&def)
	err = s.defRepo.UpdateYaml(def)

	return
}

func (s *DefService) dataToYaml(def *model.ZdDef) (str string) {
	root, err := s.fieldRepo.GetDefFieldTree(def.ID)
	if err != nil {
		return
	}

	defData := model.DefData{}
	s.defRepo.GenDef(*def, &defData)

	for _, child := range root.Children { // ignore the root
		defField := model.DefField{}
		convertToConfModel(*child, &defField)

		defData.Fields = append(defData.Fields, defField)
	}

	bytes, err := yaml.Marshal(defData)
	def.Yaml = string(bytes)

	return
}

func (s *DefService) addExt(pth string) string {
	if strings.LastIndex(pth, ".yaml") != len(pth) - 4 {
		pth += ".yaml"
	}

	return pth
}

func (s *DefService) dealWithPathSepRight(pth string) string {
	pth = fileUtils.RemovePathSepLeftIfNeeded(pth)
	pth = fileUtils.AddPathSepRightIfNeeded(pth)

	return pth
}
func (s *DefService) getFolder(pth string) string {
	idx := strings.LastIndex(pth, constant.PthSep)
	return pth[:idx+1]
}

func NewDefService(defRepo *serverRepo.DefRepo, fieldRepo *serverRepo.FieldRepo, referRepo *serverRepo.ReferRepo) *DefService {
	return &DefService{defRepo: defRepo, fieldRepo: fieldRepo, referRepo: referRepo}
}
