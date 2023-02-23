package helper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	consts "github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/domain"
	"github.com/easysoft/zendata/internal/pkg/model"
	fileUtils "github.com/easysoft/zendata/pkg/utils/file"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v3"
)

func ListData() {
	res := map[string][]domain.ResFile{}
	GetFilesAndDirs(vari.ZdDir, consts.ResDirUsers, &res)

	res, nameWidth, titleWidth := LoadRes(res)
	PrintRes(res, nameWidth, titleWidth)
}

func ListRes() (ret map[string][]domain.ResFile) {
	res, nameWidth, titleWidth := GetRes()
	PrintRes(res, nameWidth, titleWidth)

	return
}

func GetRes() (ret map[string][]domain.ResFile, nameWidth, titleWidth int) {
	ret = map[string][]domain.ResFile{}

	for _, key := range consts.ResKeys {
		GetFilesAndDirs(key, key, &ret)
	}

	ret, nameWidth, titleWidth = LoadRes(ret)

	return
}

func LoadRes(res map[string][]domain.ResFile) (ret map[string][]domain.ResFile, nameWidth, titleWidth int) {
	ret = map[string][]domain.ResFile{}

	for _, key := range consts.ResKeys {
		arr := make([]domain.ResFile, 0)

		for _, item := range res[key] {
			pth := item.Path

			fileExt := filepath.Ext(pth)
			isArticleFiles := false
			var title, desc, tp string

			if key == consts.ResDirData { // data dir contains excel
				title, desc, tp = ReadExcelInfo(pth)
			} else if key == consts.ResDirYaml || key == consts.ResDirUsers {
				isArticleFiles, _ = regexp.MatchString("yaml.article", pth)

				if fileExt == ".txt" { // ignore packaged article text file
					title, desc, tp = ReadTextInfo(pth, key)
				} else {
					title, desc, tp = ReadYamlInfo(pth)
				}
			}

			item.ReferName = PathToName(pth, key, tp)
			item.Title = title
			item.Desc = desc
			item.ResType = tp

			lent := runewidth.StringWidth(item.ReferName)
			if lent > nameWidth {
				nameWidth = lent
			}

			if key == consts.ResDirData {
				sheets := strings.Split(title, "|")
				for _, sheet := range sheets {
					lent2 := runewidth.StringWidth(sheet)
					if lent2 > titleWidth {
						titleWidth = lent2
					}
				}
			} else {
				lent2 := runewidth.StringWidth(title)
				if lent2 > titleWidth {
					titleWidth = lent2
				}
			}

			if !isArticleFiles {
				arr = append(arr, item)
			}
		}

		ret[key] = SortByName(arr)
	}

	return
}

func PrintRes(res map[string][]domain.ResFile, nameWidth, titleWidth int) {
	dataMsg := ""
	yamlMsg := ""
	usersMsg := ""
	idx := 0
	for _, key := range consts.ResKeys {
		arr := res[key]

		for _, item := range arr {
			name := item.ReferName
			desc := item.Desc
			title := item.Title
			titles := strings.Split(title, "|")

			idx2 := 0
			for _, title := range titles {
				if idx2 > 0 {
					name = ""
				}
				name = name + strings.Repeat(" ", nameWidth-runewidth.StringWidth(name))

				title = title + strings.Repeat(" ", titleWidth-runewidth.StringWidth(title))
				msg := fmt.Sprintf("%s  %s  %s\n", name, title, desc)

				if key == consts.ResDirData {
					dataMsg = dataMsg + msg
				} else if key == consts.ResDirYaml {
					yamlMsg = yamlMsg + msg
				} else if key == consts.ResDirUsers {
					usersMsg = usersMsg + msg
				}

				idx2++
			}

			idx++
		}
	}

	logUtils.PrintTo(dataMsg + "\n" + yamlMsg + "\n" + usersMsg)
}

