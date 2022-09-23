package commonTestHelper

import (
	"os"
	"runtime"
	"strings"
)

const (
	FilePthSep = string(os.PathSeparator)
)

var (
	NewLine  = "/n"
	RootPath = ""
)

func init() {
	if runtime.GOOS == "windows" {
		NewLine = "\r\n"
	}
	apath, _ := os.Getwd()
	RootPath = apath[:strings.Index(apath, "test")]
}
