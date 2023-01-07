package serverService

import (
	"os"
	"regexp"
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

type ConfigService struct {
	ConfigRepo  *serverRepo.ConfigRepo  `inject:""`
	ResService  *ResService             `inject:""`
	SectionRepo *serverRepo.SectionRepo `inject:""`
}

func (s *ConfigService) List(keywords string, page int) (list []*model.ZdConfig, total int) {
	list, total, _ = s.ConfigRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *ConfigService) Get(id int) (config model.ZdConfig, dirs []model.Dir) {
	if id > 0 {
		config, _ = s.ConfigRepo.Get(uint(id))
	}

	serverUtils.GetDirs(consts.ResDirYaml, &dirs)

	return
}

func (s *ConfigService) Save(config *model.ZdConfig) (err error) {
	config.Folder = serverUtils.DealWithPathSepRight(config.Folder)
	config.Path = vari.ZdDir + config.Folder + serverUtils.AddExt(config.FileName, ".yaml")
	config.ReferName = helper.PathToName(config.Path, consts.ResDirYaml, consts.ResTypeConfig)

	if config.ID == 0 {
		err = s.Create(config)
	} else {
		err = s.Update(config)
	}

	return
}

func (s *ConfigService) Create(config *model.ZdConfig) (err error) {
	err = s.ConfigRepo.Create(config)
	s.updateYaml(config.ID)

	return
}

func (s *ConfigService) Update(config *model.ZdConfig) (err error) {
	var old model.ZdConfig
	old, err = s.ConfigRepo.Get(config.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if config.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.ConfigRepo.Update(config)
	s.updateYaml(config.ID)

	return
}

func (s *ConfigService) Remove(id int) (err error) {
	var old model.ZdConfig
	old, err = s.ConfigRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}

	fileUtils.RemoveExist(old.Path)
	err = s.ConfigRepo.Remove(uint(id))

	return
}

func (s *ConfigService) updateYaml(id uint) (err error) {
	var po model.ZdConfig
	po, _ = s.ConfigRepo.Get(id)

	s.genYaml(&po)
	err = s.ConfigRepo.UpdateYaml(po)
	fileUtils.WriteFile(po.Path, po.Yaml)

	return
}
func (s *ConfigService) genYaml(config *model.ZdConfig) (str string) {
	yamlObj := model.ResConfig{}
	s.ConfigRepo.GenConfigRes(*config, &yamlObj)

	bytes, _ := yaml.Marshal(yamlObj)
	config.Yaml = helper.ConvertYamlStringToMapFormat(bytes)

	return
}

func (s *ConfigService) Sync(files []model.ResFile) (err error) {
	list := s.ConfigRepo.ListAll()

	mp := map[string]*model.ZdConfig{}
	for _, item := range list {
		mp[item.Path] = item
	}

	for _, fi := range files {
		_, found := mp[fi.Path]
		//logUtils.PrintTo(fi.UpdatedAt.String() + ", " + mp[fi.Path].UpdatedAt.String())
		if !found { // no record
			s.SyncToDB(fi)
		} else if fi.UpdatedAt.Unix() > mp[fi.Path].UpdatedAt.Unix() { // db is old
			s.ConfigRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		} else { // db is new

		}
	}

	return
}
func (s *ConfigService) SyncToDB(fi model.ResFile) (err error) {
	content, _ := os.ReadFile(fi.Path)
	yamlContent := helper.ReplaceSpecialChars(content)

	po := model.ZdConfig{}
	err = yaml.Unmarshal(yamlContent, &po)

	po.Title = fi.Title
	po.Desc = fi.Desc
	po.Path = fi.Path
	po.Folder = serverUtils.GetRelativePath(po.Path)
	if strings.Index(po.Path, vari.ZdDir+consts.ResDirYaml) > -1 {
		po.ReferName = helper.PathToName(po.Path, consts.ResDirYaml, consts.ResTypeConfig)
	} else {
		po.ReferName = helper.PathToName(po.Path, consts.ResDirUsers, consts.ResTypeConfig)
	}
	po.FileName = fileUtils.GetFileName(po.Path)
	po.Yaml = string(content)

	reg := regexp.MustCompile(`\t`)
	if reg.MatchString(po.Prefix) {
		po.Prefix = strings.ReplaceAll(po.Prefix, "\t", `\t`)
		po.Prefix = `"` + po.Prefix + `"`
	}

	s.ConfigRepo.Create(&po)

	rangeSections := gen.ParseRangeProperty(po.Range)
	for i, rangeSection := range rangeSections {
		s.SectionRepo.SaveFieldSectionToDB(rangeSection, i, po.ID, "config")
	}

	return
}

func (s *ConfigService) GConfigItemTree(configId int) (root model.ZdRangesItem) {
	root.ID = 0
	root.Field = "字段"

	item := model.ZdRangesItem{ParentID: root.ID, Field: "配置"}
	item.ID = uint(configId)
	root.Fields = append(root.Fields, &item)

	return
}
