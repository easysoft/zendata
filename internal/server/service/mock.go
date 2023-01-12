package serverService

import (
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

type MockService struct {
	ResService *ResService `inject:""`
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

		for k, v := range data.Paths {
			vari.GlobalVars.MockData.Paths[k] = v
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
	reqPath = s.AddPrefixIfNeeded(reqPath)
	reqMethod = strings.ToLower(reqMethod)

	if vari.GlobalVars.MockData.Paths[reqPath] == nil || // no such a path
		vari.GlobalVars.MockData.Paths[reqPath][reqMethod] == nil { // no such a method
		return
	}

	ret, _ = s.GenData(vari.GlobalVars.MockData.Paths[reqPath][reqMethod])

	return
}

func (s *MockService) GenData(endpoint *model.EndPoint) (ret interface{}, err error) {
	ret = endpoint

	return
}

func (s *MockService) AddPrefixIfNeeded(pth string) (ret string) {
	ret = "/" + strings.TrimPrefix(pth, "/")
	return
}
