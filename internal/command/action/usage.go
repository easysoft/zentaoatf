package action

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/pkg/consts"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	resUtils "github.com/easysoft/zentaoatf/pkg/lib/res"
	"github.com/fatih/color"
	"os"
	"regexp"
)

var (
	usageFile  = fmt.Sprintf("res%sdoc%susage.txt", string(os.PathSeparator), consts.FilePthSep)
	sampleFile = fmt.Sprintf("res%sdoc%ssample.txt", consts.FilePthSep, string(os.PathSeparator))
)

func PrintUsage() {
	logUtils.Info("\n" + color.CyanString(i118Utils.Sprintf("usage")))

	usageData, _ := resUtils.ReadRes(usageFile)
	exeFile := commConsts.App
	if commonUtils.IsWin() {
		exeFile += ".exe"
	}
	usage := fmt.Sprintf(string(usageData), exeFile)
	fmt.Printf("%s", usage)

	logUtils.Info("\n" + color.CyanString(i118Utils.Sprintf("example")))

	sampleData, _ := resUtils.ReadRes(sampleFile)
	sample := string(sampleData)
	if !commonUtils.IsWin() {
		regx, _ := regexp.Compile(`\\`)
		sample = regx.ReplaceAllString(sample, "/")

		regx, _ = regexp.Compile(commConsts.App + `.exe`)
		sample = regx.ReplaceAllString(sample, commConsts.App)

		regx, _ = regexp.Compile(`/bat/`)
		sample = regx.ReplaceAllString(sample, "/shell/")

		regx, _ = regexp.Compile(`\.bat\s{4}`)
		sample = regx.ReplaceAllString(sample, ".shell")
	}
	fmt.Printf("%s\n", sample)
}
