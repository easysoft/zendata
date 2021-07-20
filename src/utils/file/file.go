package fileUtils

import (
	"fmt"
	"github.com/easysoft/zendata/res"
	"github.com/easysoft/zendata/src/model"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	i118Utils "github.com/easysoft/zendata/src/utils/i118"
	"github.com/easysoft/zendata/src/utils/vari"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ReadFile(filePath string) string {
	if !FileExist(filePath) {
		return ""
	}
	buf := ReadFileBuf(filePath)
	str := string(buf)
	str = commonUtils.RemoveBlankLine(str)
	return str
}

func ReadFileBuf(filePath string) []byte {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte(err.Error())
	}

	return buf
}

func WriteFile(filePath string, content string) {
	dir := filepath.Dir(filePath)
	MkDirIfNeeded(dir)

	var d1 = []byte(content)
	err2 := ioutil.WriteFile(filePath, d1, 0666) //写入文件(字节数组)
	check(err2)
}

func RemoveExist(path string) {
	os.Remove(path)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func MkDirIfNeeded(dir string) error {
	if !FileExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		return err
	}

	return nil
}
func RmFile(dir string) error {
	if FileExist(dir) {
		err := os.RemoveAll(dir)
		return err
	}

	return nil
}

func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func GetAbsolutePath(pth string) string {
	if !IsAbsPath(pth) {
		pth, _ = filepath.Abs(pth)
	}

	return pth
}
func GetAbsoluteDir(pth string) string {
	pth = GetAbsolutePath(pth)
	pth = filepath.Dir(pth)
	pth = AddSepIfNeeded(pth)

	return pth
}

func GetRelatPath(pth string) string {
	pth = strings.TrimPrefix(pth, vari.ZdPath)

	return pth
}

func IsAbsPath(pth string) bool {
	return filepath.IsAbs(pth) ||
		strings.Index(pth, ":") == 1 // windows
}

func AddSepIfNeeded(pth string) string {
	sepa := string(os.PathSeparator)

	if strings.LastIndex(pth, sepa) < len(pth)-1 {
		pth += sepa
	}
	return pth
}

func ReadResData(path string) string {
	isRelease := commonUtils.IsRelease()

	var jsonStr string
	if isRelease {
		data, _ := res.Asset(path)
		jsonStr = string(data)
	} else {
		jsonStr = ReadFile(path)
	}

	return jsonStr
}

func GetExeDir() string { // where zd.exe file in
	var dir string
	arg1 := strings.ToLower(os.Args[0])

	name := filepath.Base(arg1)
	if strings.Index(name, "zd") == 0 && strings.Index(arg1, "go-build") < 0 {
		p, _ := exec.LookPath(os.Args[0])
		if strings.Index(p, string(os.PathSeparator)) > -1 {
			dir = p[:strings.LastIndex(p, string(os.PathSeparator))]
		}
	} else { // debug
		dir, _ = os.Getwd()
		//if commonUtils.IsMac() {
		//	dir = "/Users/aaron/rd/project/zentao/go/zd"
		//}
	}

	dir, _ = filepath.Abs(dir)
	dir = AddSepIfNeeded(dir)

	//fmt.Printf("Debug: Launch %s in %s \n", arg1, dir)
	return dir
}

func GetWorkDir() string { // where we run file in
	dir, _ := os.Getwd()

	dir, _ = filepath.Abs(dir)
	dir = AddSepIfNeeded(dir)

	//fmt.Printf("Debug: Launch %s in %s \n", arg1, dir)
	return dir
}

func GetAbsDir(path string) string {
	abs, _ := filepath.Abs(filepath.Dir(path))
	abs = AddSepIfNeeded(abs)
	return abs
}

func GetResProp(from, currFileDir string) (resFile, resType, sheet string) { // from resource
	if strings.LastIndex(from, ".yaml") > -1 { // yaml, ip.v1.yaml
		resFile = ConvertResYamlPath(from, currFileDir)
		resType = "yaml"
	} else if strings.LastIndex(from, ".txt") > -1 {
		resFile = ConvertResYamlPath(from, currFileDir)
		resType = "text"
	} else { // excel, like address.cn.v1.china
		resFile, sheet = ConvertResExcelPath(from, currFileDir)
		resType = "excel"
	}

	if resFile == "" {
		resPath := vari.ConfigFileDir + from
		if !FileExist(resPath) { // in same folder with passed config file, like dir/name.yaml
			resPath = vari.ZdPath + from
			if !FileExist(resPath) { // in res file
				resPath = ""
			}
		}
		resFile = resPath
	}

	if resFile == "" {
		color.New(color.FgCyan).Fprintf(color.Output, i118Utils.I118Prt.Sprintf("fail_to_find_res", from)+"\n")
	}

	return
}

func ConvertReferRangeToPath(file, currFile string) (path string) {
	dir := GetAbsDir(currFile)
	path = ConvertResYamlPath(file, dir)

	if path == "" {
		resPath := GetAbsDir(currFile) + file
		if !FileExist(resPath) { // in same folder
			resPath = vari.ZdPath + file
			if !FileExist(resPath) { // in res file
				resPath = ""
			}
		}
		path = resPath
	}

	return
}

