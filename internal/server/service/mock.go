package serverService

import (
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/internal/pkg/service"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type MockService struct {
	ResService    *ResService            `inject:""`
	MainService   *service.MainService   `inject:""`
	OutputService *service.OutputService `inject:""`
}

func (s *MockService) Init() (err error) {
	vari.GlobalVars.MockData = &model.MockData{}
	vari.GlobalVars.MockData.Paths = map[string]map[string]*model.EndPoint{}
	var files []string

	s.LoadDef(vari.GlobalVars.MockDir, &files, 0)

	for _, file := range files {
		data := model.MockData{}

		content := fileUtils.ReadFileBuf(file)
		err := yaml.Unmarshal(content, &data)

		if err != nil {
			continue
		}

		for pth, mp := range data.Paths {
			pattern := s.getPathPatten(pth)
			vari.GlobalVars.MockData.Paths[pattern] = mp
		}
	}

	return
}

func (s *MockService) LoadDef(pth string, files *[]string, level int) (err error) {
	if !fileUtils.IsDir(pth) {
		*files = append(*files, pth)
		return
	}

	dir, err := os.ReadDir(pth)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		childPath := filepath.Join(pth, fi.Name())

		if fi.IsDir() && level < 3 {
			s.LoadDef(childPath, files, level+1)

		} else {
			*files = append(*files, childPath)

		}
	}

	return nil
}

func (s *MockService) GetResp(reqPath, reqMethod string) (ret interface{}, err error) {
	reqPath = s.addPrefixIfNeeded(reqPath)
	reqMethod = strings.ToLower(reqMethod)

	for pth, mp := range vari.GlobalVars.MockData.Paths {
		if !regexp.MustCompile(pth).MatchString(reqPath) { // no match
			continue
		}

		if mp[reqMethod] == nil { // no such a method
			continue
		}

		ret, _ = s.GenData(mp[reqMethod])
	}

	return
}

func (s *MockService) GenData(endpoint *model.EndPoint) (ret interface{}, err error) {
	vari.GlobalVars.RunMode = consts.RunModeServerRequest
	vari.GlobalVars.Total = endpoint.Lines
	vari.GlobalVars.OutputFormat = "json"
	vari.GlobalVars.ExportFields = strings.Split(endpoint.Fields, ",")

	dataType := endpoint.Type
	if dataType == "item" {
		vari.GlobalVars.Total = 1
	}

	configFile := filepath.Join(vari.ZdDir, endpoint.Config)
	vari.GlobalVars.ConfigFileDir = fileUtils.GetAbsDir(configFile)

	configContent := fileUtils.ReadFileBuf(configFile)

	contents := [][]byte{configContent}

	s.MainService.GenerateDataByContents(contents)

	records := s.OutputService.GenRecords()
	if dataType == "item" {
		ret = records[0]
	} else {
		ret = records
	}

	return
}

func (s *MockService) addPrefixIfNeeded(pth string) (ret string) {
	ret = "/" + strings.TrimPrefix(pth, "/")
	return
}

func (s *MockService) getPathPatten(pth string) (ret string) {
	regx := regexp.MustCompile(`({[^}]+?})`)
	ret = regx.ReplaceAllString(pth, "(.+)")

	ret = "^" + ret + "$"

	return
}
