package serverUtils

import (
	serverModel "github.com/easysoft/zentaoatf/src/server/domain"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/mholt/archiver/v3"
	uuid "github.com/satori/go.uuid"
	"os"
	"strings"
)

func GetTestScript(build *serverModel.Build) (result serverModel.OptResult) {
	if build.ScmAddress != "" {
		CheckoutCodes(build)
	} else if strings.Index(build.ScriptUrl, "http://") == 0 {
		DownloadCodes(build)
	} else {
		build.ProjectDir = fileUtils.AddPathSepIfNeeded(build.ScriptUrl)
	}

	result.Success("")
	return result
}

func CheckoutCodes(task *serverModel.Build) {
	url := task.ScmAddress
	userName := task.ScmAccount
	password := task.ScmPassword

	projectDir := task.WorkDir + GetGitProjectName(url) + string(os.PathSeparator)

	fileUtils.MkDirIfNeeded(projectDir)

	options := git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}
	if userName != "" {
		options.Auth = &http.BasicAuth{
			Username: userName,
			Password: password,
		}
	}
	_, err := git.PlainClone(projectDir, false, &options)
	if err != nil {
		return
	}

	task.ProjectDir = projectDir
}

func DownloadCodes(task *serverModel.Build) {
	zipPath := task.WorkDir + uuid.NewV4().String() + fileUtils.GetExtName(task.ScriptUrl)
	Download(task.ScriptUrl, zipPath)

	scriptFolder := GetZipSingleDir(zipPath)
	if scriptFolder != "" { // single dir in zip
		task.ProjectDir = task.WorkDir + scriptFolder
		archiver.Unarchive(zipPath, task.WorkDir)
	} else { // more then one dir, unzip to a folder
		fileNameWithoutExt := fileUtils.GetFileNameWithoutExt(zipPath)
		task.ProjectDir = task.WorkDir + fileNameWithoutExt + string(os.PathSeparator)
		archiver.Unarchive(zipPath, task.ProjectDir)
	}
}
