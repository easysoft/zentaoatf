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

func CheckoutCodes(build *serverModel.Build) {
	url := build.ScmAddress
	userName := build.ScmAccount
	password := build.ScmPassword

	build.ProjectDir = build.WorkDir + GetGitProjectName(url) + string(os.PathSeparator)

	fileUtils.MkDirIfNeeded(build.ProjectDir)

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
	_, err := git.PlainClone(build.ProjectDir, false, &options)
	if err != nil {
		return
	}
}

func DownloadCodes(build *serverModel.Build) {
	zipPath := build.WorkDir + uuid.NewV4().String() + ".zip"
	Download(build.ScriptUrl, zipPath)

	scriptFolder := GetZipSingleDir(zipPath)
	if scriptFolder != "" { // single dir in zip
		build.ProjectDir = build.WorkDir + scriptFolder
		archiver.Unarchive(zipPath, build.WorkDir)
	} else { // more then one dir, unzip to a folder
		fileNameWithoutExt := fileUtils.GetFileNameWithoutExt(zipPath)
		build.ProjectDir = build.WorkDir + fileNameWithoutExt + string(os.PathSeparator)
		archiver.Unarchive(zipPath, build.ProjectDir)
	}
}
