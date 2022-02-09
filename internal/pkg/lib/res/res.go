package resUtils

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/res"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func ReadRes(path string) (ret []byte, err error) {
	isRelease := commonUtils.IsRelease()

	if isRelease {
		ret, err = res.Asset(path)
	} else {
		ret, err = ioutil.ReadFile(path)
	}

	dir, _ := os.Getwd()

	msg := fmt.Sprintf("isRelease=%t, path=%s, dir=%s", isRelease, path, dir)
	if logUtils.LoggerStandard != nil {
		logUtils.Info(msg)
	} else {
		log.Println(msg)
	}

	return
}

var (
	usageFile  = fmt.Sprintf("res%sdoc%susage.txt", string(os.PathSeparator), consts.PthSep)
	sampleFile = fmt.Sprintf("res%sdoc%ssample.txt", consts.PthSep, string(os.PathSeparator))
)

func PrintUsage() {
	logUtils.ExecConsolef(color.FgCyan, "Usage: ")

	usageData, _ := ReadRes(usageFile)
	exeFile := commConsts.App
	if commonUtils.IsWin() {
		exeFile += ".exe"
	}
	usage := fmt.Sprintf(string(usageData), exeFile)
	fmt.Printf("%s\n", usage)

	logUtils.ExecConsole(color.FgCyan, "\nExample: ")
	sampleData, _ := ReadRes(sampleFile)
	sample := ""
	if !commonUtils.IsWin() {
		regx, _ := regexp.Compile(`\\`)
		sample = regx.ReplaceAllString(string(sampleData), "/")

		regx, _ = regexp.Compile(commConsts.App + `.exe`)
		sample = regx.ReplaceAllString(sample, commConsts.App)

		regx, _ = regexp.Compile(`/bat/`)
		sample = regx.ReplaceAllString(sample, "/shell/")

		regx, _ = regexp.Compile(`\.bat\s{4}`)
		sample = regx.ReplaceAllString(sample, ".shell")
	}
	fmt.Printf("%s\n", sample)
}
