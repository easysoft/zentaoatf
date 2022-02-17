package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"os"
)

func Clean() {
	path := commConsts.WorkDir + commConsts.LogDirName + consts.PthSep
	bak := path[:len(path)-1] + "-bak" + consts.PthSep + path[len(path):]

	os.RemoveAll(path)
	os.RemoveAll(bak)

	logUtils.Info(i118Utils.Sprintf("success_to_clean_logs"))
}
