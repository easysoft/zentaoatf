package script

import "github.com/easysoft/zentaoatf/src/utils"

func LoadTestAssets() ([]string, []string) {
	config := utils.ReadCurrConfig()
	ext := GetLangMap()[config.LangType]["extName"]

	caseFiles := make([]string, 0)
	suitesFiles := make([]string, 0)

	utils.GetAllFiles(utils.Prefer.WorkDir+utils.ScriptDir, ext, &caseFiles)
	utils.GetAllFiles(utils.Prefer.WorkDir+utils.ScriptDir, "suite", &suitesFiles)

	return caseFiles, suitesFiles
}
