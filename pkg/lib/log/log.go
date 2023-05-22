package logUtils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/pkg/consts"
	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"

	"github.com/snowlyg/helper/dir"
)

func GetLogDir(workspacePath string) string {
	logBase := filepath.Join(workspacePath, commConsts.LogDirName)

	days := geWeekDays()
	files1, _ := ioutil.ReadDir(logBase)
	for _, fi := range files1 {
		name := fi.Name()
		if fi.IsDir() && !stringUtils.FindInArr(name, days) {
			os.RemoveAll(filepath.Join(logBase, name))
		}
	}

	logDir := filepath.Join(logBase, dateUtils.DateStrShort(time.Now()))
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.MkdirAll(logDir, os.ModePerm)
	}

	files2, _ := ioutil.ReadDir(logDir)
	regx := `^\d\d\d$`
	numb := 0
	for _, fi := range files2 {
		if !fi.IsDir() {
			continue
		}

		name := fi.Name()
		isLog, _ := regexp.MatchString(regx, name)
		if !isLog {
			continue
		}

		name = strings.TrimLeft(name, "0")
		nm, _ := strconv.Atoi(name)

		if nm >= numb {
			numb = nm
		}
	}

	num := getLogNumb(numb + 1)
	ret := addPathSepIfNeeded(filepath.Join(logDir, num))

	if !dir.IsExist(ret) { // 判断是否有Director文件夹
		dir.InsureDir(ret)
	}

	return ret
}

func geWeekDays() (ret []string) {
	for i := 0; i < 7; i++ {
		today := time.Now()
		newDay := today.AddDate(0, 0, i*-1)
		newDayStr := dateUtils.DateStrShort(newDay)

		ret = append(ret, newDayStr)
	}

	return
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
