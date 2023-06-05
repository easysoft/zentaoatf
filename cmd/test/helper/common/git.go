package commonTestHelper

import (
	"os"
	"os/exec"

	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	"github.com/go-git/go-git/v5"
)

func CloneGit(gitUrl string, name string) error {
	projectDir := name
	fileUtils.MkDirIfNeeded(projectDir)

	options := git.CloneOptions{
		URL:      gitUrl,
		Progress: os.Stdout,
	}
	_, err := git.PlainClone(projectDir, false, &options)
	return err
}

func Pull() (err error) {
	var cmd *exec.Cmd
	cmd = exec.Command("git", "pull")
	cmd.Dir = constTestHelper.RootPath
	_, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	return
}
