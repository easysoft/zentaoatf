package logUtils

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/snowlyg/helper/dir"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func GetLogDir(projectPath string) string {
	pth := filepath.Join(projectPath, commConsts.LogDirName)

	d, _ := ioutil.ReadDir(pth)

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

	if numb >= 9 {
		numb = 0

		tempDir := pth[:len(pth)-1] + "-bak" + string(os.PathSeparator) + pth[len(pth):]
		childDir := pth + "bak" + string(os.PathSeparator) + pth[len(pth):]

		os.RemoveAll(childDir)
		os.Rename(pth, tempDir)

		err := os.Rename(tempDir, childDir)
		_ = err
	}

	num := getLogNumb(numb + 1)
	ret := addPathSepIfNeeded(filepath.Join(pth, num))

	if !dir.IsExist(ret) { // 判断是否有Director文件夹
		dir.InsureDir(ret)
	}

	return ret
}

func getLogNumb(numb int) string {
	return fmt.Sprintf("%03s", strconv.Itoa(numb))
}

func addPathSepIfNeeded(pth string) string { // not to call fileUtils cycle
	sep := consts.PthSep

	if strings.LastIndex(pth, sep) < len(pth)-1 {
		pth += sep
	}
	return pth
}
