package commonUtils

import (
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/const"
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/emirpasic/gods/maps"
	"net"
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

func AddSlashForUrl(url string) string {
	if strings.LastIndex(url, "/") < len(url)-1 {
		url += "/"
	}

	return url
}

func IgnoreFile(path string) bool {
	path = filepath.Base(path)

	if strings.Index(path, ".") == 0 ||
		path == "bin" || path == "release" || path == "logs" || path == "xdoc" {
		return true
	} else {
		return false
	}
}

func GetFieldVal(config model.Config, key string) string {
	key = stringUtils.UcFirst(key)

	immutable := reflect.ValueOf(config)
	val := immutable.FieldByName(key).String()

	return val
}

func SetFieldVal(config *model.Config, key string, val string) string {
	key = stringUtils.UcFirst(key)

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

			ipType := getIpType(ip)
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

func getIpType(IP net.IP) string {
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

func GetDebugParamForRun(args []string) (debug string, ret []string) {
	index := -1
	for i, item := range args {
		if item == "-debug" || item == "--debug" {
			index = i
			break
		}
	}

	if index > -1 && len(args) > index+1 {
		debug = args[index+1]
		ret = append(args[0:index], args[index+2:]...)
	} else {
		ret = args
	}

	return
}
