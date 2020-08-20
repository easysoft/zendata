package fileUtils

import (
	"github.com/easysoft/zendata/res"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/vari"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func ReadFile(filePath string) string {
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

func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func AbosutePath(pth string) string {
	if !IsAbosutePath(pth) {
		pth, _ = filepath.Abs(pth)
	}

	pth = AddSepIfNeeded(pth)

	return pth
}

func IsAbosutePath(pth string) bool {
	return path.IsAbs(pth) ||
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
	abs := ""
	if !IsAbosutePath(path) {
		path = vari.WorkDir + path
	}

	abs, _ = filepath.Abs(filepath.Dir(path))
	abs = AddSepIfNeeded(abs)
	return abs
}

func GetResProp(from string) (resFile, resType, sheet string) { // from resource

	index := strings.LastIndex(from, ".yaml")
	if index > -1 { // yaml, ip.v1.yaml
		resFile = convertYamlPath(from)
		resType = "yaml"
	} else { // excel, like address.cn.v1.china
		resFile, sheet = convertExcelPath(from)
		resType = "excel"
	}

	if resFile == "" {
		resPath := vari.ConfigDir + resFile
		if !FileExist(resPath) { // in same folder with passed config file

			resPath = vari.WorkDir + resFile
			if !FileExist(resPath) {  // in res file
				resPath = ""
			}
		}
		resFile = resPath
	}

	return
}

func convertYamlPath(from string) (ret string) {
	arr := strings.Split(from, ".")
	for i := 0; i < len(arr); i++ {
		dir := ""
		if i > 0 {
			dir = strings.Join(arr[:i], constant.PthSep)
		}
		file := strings.Join(arr[i:], ".")

		if dir != "" {
			ret = dir + constant.PthSep + file
		} else {
			ret = file
		}

		realPth1 := vari.WorkDir + constant.ResDirYaml + constant.PthSep + ret
		realPth2 := vari.WorkDir + constant.ResDirUsers + constant.PthSep + ret
		if FileExist(realPth1) {
			ret = realPth1
			break
		} else if FileExist(realPth2) {
			ret = realPth2
			break
		}
	}

	return
}

func convertExcelPath(from string) (ret, sheet string) {
	path1 := from // address.cn.v1
	index := strings.LastIndex(from, ".")
	path2 := from[:index] // address.cn.v1.china

	paths := [2]string{path1, path2}
	for index, path := range paths {

		arr := strings.Split(path, ".")
		for i := 0; i < len(arr); i++ {
			dir := ""
			if i > 0 {
				dir = strings.Join(arr[:i], constant.PthSep)
			}
			file := strings.Join(arr[i:], ".") + ".xlsx"

			if dir != "" {
				ret = dir + constant.PthSep + file
			} else {
				ret = file
			}

			realPth := vari.WorkDir + constant.ResDirData + constant.PthSep + ret
			if FileExist(realPth) {
				if index == 1 {
					sheet = from[strings.LastIndex(from, ".")+1:]
				}
				ret = realPth
				return
			}
		}
	}

	return
}