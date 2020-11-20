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

type ConfigService struct {
	configRepo *serverRepo.ConfigRepo
	resService *ResService
}

func (s *ConfigService) List(keywords string, page int) (list []*model.ZdConfig, total int) {
	config := s.resService.LoadRes("config")
	list, _, _ = s.configRepo.List("", -1)

	s.importResToDB(config, list)
	list, total, _ = s.configRepo.List(strings.TrimSpace(keywords), page)

	return
}

func (s *ConfigService) Get(id int) (config model.ZdConfig, dirTree model.Dir) {
	if id > 0 {
		config, _ = s.configRepo.Get(uint(id))
	}

	dirTree = model.Dir{Name: fileUtils.AddSepIfNeeded(constant.ResDirYaml)}
	serverUtils.GetDirTree(&dirTree)

	return
}

func (s *ConfigService) Save(config *model.ZdConfig) (err error) {
	config.Folder = serverUtils.DealWithPathSepRight(config.Folder)
	config.Path = vari.WorkDir + config.Folder + serverUtils.AddExt(config.Title, ".yaml")
	config.Name = service.PathToName(config.Path, constant.ResDirYaml)

	if config.ID == 0 {
		err = s.Create(config)
	} else {
		err = s.Update(config)
	}

	return
}

func (s *ConfigService) Create(config *model.ZdConfig) (err error) {
	s.dataToYaml(config)
	err = s.configRepo.Create(config)

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

	s.dataToYaml(config)
	err = s.configRepo.Update(config)

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

func (s *ConfigService) dataToYaml(config *model.ZdConfig) (str string) {

	return
}

func (s *ConfigService) importResToDB(config []model.ResFile, list []*model.ZdConfig) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range config {
		if !stringUtils.FindInArrBool(item.Path, names) {
			content, _ := ioutil.ReadFile(item.Path)
			yamlContent := stringUtils.ReplaceSpecialChars(content)
			config := model.ZdConfig{}
			err = yaml.Unmarshal(yamlContent, &config)
			config.Title = item.Title
			config.Name = item.Name
			config.Desc = item.Desc
			config.Path = item.Path
			config.Folder = serverUtils.GetRelativePath(config.Path)
			config.Field = item.Title
			config.Note = item.Desc
			config.Yaml = string(content)

			s.configRepo.Create(&config)
		}
	}

	return
}

func NewConfigService(configRepo *serverRepo.ConfigRepo) *ConfigService {
	return &ConfigService{configRepo: configRepo}
}
