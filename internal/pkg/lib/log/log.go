package logUtils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/snowlyg/helper/dir"
)

func GetLogDir(workspacePath string) string {
	logDir := filepath.Join(workspacePath, commConsts.LogDirName)

	d, _ := ioutil.ReadDir(logDir)

	regx := `^\d\d\d$`

	numb := 0
	for _, fi := range d {
		if fi.IsDir() {
			name := fi.Name()
			pass, _ := regexp.MatchString(regx, name)

			if pass { // 999
				name = strings.TrimLeft(name, "0")
				nm, _ := strconv.Atoi(name)

				if nm >= numb {
					numb = nm
				}
			}
		}
	}

	if numb > 9 {
		numb = 0

		tempDir := logDir[:len(logDir)-1] + "-bak" + string(os.PathSeparator) + logDir[len(logDir):]
		childDir := logDir + "-bak" + string(os.PathSeparator) + logDir[len(logDir):]

		if err := os.RemoveAll(childDir); err != nil {
			panic(err)
		}

		if err := os.Rename(logDir, tempDir); err != nil {
			panic(err)
		}

		if err := os.Rename(tempDir, childDir); err != nil {
			panic(err)
		}
	}

	num := getLogNumb(numb + 1)
	ret := addPathSepIfNeeded(filepath.Join(logDir, num))

	if !dir.IsExist(ret) { // 判断是否有Director文件夹
		dir.InsureDir(ret)
	}

	return ret
}

func getLogNumb(numb int) string {
	return fmt.Sprintf("%03s", strconv.Itoa(numb))
}

func addPathSepIfNeeded(pth string) string { // not to call fileUtils cycle
	sep := consts.FilePthSep

	if strings.LastIndex(pth, sep) < len(pth)-1 {
		pth += sep
	}
	return pth
}
