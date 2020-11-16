package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/easysoft/zendata/src/model"
	constant "github.com/easysoft/zendata/src/utils/const"
	fileUtils "github.com/easysoft/zendata/src/utils/file"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	logUtils "github.com/easysoft/zendata/src/utils/log"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/mattn/go-runewidth"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"path"
	"sort"
	"strings"
)

func ListRes () {
	res, nameWidth, titleWidth := LoadRes("")
	PrintRes(res, nameWidth, titleWidth)
}

func LoadRes(resType string) (res map[string][]model.ResFile, nameWidth, titleWidth int) {
	res = map[string][]model.ResFile{}

	for _, key := range constant.ResKeys {
		GetFilesAndDirs(key, key, &res)
	}

	for _, key := range constant.ResKeys {
		arr := make([]model.ResFile, 0)

		for _, item := range res[key] {
			pth := item.Path
			name := PathToName(pth, key)
			var title, desc, tp string

			if key == constant.ResDirData { // data dir contains excel
				title, desc, tp = ReadExcelInfo(pth)
			} else if key == constant.ResDirYaml || key == constant.ResDirUsers {
				fileExt := path.Ext(pth)
				if fileExt == ".txt" {
					title, desc, tp = ReadTextInfo(pth, key)
				} else {
					title, desc, tp = ReadYamlInfo(pth)
				}
			}

			item.Name = name
			item.Title = title
			item.Desc = desc
			item.ResType = tp

			lent := runewidth.StringWidth(name)
			if lent > nameWidth {
				nameWidth = lent
			}

			if key == constant.ResDirData {
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

			if resType == "" || resType == item.ResType {
				arr = append(arr, item)
			}
		}

		res[key] = SortByName(arr)
	}

	return
}

func PrintRes(res map[string][]model.ResFile, nameWidth, titleWidth int) {
	dataMsg := ""
	yamlMsg := ""
	usersMsg := ""
	idx := 0
	for _, key := range constant.ResKeys {
		arr := res[key]

		for _, item := range arr {
			name := item.Name
			desc := item.Desc
			title := item.Title
			titles := strings.Split(title, "|")

			idx2 := 0
			for _, title := range titles {
				if idx2 > 0 {
					name = ""
				}
				name = name + strings.Repeat(" ", nameWidth - runewidth.StringWidth(name))

				title = title  + strings.Repeat(" ", titleWidth - runewidth.StringWidth(title))
				msg := fmt.Sprintf("%s  %s  %s\n", name, title, desc)

				if key == constant.ResDirData {
					dataMsg = dataMsg + msg
				} else if key == constant.ResDirYaml {
					yamlMsg = yamlMsg + msg
				} else if key == constant.ResDirUsers {
					usersMsg = usersMsg + msg
				}

				idx2++
			}

			idx++
		}
	}

	logUtils.PrintTo(dataMsg + "\n" + yamlMsg + "\n" + usersMsg)
}

func GetFilesAndDirs(pth, typ string, res *map[string][]model.ResFile)  {
	if !fileUtils.IsAbosutePath(pth) {
		pth = vari.WorkDir + pth
	}

	dir, err := ioutil.ReadDir(pth)
	if err != nil {
		return
	}

	for _, fi := range dir {
		if fi.IsDir() {
			GetFilesAndDirs(pth + constant.PthSep + fi.Name(), typ, res)
		} else {
			name := fi.Name()
			fileExt := path.Ext(name)
			if fileExt != ".yaml" && fileExt != ".xlsx" && fileExt != ".txt" {
				continue
			}

			file := model.ResFile{Path: pth + constant.PthSep + name}
			(*res)[typ] = append((*res)[typ], file)
		}
	}
}

func ReadYamlInfo(path string) (title, desc, resType string) {
	info := model.DefInfo{}

	if strings.Index(path, "zentao/number/") > -1 {
		log.Println(path)
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
	resType = constant.ResTypeExcel
	return
}

func ReadTextInfo(path, key string) (title, desc, resType string) {
	title = PathToName(path, key)
	desc = i118Utils.I118Prt.Sprintf("text_data")
	resType = constant.ResTypeText
	return
}

func PathToName(path, key string) string {
	//if strings.Index(path, "users.txt") > -1 {
	//	logUtils.PrintToScreen("111")
	//}

	name := strings.ReplaceAll(path, constant.PthSep,".")
	sep := "." + key + "."
	name = name[strings.Index(name, sep)+len(sep):]
	if key == constant.ResDirData { // remove .xlsx postfix
		name = name[:strings.LastIndex(name, ".")]
	}

	return name
}

func SortByName(arr []model.ResFile) []model.ResFile {
	sort.Slice(arr, func(i, j int) bool {
		flag := false
		if arr[i].Name > (arr[j].Name) {
			flag = true
		}
		return flag
	})
	return arr
}

func GetYamlResType(def model.DefInfo) string {
	if def.Range != "" {
		return constant.ResTypeConfig
	} else if def.Ranges != nil {
		return constant.ResTypeRanges
	} else if def.Instances != nil {
		return constant.ResTypeInstances
	}

	return ""
}