func ConvertResYamlPath(from, workDir string) (ret string) {
	pth := namedFileExistInDir(from, workDir)
	if pth != "" {
		ret = pth
		return
	}

	arr := strings.Split(from, ".")
	for i := 0; i < len(arr); i++ {
		dir := ""
		if i > 0 {
			dir = strings.Join(arr[:i], constant.PthSep)
		}
		file := strings.Join(arr[i:], ".")

		relatPath := ""
		if dir != "" {
			relatPath = dir + constant.PthSep + file
		} else {
			relatPath = file
		}

		realPth0 := filepath.Join(workDir, relatPath)
		realPth1 := vari.ZdPath + constant.ResDirYaml + constant.PthSep + relatPath
		realPth2 := vari.ZdPath + constant.ResDirUsers + constant.PthSep + relatPath
		if FileExist(realPth0) {
			ret = realPth0
			break
		} else if FileExist(realPth1) {
			ret = realPth1
			break
		} else if FileExist(realPth2) {
			ret = realPth2
			break
		}
	}

	return
}

func ConvertResExcelPath(from, dir string) (ret, sheet string) {
	pth := namedFileExistInDir(from, dir)
	if pth != "" {
		ret = pth
		return
	}

	path1 := from // address.cn.v1
	index := strings.LastIndex(from, ".")
	path2 := from[:index] // address.cn.v1.china

	paths := [2]string{path1, path2}
	for index, filePath := range paths {

		arr := strings.Split(filePath, ".")
		for i := 0; i < len(arr); i++ {
			dir := ""
			if i > 0 {
				dir = strings.Join(arr[:i], constant.PthSep)
			}

			tagFile := strings.Join(arr[i:], ".") + ".xlsx"

			relatPath := ""
			if dir != "" {
				relatPath = dir + constant.PthSep + tagFile
			} else {
				relatPath = tagFile
			}

			realPth := vari.ZdPath + constant.ResDirData + constant.PthSep + relatPath
			if FileExist(realPth) {
				if index == 1 {
					sheet = from[strings.LastIndex(from, ".")+1:]
				}
				ret = realPth
				return
			}
		}
	}

	if ret == "" { // try excel dir
		realPth := vari.ZdPath + constant.ResDirData + constant.PthSep +
			strings.Replace(from, ".", constant.PthSep, -1)
		if IsDir(realPth) {
			ret = realPth
			return
		}
	}

	if ret == "" {
		color.New(color.FgCyan).Fprintf(color.Output, i118Utils.I118Prt.Sprintf("fail_to_find_res", from)+"\n")
	}

	return
}

func ComputerReferFilePath(file string, field *model.DefField) (resPath string) {
	resPath = file
	if IsAbsPath(resPath) && FileExist(resPath) {
		return
	}

	resPath = field.FileDir + file
	if FileExist(resPath) {
		return
	}

	resPath = vari.ConfigFileDir + file
	if FileExist(resPath) {
		return
	}

	resPath = vari.DefaultFileDir + file
	if FileExist(resPath) {
		return
	}

	resPath = vari.ZdPath + constant.ResDirUsers + constant.PthSep + file
	if FileExist(resPath) {
		return
	}
	resPath = vari.ZdPath + constant.ResDirYaml + constant.PthSep + file
	if FileExist(resPath) {
		return
	}

	resPath = vari.ZdPath + file
	if FileExist(resPath) {
		return
	}

	return
}

func GetFilesByExtInDir(folder, ext string, files *[]string) {
	folder, _ = filepath.Abs(folder)

	if !IsDir(folder) {
		if ext == "" || filepath.Ext(folder) == ext {
			*files = append(*files, folder)
		}

		return
	}

	dir, err := ioutil.ReadDir(folder)
	if err != nil {
		return
	}

	for _, fi := range dir {
		name := fi.Name()
		if commonUtils.IngoreFile(name) {
			continue
		}

		filePath := AddSepIfNeeded(folder) + name
		if fi.IsDir() {
			GetFilesByExtInDir(filePath, ext, files)
		} else if strings.Index(name, "~") != 0 && (ext == "" || filepath.Ext(filePath) == ext) {
			*files = append(*files, filePath)
		}
	}
}

func GetFileName(filePath string) string {
	fileName := filepath.Base(filePath)
	fileName = strings.TrimSuffix(fileName, filepath.Ext(filePath))

	return fileName
}

func GetFilesInDir(folder, ext string, files *[]string) {
	folder, _ = filepath.Abs(folder)

	if !IsDir(folder) {
		if filepath.Ext(folder) == ext {
			*files = append(*files, folder)
		}

		return
	}

	dir, err := ioutil.ReadDir(folder)
	if err != nil {
		return
	}

	for _, fi := range dir {
		name := fi.Name()
		if commonUtils.IngoreFile(name) {
			continue
		}

		filePath := AddSepIfNeeded(folder) + name
		if fi.IsDir() {
			GetFilesInDir(filePath, ext, files)
		} else if strings.Index(name, "~") != 0 && filepath.Ext(filePath) == ext {
			*files = append(*files, filePath)
		}
	}
}

func ChangeFileExt(filePath, ext string) string {
	ret := strings.TrimSuffix(filePath, filepath.Ext(filePath))
	ret += ext

	return ret
}

func AddPathSepRightIfNeeded(pth string) string {
	if pth[len(pth)-1:] != constant.PthSep {
		pth += constant.PthSep
	}

	return pth
}
func RemovePathSepLeftIfNeeded(pth string) string {
	if strings.Index(pth, constant.PthSep) == 0 {
		pth = pth[1:]
	}

	return pth
}

func namedFileExistInDir(file, dir string) (pth string) {
	if IsAbsPath(file) { // abs path, return it
		if FileExist(file) {
			pth = file
			return
		} else {
			return
		}
	} else {
		file = filepath.Join(dir, file)
		if FileExist(file) {
			pth = file
			return
		}
	}

	return
}

func GenArticleFiles(pth string, index int) (ret string) {
	pfix := fmt.Sprintf("%03d", index+1)

	ret = strings.TrimSuffix(pth, filepath.Ext(pth))
	ret += "-" + pfix + filepath.Ext(pth)

	return
}
