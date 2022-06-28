package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	execHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
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
	conf := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)

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
	conf := configHelper.LoadByWorkspacePath(commConsts.ZtfDir)

	out, _ := execHelper.RunFile(file, commConsts.WorkDir, conf, nil, nil)

	expFile := filepath.Join(filepath.Dir(file), fileUtils.GetFileNameWithoutExt(file)+".exp")
	fileUtils.WriteFile(expFile, out)
}
