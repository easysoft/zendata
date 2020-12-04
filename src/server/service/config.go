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
	"regexp"
	"strings"
)

type ConfigService struct {
	configRepo  *serverRepo.ConfigRepo
	resService  *ResService
	sectionRepo *serverRepo.SectionRepo
}

func (s *ConfigService) List(keywords string, page int) (list []*model.ZdConfig, total int) {
	list, total, _ = s.configRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *ConfigService) Get(id int) (config model.ZdConfig, dirs []model.Dir) {
	if id > 0 {
		config, _ = s.configRepo.Get(uint(id))
	}

	serverUtils.GetDirs(constant.ResDirYaml, &dirs)

	return
}

func (s *ConfigService) Save(config *model.ZdConfig) (err error) {
	config.Folder = serverUtils.DealWithPathSepRight(config.Folder)
	config.Path = vari.WorkDir + config.Folder + serverUtils.AddExt(config.FileName, ".yaml")
	config.ReferName = service.PathToName(config.Path, constant.ResDirYaml, constant.ResTypeConfig)

	if config.ID == 0 {
		err = s.Create(config)
	} else {
		err = s.Update(config)
	}

	return
}

func (s *ConfigService) Create(config *model.ZdConfig) (err error) {
	err = s.configRepo.Create(config)
	s.updateYaml(config.ID)

	return
}

func (s *ConfigService) Update(config *model.ZdConfig) (err error) {
	var old model.ZdConfig
	old, err = s.configRepo.Get(config.ID)
	if err == gorm.ErrRecordNotFound {
		return
	}
	if config.Path != old.Path {
		fileUtils.RemoveExist(old.Path)
	}

	err = s.configRepo.Update(config)
	s.updateYaml(config.ID)

	return
}

func (s *ConfigService) Remove(id int) (err error) {
	var old model.ZdConfig
	old, err = s.configRepo.Get(uint(id))
	if err == gorm.ErrRecordNotFound {
		return
	}

	fileUtils.RemoveExist(old.Path)
	err = s.configRepo.Remove(uint(id))

	return
}

func (s *ConfigService) updateYaml(id uint) (err error) {
	var po model.ZdConfig
	po, _ = s.configRepo.Get(id)

	s.genYaml(&po)
	err = s.configRepo.UpdateYaml(po)
	fileUtils.WriteFile(po.Path, po.Yaml)

	return
}
func (s *ConfigService) genYaml(config *model.ZdConfig) (str string) {
	yamlObj := model.ResConfig{}
	s.configRepo.GenConfigRes(*config, &yamlObj)

	bytes, _ := yaml.Marshal(yamlObj)
	config.Yaml = stringUtils.ConvertYamlStringToMapFormat(bytes)

	return
}

func (s *ConfigService) Sync(files []model.ResFile) (err error) {
	list := s.configRepo.ListAll()

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
			s.configRepo.Remove(mp[fi.Path].ID)
			s.SyncToDB(fi)
		} else { // db is new

		}
	}

	return
}
func (s *ConfigService) SyncToDB(fi model.ResFile) (err error) {
	content, _ := ioutil.ReadFile(fi.Path)
	yamlContent := stringUtils.ReplaceSpecialChars(content)

	po := model.ZdConfig{}
	err = yaml.Unmarshal(yamlContent, &po)

	po.Title = fi.Title
	po.Desc = fi.Desc
	po.Path = fi.Path
	po.Folder = serverUtils.GetRelativePath(po.Path)
	if strings.Index(po.Path, vari.WorkDir+constant.ResDirYaml) > -1 {
		po.ReferName = service.PathToName(po.Path, constant.ResDirYaml, constant.ResTypeConfig)
	} else {
		po.ReferName = service.PathToName(po.Path, constant.ResDirUsers, constant.ResTypeConfig)
	}
	po.FileName = fileUtils.GetFileName(po.Path)
	po.Yaml = string(content)

	reg := regexp.MustCompile(`\t`)
	if reg.MatchString(po.Prefix) {
		po.Prefix = strings.ReplaceAll(po.Prefix, "\t", `\t`)
		po.Prefix = `"` + po.Prefix + `"`
	}

	s.configRepo.Create(&po)

	rangeSections := gen.ParseRangeProperty(po.Range)
	for i, rangeSection := range rangeSections {
		s.sectionRepo.SaveFieldSectionToDB(rangeSection, i, po.ID, "config")
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

func NewConfigService(configRepo *serverRepo.ConfigRepo, sectionRepo *serverRepo.SectionRepo) *ConfigService {
	return &ConfigService{configRepo: configRepo, sectionRepo: sectionRepo}
}
