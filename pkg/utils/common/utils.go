package commonUtils

import (
	"github.com/easysoft/zendata/internal/pkg/const"
	"github.com/easysoft/zendata/internal/pkg/model"
	stringUtils "github.com/easysoft/zendata/pkg/utils/string"
	"github.com/emirpasic/gods/maps"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func Base(pathStr string) string {
	pathStr = filepath.ToSlash(pathStr)
	return filepath.Base(pathStr)
}

func RemoveBlankLine(str string) string {
	myExp := regexp.MustCompile(`\n{3,}`) // 连续换行
	ret := myExp.ReplaceAllString(str, "\n\n")

	ret = strings.Trim(ret, "\n")
	//ret = strings.TrimSpace(ret)

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
		return strings.Replace(osName, "windows", "win", 1)
	}
}
func IsWin() bool {
	return GetOs() == "win"
}
func IsLinux() bool {
	return GetOs() == "linux"
}
func IsMac() bool {
	return GetOs() == "mac"
}

func IsRelease() bool {
	arg1 := strings.ToLower(os.Args[0])
	name := filepath.Base(arg1)
	return strings.Index(name, constant.AppName) == 0 && strings.Index(arg1, "go-build") < 0

	//if _, err := os.Stat("res"); os.IsNotExist(err) {
	//	return true
	//}
	//
	//return false
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

func GetIp() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	ipMap := map[string]string{}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}

			ipType := GetIpType(ip)
			ipMap[ipType] = ip.String()
		}
	}

	if ipMap["public"] != "" {
		return ipMap["public"]
	} else if ipMap["private"] != "" {
		return ipMap["private"]
	} else {
		return ""
	}
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func GetIpType(IP net.IP) string {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return ""
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return "private"
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return "private"
		case ip4[0] == 192 && ip4[1] == 168:
			return "private"
		default:
			return "public"
		}
	}
	return ""
}

func RandNum(length int) int {
	randSeeds := time.Now().Unix() + int64(rand.Intn(100000000))
	rand.Seed(randSeeds)

	seedInt := rand.Intn(length)
	return seedInt
}
func RandNum64(length int64) int64 {
	randSeeds := time.Now().Unix() + int64(rand.Intn(100000000))
	rand.Seed(randSeeds)

	seedInt := rand.Int63n(length)
	return seedInt
}
