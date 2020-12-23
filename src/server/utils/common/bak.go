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
	"regexp"
	"time"
)

func BakLog(src string) {
	now := time.Now()
	dateStr := dateUtils.DateStrNoSep(now)
	timeStr := dateUtils.TimeStrNoSep(now)
	dateDir := vari.AgentLogDir + dateStr + constant.PthSep
	dist := dateDir + timeStr + ".zip"

	fileUtils.MkDirIfNeeded(vari.AgentLogDir)

	err := fileUtils.ZipFiles(dist, src)
	if err != nil {
		logUtils.Logger.Error(fmt.Sprintf("fail to zip test results '%s' to '%s', error %s", src, dist, err.Error()))
	}

	removeHistoryLog(vari.AgentLogDir)
}
func removeHistoryLog(root string) {
	dirs, _ := ioutil.ReadDir(root)

	for _, dir := range dirs {
		name := dir.Name()
		pass, _ := regexp.MatchString(`^[0-9]{8}$`, name)
		if !pass {
			continue
		}

		tm, err := time.Parse("20060102", name)
		if err == nil && time.Now().Unix()-tm.Unix() > 7*24*3600 {
			fileUtils.RmDir(root + name)
		}
	}
}

func ListHistoryLog() (ret []map[string]string) {
	dirs, _ := ioutil.ReadDir(vari.AgentLogDir)

	for _, dir := range dirs {
		dirName := dir.Name()
		pass, _ := regexp.MatchString(`^[0-9]{8}$`, dirName)
		if !pass {
			continue
		}

		files, _ := ioutil.ReadDir(vari.AgentLogDir + dirName)

		for _, fi := range files {
			name := fi.Name()
			if path.Ext(name) != ".zip" {
				continue
			}

			item := map[string]string{"name": dirName + "-" + name}
			ret = append(ret, item)
		}
	}

	return
}
