package serverUtils

import (
	"encoding/json"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"io/ioutil"
	"regexp"
	"strings"
)

func ConvertDef(data interface{}) (def model.ZdDef) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &def)

	return
}
func ConvertField(data interface{}) (field model.ZdField) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &field)

	return
}

func ConvertSection(data interface{}) (section model.ZdSection) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &section)

	return
}

func ConvertRefer(data interface{}) (refer model.ZdRefer) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &refer)

	return
}

func ConvertRanges(data interface{}) (ranges model.ZdRanges) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}
func ConvertRangesItem(data interface{}) (item model.ZdRangesItem) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &item)

	return
}
func ConvertInstances(data interface{}) (ranges model.ZdInstances) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}
func ConvertInstancesItem(data interface{}) (item model.ZdInstancesItem) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &item)

	return
}
func ConvertExcel(data interface{}) (ranges model.ZdExcel) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}
func ConvertConfig(data interface{}) (ranges model.ZdConfig) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}
func ConvertText(data interface{}) (ranges model.ZdText) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &ranges)

	return
}

func ConvertParams(data interface{}) (mp map[string]string) {
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &mp)

	return
}

func GetRelativePath(pth string) string {
	idx := strings.LastIndex(pth, consts.PthSep)
	folder := pth[:idx+1]
	folder = strings.Replace(folder, vari.ZdDir, "", 1)

	return folder
}

func DealWithPathSepRight(pth string) string {
	pth = fileUtils.RemovePathSepLeftIfNeeded(pth)
	pth = fileUtils.AddPathSepRightIfNeeded(pth)

	return pth
}

func AddExt(pth, ext string) string {
	if strings.LastIndex(pth, ext) != len(pth)-4 {
		pth += ext
	}

	return pth
}

func GetDirs(dir string, dirs *[]domain.Dir) {
	isArticleFiles, _ := regexp.MatchString("yaml.article", dir)
	if isArticleFiles {
		return
	}

	folder := fileUtils.AddSepIfNeeded(dir)
	*dirs = append(*dirs, domain.Dir{Name: folder})

	files, _ := ioutil.ReadDir(folder)
	for _, fi := range files {
		name := fi.Name()
		if !fi.IsDir() || strings.Index(name, ".") == 0 {
			continue
		}

		childFolder := folder + name
		GetDirs(childFolder, dirs)
	}

	return
}
