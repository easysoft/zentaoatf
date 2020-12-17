package serverUtils

import (
	"strings"
)

func GetGitProjectName(gitUrl string) string {
	index := strings.LastIndex(gitUrl, "/")

	name := gitUrl[index+1:]
	name = strings.Split(name, ".git")[0]
	return name
}
