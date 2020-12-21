package serverUtils

import (
	"fmt"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	dateUtils "github.com/easysoft/zentaoatf/src/utils/date"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io/ioutil"
	"path"
	"time"
)

func BakLog(src string) {
	now := time.Now()
	dateStr := dateUtils.DateStr(now)
	timeStr := dateUtils.TimeStr(now)
	logDir := vari.ZTFDir + "log-agent" + constant.PthSep
	dateDir := logDir + dateStr + constant.PthSep
	dist := dateDir + timeStr + ".zip"

	fileUtils.MkDirIfNeeded(logDir)

	err := fileUtils.ZipFiles(dist, src)
	if err != nil {
		logUtils.Logger.Error(fmt.Sprintf("fail to zip test results '%s' to '%s', error %s", src, dist, err.Error()))
	}

	removeHistoryLog(logDir)
}
func removeHistoryLog(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, fi := range files {
		name := fi.Name()
		tm, err := dateUtils.StrToDate(name)
		if err == nil && time.Now().Unix()-tm.Unix() > 7*24*3600 {
			fileUtils.RmDir(dir + name)
		}
	}
}

func ListHistoryLog() (ret []map[string]string) {
	logDir := vari.ServerWorkDir + "log-agent" + constant.PthSep

	dirs, _ := ioutil.ReadDir(logDir)

	for _, dir := range dirs {
		dirName := dir.Name()
		files, _ := ioutil.ReadDir(logDir + dirName)

		for _, fi := range files {
			name := fi.Name()
			if path.Ext(name) != ".zip" {
				continue
			}

			item := map[string]string{"name": dirName + constant.PthSep + name}
			ret = append(ret, item)
		}
	}

	return
}
