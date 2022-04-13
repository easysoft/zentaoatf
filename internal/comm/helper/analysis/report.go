package analysisHelper

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/easysoft/zentaoatf/internal/comm/domain"
	"github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func ListReport(workspacePath string) (reportFiles []string) {
	dir := filepath.Join(workspacePath, commConsts.LogDirName)

	files, _ := ioutil.ReadDir(dir)
	for _, fi := range files {
		if fi.IsDir() {
			reportFiles = append(reportFiles, fi.Name())
		}
	}

	return
}

func ReadReportByWorkspaceSeq(workspacePath string, seq string) (report commDomain.ZtfReport, err error) {
	pth := ""
	if commConsts.ExecFrom == commConsts.FromCmd {
		seqPath := seq
		if !filepath.IsAbs(seqPath) {
			seqPath = filepath.Join(workspacePath, seq)
		}
		pth = filepath.Join(seqPath, commConsts.ResultJson)
	} else {
		pth = filepath.Join(workspacePath, commConsts.LogDirName, seq, commConsts.ResultJson)
	}

	return ReadReportByPath(pth)
}

func ReadReportByPath(pth string) (report commDomain.ZtfReport, err error) {
	content := fileUtils.ReadFileBuf(pth)
	if commConsts.ExecFrom == commConsts.FromCmd {
		contentData := strings.Replace(string(content), "\n", "", -1)
		contentData = strings.ReplaceAll(contentData, "\"status\":false", "\"status\":\"fail\"")
		contentData = strings.ReplaceAll(contentData, "\"status\":true", "\"status\":\"pass\"")
		contentData = strings.ReplaceAll(contentData, "\"env\"", "\"testEnv\"")
		contentData = strings.ReplaceAll(contentData, "\"testFrame\"", "\"testTool\"")
		content = []byte(contentData)
	}

	err = json.Unmarshal(content, &report)

	return
}
