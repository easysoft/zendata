package serverService

import (
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type MockService struct {
	ResService *ResService `inject:""`
}

func (s *MockService) Init() (err error) {
	vari.GlobalVars.MockData = &model.MockData{}
	vari.GlobalVars.MockData.Paths = map[string]model.Path{}
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
