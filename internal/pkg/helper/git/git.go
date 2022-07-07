package gitHelper

import (
	"os"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/mholt/archiver/v3"
	uuid "github.com/satori/go.uuid"
)

type Build struct {
	ScriptUrl  string
	ScmAddress string
	Username   string
	Password   string
	WorkDir    string
	ProjectDir string
	RsaKey     string
}

func GetTestScript(build *Build) (err error) {
	if build.ScmAddress != "" {
		err = CheckoutCodes(build)
	} else if strings.Index(build.ScriptUrl, "http://") == 0 {
		err = DownloadCodes(build)
	}
	//else {
	//	build.ProjectDir = fileUtils.AddPathSepIfNeeded(build.ScriptUrl)
	//}

	return
}

func CheckoutCodes(build *Build) (err error) {
	url := build.ScmAddress
	userName := build.Username
	password := build.Password

	projectDir := build.WorkDir + GetGitProjectName(url) + commConsts.PthSep
	build.ProjectDir = projectDir
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
	_, err = git.PlainClone(projectDir, false, &options)
	if err != nil {
		return
	}

	return
}

func DownloadCodes(build *Build) (err error) {
	zipPath := build.WorkDir + uuid.NewV4().String() + fileUtils.GetExtName(build.ScriptUrl)
	err = fileUtils.Download(build.ScriptUrl, zipPath)

	if err != nil {
		return
	}

	scriptFolder := fileUtils.GetZipSingleDir(zipPath)
	if scriptFolder != "" { // single dir in zip
		build.ProjectDir = build.WorkDir + scriptFolder
		err = archiver.Unarchive(zipPath, build.WorkDir)
	} else { // more then one dir, unzip to a folder
		fileNameWithoutExt := fileUtils.GetFileNameWithoutExt(zipPath)
		build.ProjectDir = build.WorkDir + fileNameWithoutExt + commConsts.PthSep
		err = archiver.Unarchive(zipPath, build.ProjectDir)
	}

	return
}

func GetGitProjectName(gitUrl string) string {
	index := strings.LastIndex(gitUrl, "/")

	name := gitUrl[index+1:]
	name = strings.Split(name, ".git")[0]
	return name
}

func CheckAuth(build *Build) (err error) {
	userName := build.Username
	password := build.Password
	RsaKey := build.RsaKey
	projectDir := build.WorkDir + commConsts.PthSep
	build.ProjectDir = projectDir
	fileUtils.MkDirIfNeeded(projectDir)
	options := git.FetchOptions{
		Progress: os.Stdout,
	}
	if userName != "" {
		options.Auth = &http.BasicAuth{
			Username: userName,
			Password: password,
		}
	}
	if RsaKey != "" {
		_, err = os.Stat(RsaKey)
		if err != nil {
			return
		}
		options.Auth, err = ssh.NewPublicKeysFromFile("git", RsaKey, password)
		if err != nil {
			return
		}
	}
	r, err := git.PlainOpen(projectDir)
	if err != nil {
		return
	}
	remote, err := r.Remote("origin")
	if err != nil {
		return
	}
	err = remote.Fetch(&options)
	if err != nil && err.Error() != "already up-to-date" {
		return
	}
	return nil
}
