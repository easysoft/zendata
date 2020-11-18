package serverService

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/server/repo"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
)

type ConfigService struct {
	configRepo *serverRepo.ConfigRepo
	resService *ResService
}

func (s *ConfigService) List() (list []*model.ZdConfig) {
	config := s.resService.LoadRes("config")
	list, _ = s.configRepo.List()

	s.saveResToDB(config, list)
	list, _ = s.configRepo.List()

	return
}

func (s *ConfigService) Get(id int) (config model.ZdConfig) {
	config, _ = s.configRepo.Get(uint(id))

	return
}

func (s *ConfigService) Save(config *model.ZdConfig) (err error) {
	err = s.configRepo.Save(config)

	return
}

func (s *ConfigService) Remove(id int) (err error) {
	err = s.configRepo.Remove(uint(id))
	if err != nil {
		return
	}

	config, _ := s.configRepo.Get(uint(id))
	logUtils.PrintTo(config.Path)
	//fileUtils.RemoveExist(config.Path)

	return
}

func (s *ConfigService) saveResToDB(config []model.ResFile, list []*model.ZdConfig) (err error) {
	names := make([]string, 0)
	for _, item := range list {
		names = append(names, item.Path)
	}

	for _, item := range config {
		if !stringUtils.FindInArrBool(item.Path, names) {
			config := model.ZdConfig{Title: item.Title, Name: item.Name, Desc: item.Desc, Field: item.Title, Path: item.Path, Note: item.Desc}
			s.configRepo.Save(&config)
		}
	}

	return
}

func NewConfigService(configRepo *serverRepo.ConfigRepo) *ConfigService {
	return &ConfigService{configRepo: configRepo}
}
