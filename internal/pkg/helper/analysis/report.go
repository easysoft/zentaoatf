package analysisHelper

import (
	"encoding/json"
	"path/filepath"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	"github.com/ergoapi/util/file"

	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
)

func ListReport(workspacePath string) (reportSeqs []string) {
	return ListReport2(workspacePath, -1)
}

func ListReport2(workspacePath string, maxSize int) (reportSeqs []string) {
	logRoot := filepath.Join(workspacePath, commConsts.LogDirName)

	if maxSize <= 0 {
		maxSize = int(^uint(0) >> 1)
	}

	var count int = 0
	// read log dir
	dailyFiles, _ := file.ReadDir(logRoot)
	for i := len(dailyFiles) - 1; i > -1; i-- {
		daily := dailyFiles[i]

		if daily.IsDir() {
			// read daily log dir
			files, _ := file.ReadDir(filepath.Join(logRoot, daily.Name()))
			for j := len(files) - 1; j > -1; j-- {
				count++
				if count > maxSize {
					break
				}

				fi := files[j]
				reportSeqs = append(reportSeqs, EncodeSeq(filepath.Join(daily.Name(), fi.Name())))
			}
		}
	}

	return
}

func EncodeSeq(seq string) string {
	return strings.ReplaceAll(seq, string(filepath.Separator), ":")
}

func DecodeSeq(seq string) string {
	return strings.ReplaceAll(seq, ":", string(filepath.Separator))
}

func ReadReportByWorkspaceSeq(workspacePath string, resultDir string) (report commDomain.ZtfReport, pth string, err error) {
	if commConsts.ExecFrom == commConsts.FromCmd {
		seqPath := resultDir
		if !filepath.IsAbs(seqPath) {
			seqPath = filepath.Join(workspacePath, resultDir)
		}

		pth = filepath.Join(seqPath, commConsts.ResultJson)
	} else {
		resultDir = DecodeSeq(resultDir)
		pth = filepath.Join(workspacePath, commConsts.LogDirName, resultDir, commConsts.ResultJson)
	}

	report, err = ReadReportByPath(pth)

	return
}

func ReadReportByPath(pth string) (report commDomain.ZtfReport, err error) {
	content := fileUtils.ReadFileBuf(pth)
	if commConsts.ExecFrom == commConsts.FromCmd {
		contentData := strings.Replace(string(content), "\n", "", -1)
		contentData = strings.ReplaceAll(contentData, "\"status\":false", "\"status\":\"fail\"")
		contentData = strings.ReplaceAll(contentData, "\"status\":true", "\"status\":\"pass\"")
		contentData = strings.ReplaceAll(contentData, "\"env\"", "\"platform\"")
		contentData = strings.ReplaceAll(contentData, "\"testEnv\"", "\"platform\"")
		contentData = strings.ReplaceAll(contentData, "\"testFrame\"", "\"testTool\"")
		content = []byte(contentData)
	}

	err = json.Unmarshal(content, &report)

	return
}
