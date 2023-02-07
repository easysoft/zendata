package serverService

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	"github.com/easysoft/zendata/internal/pkg/service"
	serverRepo "github.com/easysoft/zendata/internal/server/repo"
	dateUtils "github.com/easysoft/zendata/pkg/utils/date"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/helper/dir"
	"gopkg.in/yaml.v2"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type MockService struct {
	MainService   *service.MainService   `inject:""`
	OutputService *service.OutputService `inject:""`
	MockService   *service.MockService   `inject:""`
	MockRepo      *serverRepo.MockRepo   `inject:""`
}

func (s *MockService) List(keywords string, page int) (pos []*model.ZdMock, total int, err error) {
	pos, total, err = s.MockRepo.List(strings.TrimSpace(keywords), page)
	return
}

func (s *MockService) Get(id int) (po model.ZdMock, err error) {
	po, err = s.MockRepo.Get(uint(id))

	return
}

func (s *MockService) Init() (err error) {
	vari.GlobalVars.MockData = &model.MockData{}
	vari.GlobalVars.MockData.Paths = map[string]map[string]map[string]map[string]*model.EndPoint{}
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

func (s *MockService) GetResp(reqPath, reqMethod, respCode, mediaType string) (ret interface{}, err error) {
	reqPath = s.addPrefixIfNeeded(reqPath)
	reqMethod = strings.ToLower(reqMethod)

	for pth, mp := range vari.GlobalVars.MockData.Paths {
		if !regexp.MustCompile(pth).MatchString(reqPath) { // no match
			continue
		}

		if mp[reqMethod] == nil { // no such a method
			continue
		}

		ret, _ = s.GenData(mp[reqMethod][respCode][mediaType])
	}

	return
}

func (s *MockService) GenData(endpoint *model.EndPoint) (ret interface{}, err error) {
	vari.GlobalVars.RunMode = consts.RunModeServerRequest
	vari.GlobalVars.Total = endpoint.Lines
	vari.GlobalVars.OutputFormat = "json"
	vari.GlobalVars.ExportFields = strings.Split(endpoint.Fields, ",")

	dataType := endpoint.Type
	if dataType != consts.SchemaTypeArray {
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

func (s *MockService) Upload(ctx iris.Context, fh *multipart.FileHeader) (content, mockConf, dataConf, pth string, err error) {
	filename, err := fileUtils.GetUploadFileName(fh.Filename)
	if err != nil {
		logUtils.PrintTo(fmt.Sprintf("获取文件名失败，错误%s", err.Error()))
		return
	}

	targetDir := filepath.Join("upload", dateUtils.DateStr(time.Now()))
	absDir := filepath.Join(dir.GetCurrentAbPath(), targetDir)

	err = dir.InsureDir(targetDir)
	if err != nil {
		logUtils.PrintTo(fmt.Sprintf("文件上传失败，错误%s", err.Error()))
		return
	}

	pth = filepath.Join(absDir, filename)
	_, err = ctx.SaveFormFile(fh, pth)
	if err != nil {
		logUtils.PrintTo(fmt.Sprintf("文件上传失败，错误%s", err.Error()))
		return
	}

	content = fileUtils.ReadFile(pth)

	vari.GlobalVars.Output = fileUtils.GetFileOrFolderDir(pth)
	mockPath, dataPath, err := s.MockService.GenMockDef(pth)

	if err == nil {
		mockConf = fileUtils.ReadFile(mockPath)
		dataConf = fileUtils.ReadFile(dataPath)
	}

	return
}
