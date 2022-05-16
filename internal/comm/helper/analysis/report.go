package analysisHelper

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
)

func ListReport(workspaceLogPath string) (reportFiles []string) {
	//dir := filepath.Join(workspacePath, commConsts.LogDirName)

	files, _ := ioutil.ReadDir(workspaceLogPath)
	for _, fi := range files {
		if fi.IsDir() {
			reportFiles = append(reportFiles, fi.Name())
		}
	}

	return
}

func ListReportByModTime(workspaceLogPath string) (reportFiles []fs.FileInfo) {

	files, _ := ioutil.ReadDir(workspaceLogPath)
	reportFiles = make([]fs.FileInfo, 0, len(files))
	for _, fi := range files {
		if fi.IsDir() {
			reportFiles = append(reportFiles, fi)
		}
	}

	sort.Slice(reportFiles, func(i, j int) bool {
		return reportFiles[i].ModTime().After(reportFiles[j].ModTime())
	})

	return
}

func ReadReportByWorkspaceSeq(workspacePath string, seq string) (report commDomain.ZtfReport, err error) {

	report, err = ReadReportByWorkspaceSeq2(workspacePath, seq, false)
	if err != nil {
		report, err2 := ReadReportByWorkspaceSeq2(workspacePath, seq, true)

		if err2 != nil {
			return report, err
		} else {
			return report, nil
		}
	}

	return
}

func ReadReportByWorkspaceSeq2(workspacePath string, seq string, isBak bool) (report commDomain.ZtfReport, err error) {
	pth := ""
	if commConsts.ExecFrom == commConsts.FromCmd {
		seqPath := seq
		if !filepath.IsAbs(seqPath) {
			seqPath = filepath.Join(workspacePath, seq)
		}
		pth = filepath.Join(seqPath, commConsts.ResultJson)
	} else {
		if isBak {
			pth = filepath.Join(workspacePath, commConsts.LogBakDirName, seq, commConsts.ResultJson)
		} else {
			pth = filepath.Join(workspacePath, commConsts.LogDirName, seq, commConsts.ResultJson)
		}
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
