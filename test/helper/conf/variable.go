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
	ZentaoPassword = "Test123456."

	UiPort = 58000
	ZtfUrl = fmt.Sprintf("http://127.0.0.1:%d/", UiPort)

	ZentaoExtUrl = "https://www.zentao.net/file-download-22700.html"
)

func init() {
	if runtime.GOOS == "windows" {
		NewLine = "\r\n"
	}
	RootPath, _ = os.Getwd()
	if strings.Index(RootPath, "test") != -1 {
		RootPath = RootPath[:strings.Index(RootPath, "test")]
	}
	if RootPath[len(RootPath)-1:] != FilePthSep {
		RootPath += FilePthSep
	}
}