func GetFilesAndDirs(pth, typ string, res *map[string][]domain.ResFile) {
	if !fileUtils.IsAbsPath(pth) {
		pth = vari.ZdDir + pth
	}

	dir, err := ioutil.ReadDir(pth)
	if err != nil {
		return
	}

	for _, fi := range dir {
		if fi.IsDir() {
			GetFilesAndDirs(filepath.Join(pth, fi.Name()), typ, res)
		} else {
			name := fi.Name()
			fileExt := filepath.Ext(name)
			if fileExt != ".yaml" && fileExt != ".xlsx" && fileExt != ".txt" {
				continue
			}

			file := domain.ResFile{Path: filepath.Join(pth, name), UpdatedAt: fi.ModTime()}
			(*res)[typ] = append((*res)[typ], file)
		}
	}
}

func ReadYamlInfo(path string) (title, desc, resType string) {
	info := domain.DefInfo{}

	if strings.Index(path, "apache") > -1 {
		logUtils.PrintTo("")
	}

	yamlContent, err := ioutil.ReadFile(path)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}
	err = yaml.Unmarshal(yamlContent, &info)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_parse_file", path))
		return
	}

	title = info.Title
	desc = info.Desc
	resType = GetYamlResType(info)
	return
}

func ReadExcelInfo(path string) (title, desc, resType string) {
	excel, err := excelize.OpenFile(path)
	if err != nil {
		logUtils.PrintTo(i118Utils.I118Prt.Sprintf("fail_to_read_file", path))
		return
	}

	for index, sheet := range excel.GetSheetList() {
		if index > 0 {
			title = title + "|"
		}
		title = title + sheet
	}

	desc = i118Utils.I118Prt.Sprintf("excel_data")
	resType = consts.ResTypeExcel
	return
}

func ReadTextInfo(path, key string) (title, desc, resType string) {
	title = PathToName(path, key, consts.ResTypeText)
	desc = i118Utils.I118Prt.Sprintf("text_data")
	resType = consts.ResTypeText
	return
}

func PathToName(path, key, tp string) string {
	isWorkData := strings.Index(path, vari.ZdDir) > -1
	if isWorkData { // user data in workdir
		path = strings.Replace(path, vari.ZdDir, "", 1)
	}

	nameSep := consts.PthSep
	if tp != consts.ResTypeText && tp != consts.ResTypeYaml && tp != consts.ResTypeConfig {
		nameSep = "."
		path = strings.ReplaceAll(path, consts.PthSep, nameSep)
		path = path[strings.Index(path, nameSep)+len(nameSep):]
	}
	if isWorkData {
		return path
	}

	sep := nameSep + key + nameSep
	name := path[strings.Index(path, sep)+len(sep):]
	if key == consts.ResDirData { // remove .xlsx postfix for excel data
		name = name[:strings.LastIndex(name, nameSep)]
	}

	return name
}
func removeDirPrefix(name, seq string) (ret string) {

	return
}

func SortByName(arr []domain.ResFile) []domain.ResFile {
	sort.Slice(arr, func(i, j int) bool {
		flag := false
		if arr[i].ReferName > (arr[j].ReferName) {
			flag = true
		}
		return flag
	})
	return arr
}

func GetYamlResType(def domain.DefInfo) string {
	if def.Ranges != nil {
		return consts.ResTypeRanges
	} else if def.Instances != nil {
		return consts.ResTypeInstances
	} else if def.Fields != nil {
		return consts.ResTypeYaml
	} else {
		return consts.ResTypeConfig
	}

	return ""
}

func GetDefFromYamlFile(path string) (po *model.ZdDef, content string, err error) {
	po = &model.ZdDef{}
	contentByte, _ := os.ReadFile(path)
	content = string(contentByte)

	yamlContent := ReplaceSpecialChars(contentByte)
	err = yaml.Unmarshal(yamlContent, po)

	return
}

func GetDefFromYamlString(content string) (po *model.ZdDef, err error) {
	po = &model.ZdDef{}
	yamlContent := ReplaceSpecialChars([]byte(content))
	err = yaml.Unmarshal(yamlContent, po)

	return
}
