package constTestHelper

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	FilePthSep = string(os.PathSeparator)
)

var (
	NewLine  = "\n"
	RootPath = ""

	ZentaoPort     = 58080
	ZentaoSiteUrl  = fmt.Sprintf("http://127.0.0.1:%d", ZentaoPort)
	ZentaoUsername = "admin"
	ZentaoPassword = "P2ssw0rd"

	UiPort = 58000
	ZtfUrl = fmt.Sprintf("http://127.0.0.1:%d/", UiPort)

	ZentaoExtUrl = "https://www.zentao.net/file-download-22700.html"

	WorkspaceName = "单元测试工作目录"
	SiteName      = "单元测试站点"
)

func init() {
	if runtime.GOOS == "windows" {
		NewLine = "\r\n"
	}
	RootPath, _ = os.Getwd()
	if strings.Index(RootPath, "test") != -1 {
		RootPath = RootPath[:strings.Index(RootPath, "test")-4]
	}
	if RootPath[len(RootPath)-1:] != FilePthSep {
		RootPath += FilePthSep
	}
}
