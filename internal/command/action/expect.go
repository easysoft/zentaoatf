package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/comm/helper/config"
	execHelper "github.com/easysoft/zentaoatf/internal/comm/helper/exec"
	scriptHelper "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	"path/filepath"
)

func GenExpectFiles(files []string) error {
	serverConfig.InitExecLog(commConsts.ExecLogDir)

	var cases []string
	for _, v1 := range files {
		group := scriptHelper.LoadScriptByWorkspace(v1)
		for _, v2 := range group {
			cases = append(cases, v2)
		}
	}

	if len(cases) < 1 {
		logUtils.Info("\n" + i118Utils.Sprintf("no_cases"))
		return nil
	}
	conf := configHelper.LoadByWorkspacePath(commConsts.WorkDir)
	casesToRun, _ := execHelper.FilterCases(cases, conf)

	dryRunScripts(casesToRun)
	logUtils.Info(i118Utils.Sprintf("success_to_create_expect"))

	return nil
}

func dryRunScripts(casesToRun []string) {
	for _, file := range casesToRun {
		dryRunScript(file)
	}
}
func dryRunScript(file string) {
	conf := configHelper.LoadByWorkspacePath(commConsts.WorkDir)
	out, _ := execHelper.RunScript(file, commConsts.WorkDir, conf, nil, nil)

	expFile := filepath.Join(filepath.Dir(file), fileUtils.GetFileNameWithoutExt(file)+".exp")
	fileUtils.WriteFile(expFile, out)
}
