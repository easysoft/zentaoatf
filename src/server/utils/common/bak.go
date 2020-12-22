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
	dateDir := vari.AgentLogDir + dateStr + constant.PthSep
	dist := dateDir + timeStr + ".zip"

	fileUtils.MkDirIfNeeded(vari.AgentLogDir)

	err := fileUtils.ZipFiles(dist, src)
	if err != nil {
		logUtils.Logger.Error(fmt.Sprintf("fail to zip test results '%s' to '%s', error %s", src, dist, err.Error()))
	}

	removeHistoryLog(vari.AgentLogDir)
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

	dirs, _ := ioutil.ReadDir(vari.AgentLogDir)

	for _, dir := range dirs {
		dirName := dir.Name()
		files, _ := ioutil.ReadDir(vari.AgentLogDir + dirName)

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
