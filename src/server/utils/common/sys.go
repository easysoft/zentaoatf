package commonUtils

import (
	serverModel "github.com/easysoft/zentaoatf/src/server/model"
	"os"
	"runtime"
	"strings"
)

func GetSysInfo() (info serverModel.SysInfo) {
	info.SysArch = runtime.GOARCH
	info.SysCores = runtime.GOMAXPROCS(0)

	info.OsType = runtime.GOOS
	info.OsName, _ = os.Hostname()

	envs := os.Environ()
	for _, env := range envs {
		if strings.Index(env, "LC_CTYPE=") > -1 { // LC_CTYPE=zh_CN.UTF-8
			info.Lang = strings.Split(env, "=")[1]
		}
	}

	return
}
