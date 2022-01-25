package analysisUtils

import (
	"encoding/json"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	"io/ioutil"
	"path/filepath"
)

func ListReport(projectPath string) (reportFiles []string) {
	dir := filepath.Join(projectPath, commConsts.LogDirName)

	files, _ := ioutil.ReadDir(dir)
	for _, fi := range files {
		if fi.IsDir() {
			reportFiles = append(reportFiles, fi.Name())
		}
	}

	return
}

func GetReport(projectPath string, seq string) (report commDomain.ZtfReport, err error) {
	pth := filepath.Join(projectPath, commConsts.LogDirName, seq, commConsts.ResultJson)

	content := fileUtils.ReadFileBuf(pth)
	err = json.Unmarshal(content, &report)

	return
}
