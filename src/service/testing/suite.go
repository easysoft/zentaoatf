package testingService

import (
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"strings"
)

func GenSuite(cases []string) {
	str := strings.Join(cases, "\n")

	fileUtils.WriteFile(vari.Prefer.WorkDir+constant.ScriptDir+"all."+constant.SuiteExt, str)
}
