package fileUtils

import (
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	shellUtils "github.com/easysoft/zentaoatf/pkg/lib/shell"
	"strings"
)

func GetMd5(pth string) (ret string, err error) {
	cmdStr := ""
	if commonUtils.IsWin() {
		cmdStr = "CertUtil -hashfile " + pth + " MD5"
	} else {
		cmdStr = "md5sum " + pth + " | awk '{print $1}'"
	}

	ret, _ = shellUtils.ExeSysCmd(cmdStr)

	if commonUtils.IsWin() {
		arr := strings.Split(ret, "\n")
		if len(arr) > 1 {
			ret = arr[1]
		}
	}

	ret = strings.TrimSpace(ret)

	return
}
