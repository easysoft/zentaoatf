package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/easysoft/zentaoatf/internal/pkg/consts"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"os"
)

func Clean() {
	path := commConsts.WorkDir + commConsts.LogDirName + consts.FilePthSep
	bak := path[:len(path)-1] + "-bak" + consts.FilePthSep + path[len(path):]

	os.RemoveAll(path)
	os.RemoveAll(bak)

	logUtils.Info(i118Utils.Sprintf("success_to_clean_logs"))
}
