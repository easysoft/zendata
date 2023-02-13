package serverService

import (
	"fmt"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
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
	DefService    *DefService            `inject:""`
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

func (s *MockService) Save(po *model.ZdMock) (err error) {
	err = s.MockRepo.Save(po)

	return
}

func (s *MockService) Remove(id int) (err error) {
	err = s.MockRepo.Remove(uint(id))

	return
}

func (s *MockService) Init() (err error) {
	vari.GlobalVars.MockData = &domain.MockData{}
	vari.GlobalVars.MockData.Paths = map[string]map[string]map[string]map[string]*domain.EndPoint{}
	var files []string

	s.LoadDef(vari.GlobalVars.MockDir, &files, 0)

	for _, file := range files {
		data := domain.MockData{}

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

		ret, _ = s.GenDataForServerRequest(mp[reqMethod][respCode][mediaType])
	}

	return
}

func (s *MockService) GenDataForServerRequest(endpoint *domain.EndPoint) (ret interface{}, err error) {
	vari.GlobalVars.RunMode = consts.RunModeServerRequest
	vari.GlobalVars.Total = endpoint.Lines
	vari.GlobalVars.OutputFormat = "json"
	vari.GlobalVars.ExportFields = strings.Split(endpoint.Fields, ",")

	dataType := endpoint.Type
	if dataType != consts.SchemaTypeArray {
		vari.GlobalVars.Total = 1
	}

	// eval config file path
	configFile := filepath.Join(vari.GlobalVars.MockDir, endpoint.Config)
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

func (s *MockService) GenDataForMockPreview(endpoint *domain.EndPoint, dataConfig string) (ret interface{}, err error) {
	vari.GlobalVars.RunMode = consts.RunModeMockPreview
	vari.GlobalVars.Total = endpoint.Lines
	vari.GlobalVars.OutputFormat = "json"
	vari.GlobalVars.ExportFields = strings.Split(endpoint.Fields, ",")

	dataType := endpoint.Type
	if dataType != consts.SchemaTypeArray {
		vari.GlobalVars.Total = 1
	}

	contents := [][]byte{[]byte(dataConfig)}
	s.MainService.GenerateDataByContents(contents)

	records := s.OutputService.GenRecords()
	if dataType == consts.SchemaTypeObject {
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

func (s *MockService) Upload(ctx iris.Context, fh *multipart.FileHeader) (
	name, content, mockConf, dataConf, pth string, err error, id uint) {

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
	name, mockPath, dataPath, err := s.MockService.GenMockDef(pth)

	if err == nil {
		mockConf = fileUtils.ReadFile(mockPath)
		dataConf = fileUtils.ReadFile(dataPath)
		fi := domain.ResFile{FileName: fh.Filename, Path: dataPath}
		_, id = s.DefService.SyncToDB(fi)
	}

	return
}

func (s *MockService) GetPreviewData(id int) (data domain.MockData, err error) {
	po, err := s.MockRepo.Get(uint(id))

	yaml.Unmarshal([]byte(po.MockContent), &data)
	data.Id = id

	return
}

func (s *MockService) GetPreviewResp(req domain.MockPreviewReq) (ret interface{}, err error) {
	po, err := s.MockRepo.Get(uint(req.Id))

	data := domain.MockData{}
	err = yaml.Unmarshal([]byte(po.MockContent), &data)
	if err != nil {
		return
	}

	for pth, mp := range data.Paths {
		if req.Url == pth {
			endpoint := mp[req.Method][req.Code][req.Media]
			ret, _ = s.GenDataForMockPreview(endpoint, po.DataContent)
			return
		}
	}

	return
}
