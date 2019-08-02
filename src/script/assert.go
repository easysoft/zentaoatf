package script

import "github.com/easysoft/zentaoatf/src/utils"

func LoadTestAssets() ([]string, []string) {
	config := utils.ReadCurrConfig()
	ext := GetLangMap()[config.LangType]["extName"]

	caseFiles, _ := utils.GetAllFiles(utils.Prefer.WorkDir+utils.ScriptDir, ext)
	suitesFiles, _ := utils.GetAllFiles(utils.Prefer.WorkDir+utils.ScriptDir, "suite")

	return caseFiles, suitesFiles
}
