package commonUtils

import (
	"github.com/easysoft/zendata/src/model"
	"github.com/easysoft/zendata/src/utils/const"
	stringUtils "github.com/easysoft/zendata/src/utils/string"
	"github.com/emirpasic/gods/maps"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

func Base(pathStr string) string {
	pathStr = filepath.ToSlash(pathStr)
	return path.Base(pathStr)
}

func RemoveBlankLine(str string) string {
	myExp := regexp.MustCompile(`\n{3,}`) // 连续换行
	ret := myExp.ReplaceAllString(str, "\n\n")

	ret = strings.Trim(ret, "\n")
	ret = strings.TrimSpace(ret)

	return ret
}

func BoolToPass(b bool) string {
	if b {
		return constant.PASS.String()
	} else {
		return constant.FAIL.String()
	}
}

func GetOs() string {
	osName := runtime.GOOS

	if osName == "darwin" {
		return "mac"
	} else {
		return osName
	}
}
func IsWin() bool {
	return GetOs() == "windows"
}
func IsLinux() bool {
	return GetOs() == "linux"
}
func IsMac() bool {
	return GetOs() == "mac"
}

func IsRelease() bool {
	if _, err := os.Stat("res"); os.IsNotExist(err) {
		return true
	}

	return false
}

func UpdateUrl(url string) string {
	if strings.LastIndex(url, "/") < len(url)-1 {
		url += "/"
	}

	return url
}

func IngoreFile(path string) bool {
	path = filepath.Base(path)

	if strings.Index(path, ".") == 0 ||
		path == "bin" || path == "release" || path == "logs" || path == "xdoc" {
		return true
	} else {
		return false
	}
}

func GetFieldVal(config model.Config, key string) string {
	key = stringUtils.Ucfirst(key)

	immutable := reflect.ValueOf(config)
	val := immutable.FieldByName(key).String()

	return val
}

func SetFieldVal(config *model.Config, key string, val string) string {
	key = stringUtils.Ucfirst(key)

	mutable := reflect.ValueOf(config).Elem()
	mutable.FieldByName(key).SetString(val)

	return val
}

func LinkedMapToMap(mp maps.Map) map[string]string {
	ret := make(map[string]string, 0)

	for _, keyIfs := range mp.Keys() {
		valueIfs, _ := mp.Get(keyIfs)

		key := strings.TrimSpace(keyIfs.(string))
		value := strings.TrimSpace(valueIfs.(string))

		ret[key] = value
	}

	return ret
}
