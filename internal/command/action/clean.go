package action

import (
	"os"

	"github.com/easysoft/zentaoatf/pkg/consts"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
)

func Clean() {
	path := commConsts.WorkDir + commConsts.LogDirName + consts.FilePthSep
	os.RemoveAll(path)

	logUtils.Info(i118Utils.Sprintf("success_to_clean_logs"))
}
