package scriptService

import (
	config2 "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
)

func LoadAssetFiles() ([]string, []string) {
	config := config2.ReadCurrConfig()
	ext := GetSupportedScriptLang()[config.LangType]["extName"]

	caseFiles := make([]string, 0)
	suitesFiles := make([]string, 0)

	fileUtils.GetAllFiles(vari.Prefer.WorkDir+constant.ScriptDir, ext, &caseFiles)
	fileUtils.GetAllFiles(vari.Prefer.WorkDir+constant.ScriptDir, "suite", &suitesFiles)

	return caseFiles, suitesFiles
}
