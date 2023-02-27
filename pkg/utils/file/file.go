package fileUtils

import (
	"errors"
	"fmt"
	constant "github.com/easysoft/zendata/internal/pkg/const"
	commonUtils "github.com/easysoft/zendata/pkg/utils/common"
	i118Utils "github.com/easysoft/zendata/pkg/utils/i118"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/fatih/color"
	"github.com/oklog/ulid/v2"
	"github.com/snowlyg/helper/str"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
	err := os.WriteFile(filePath, d1, 0666) //写入文件(字节数组)
	check(err)
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

func GetExeDir() string { // where zd.exe file in
	var dir string
	arg1 := strings.ToLower(os.Args[0])

	name := filepath.Base(arg1)
	if strings.Index(name, "zd") == 0 && strings.Index(arg1, "go-build") < 0 { // release
		p, _ := exec.LookPath(os.Args[0])
		if strings.Index(p, string(os.PathSeparator)) > -1 {
			name := "gui"
			if commonUtils.GetOs() == "mac" {
				name = "zd.app"
			}

			if strings.Index(p, name) > -1 {
				guiDir := p[:strings.LastIndex(p, name)]
				dir = guiDir[:strings.LastIndex(guiDir, string(os.PathSeparator))]
			} else {
				dir = p[:strings.LastIndex(p, string(os.PathSeparator))]
			}
		}
	} else { // debug
		dir = GetDevDir()
	}

	dir, _ = filepath.Abs(dir)
	dir = AddSepIfNeeded(dir)

	//fmt.Printf("Debug: Launch %s in %s \n", arg1, dir)
	return dir
}

func GetDevDir() string { // where we run file in
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

func GetFileOrFolderDir(pth string) (ret string) {
	if IsDir(pth) {
		ret = pth
	} else {
		ret = filepath.Dir(pth)
	}

	ret, _ = filepath.Abs(ret)
	ret = AddSepIfNeeded(ret)

	return
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
		resPath := vari.GlobalVars.ConfigFileDir + from
		if !FileExist(resPath) { // in same folder with passed config file, like dir/name.yaml
			resPath = vari.ZdDir + from
			if !FileExist(resPath) { // in res file
				resPath = ""
			}
		}
		resFile = resPath
	}

	if !FileExist(resFile) {
		resFile = ""
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
			resPath = vari.ZdDir + file
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
		realPth1 := vari.ZdDir + constant.ResDirYaml + constant.PthSep + relatPath
		realPth2 := vari.ZdDir + constant.ResDirUsers + constant.PthSep + relatPath
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

			realPth := vari.ZdDir + constant.ResDirData + constant.PthSep + relatPath
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
		realPth := vari.ZdDir + constant.ResDirData + constant.PthSep +
			strings.Replace(from, ".", constant.PthSep, -1)
		if IsDir(realPth) {
			ret = realPth
			return
		}
	}

	return
}

func GetFilesByExtInDir(folder, ext string, files *[]string) {
	extArr := strings.Split(ext, ",")

	folder, _ = filepath.Abs(folder)

	if !IsDir(folder) {
		if ext == "" || stringUtils.StrInArr(filepath.Ext(folder), extArr) {
			*files = append(*files, folder)
		}

		return
	}

	dir, err := os.ReadDir(folder)
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
		} else if strings.Index(name, "~") != 0 && (ext == "" || stringUtils.StrInArr(filepath.Ext(filePath), extArr)) {
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

func GetFilesFromParams(args []string) (files []string, count int) {
	for _, arg := range args {
		if strings.Index(arg, "-") != 0 {
			files = append(files, arg)
			count++
		} else {
			break
		}
	}

	return
}

func HandleFiles(files []string) []string {
	if len(files) != 2 {
		return files
	}

	if files[0] == "" && files[1] != "" { // no defaultFile
		files[0] = files[1]
		files[1] = ""
	} else if files[1] == "" && files[0] != "" { // no configFile
		files[1] = files[0]
		files[0] = ""
	}

	return files
}

func NewFileNameWithUlidPostfix(pth string) (ret string) {
	return AddFilePostfix(pth, stringUtils.Ulid())
}

func AddFilePostfix(pth, postfix string) (ret string) {
	ext := filepath.Ext(pth)

	ret = pth[:strings.LastIndex(pth, ext)] + "-" + postfix + ext

	return
}

func GetUploadFileName(name string) (ret string, err error) {
	fns := strings.Split(strings.TrimPrefix(name, "./"), ".")
	if len(fns) < 2 {
		msg := fmt.Sprintf("文件名错误 %s", name)
		err = errors.New(msg)
		return
	}

	base := fns[0]
	ext := fns[1]

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	rand, _ := ulid.New(ms, entropy)

	ret = str.Join(base, "-", strings.ToLower(rand.String()), ".", ext)

	return
}
