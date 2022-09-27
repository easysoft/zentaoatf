package gitHelper

import (
	"errors"
	"os"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type Build struct {
	ScmAddress string
	Username   string
	Password   string
	WorkDir    string
	ProjectDir string
	RsaKey     string
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
	if build.RsaKey != "" {
		_, err = os.Stat(build.RsaKey)
		if err != nil {
			return
		}
		options.Auth, err = ssh.NewPublicKeysFromFile("git", build.RsaKey, password)
		if err != nil {
			return
		}
	}
	_, err = git.PlainClone(projectDir, false, &options)
	if err != nil {
		return
	}

	return
}

func GetGitProjectName(gitUrl string) string {
	index := strings.LastIndex(gitUrl, "/")

	name := gitUrl[index+1:]
	name = strings.Split(name, ".git")[0]
	return name
}

func Pull(build *Build) (err error) {
	userName := build.Username
	password := build.Password
	RsaKey := build.RsaKey
	if build.ProjectDir == "" {
		url := build.ScmAddress
		projectDir := build.WorkDir + GetGitProjectName(url) + commConsts.PthSep
		build.ProjectDir = projectDir
	}
	options := git.PullOptions{
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
	r, err := git.PlainOpen(build.ProjectDir)
	if err != nil {
		return
	}
	w, err := r.Worktree()
	if err != nil {
		return
	}
	err = w.Pull(&options)
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return
	}
	return nil
}

func ForcePull(build *Build, force bool) (err error) {
	userName := build.Username
	password := build.Password
	RsaKey := build.RsaKey
	if build.ProjectDir == "" {
		url := build.ScmAddress
		projectDir := build.WorkDir + GetGitProjectName(url) + commConsts.PthSep
		build.ProjectDir = projectDir
	}
	options := git.PullOptions{
		Progress: os.Stdout,
		Force:    force,
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
	r, err := git.PlainOpen(build.ProjectDir)
	if err != nil {
		return
	}
	w, err := r.Worktree()
	if err != nil {
		return
	}
	ref, err := r.Head()
	if err != nil {
		return
	}
	err = w.Reset(&git.ResetOptions{Commit: ref.Hash(), Mode: git.HardReset})
	if err != nil {
		return
	}
	err = w.Pull(&options)
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return
	}
	return nil
}

func Fetch(build *Build) (err error) {
	userName := build.Username
	password := build.Password
	RsaKey := build.RsaKey
	projectDir := build.WorkDir + commConsts.PthSep
	build.ProjectDir = projectDir
	fileUtils.MkDirIfNeeded(projectDir)
	options := git.FetchOptions{
		Progress: os.Stdout,
		Force:    true,
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
	if err != nil && errors.Is(err, git.NoErrAlreadyUpToDate) {
		return
	}
	return nil
}
