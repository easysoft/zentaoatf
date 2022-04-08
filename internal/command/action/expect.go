package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	_scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/exec"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"path/filepath"
)

func GenExpectFiles(files []string) error {
	serverConfig.InitExecLog(commConsts.ExecLogDir)

	var cases []string
	for _, v1 := range files {
		group := scriptUtils.LoadScriptByWorkspace(v1)
		for _, v2 := range group {
			cases = append(cases, v2)
		}
	}

	if len(cases) < 1 {
		logUtils.Info("\n" + i118Utils.Sprintf("no_cases"))
		return nil
	}
	conf := configUtils.LoadByWorkspacePath(commConsts.WorkDir)
	casesToRun, _ := _scriptUtils.FilterCases(cases, conf)

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
	conf := configUtils.LoadByWorkspacePath(commConsts.WorkDir)
	out, _ := _scriptUtils.RunScript(file, commConsts.WorkDir, conf, nil, nil)

	expFile := filepath.Join(filepath.Dir(file), fileUtils.GetFileNameWithoutExt(file)+".exp")
	fileUtils.WriteFile(expFile, out)
}
