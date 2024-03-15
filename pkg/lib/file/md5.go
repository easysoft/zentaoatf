package fileUtils

import (
	"strings"

	shellUtils "github.com/easysoft/zentaoatf/pkg/lib/shell"
	"github.com/ergoapi/util/zos"
)

func GetMd5(pth string) (ret string, err error) {
	cmdStr := ""
	if zos.IsUnix() {
		cmdStr = "md5sum " + pth + " | awk '{print $1}'"
	} else {
		cmdStr = "CertUtil -hashfile " + pth + " MD5"
	}

	ret, _ = shellUtils.ExeSysCmd(cmdStr)

	if !zos.IsUnix() {
		arr := strings.Split(ret, "\n")
		if len(arr) > 1 {
			ret = arr[1]
		}
	}

	ret = strings.TrimSpace(ret)

	return
}
