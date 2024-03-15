package resUtils

import (
	"os"
	"path/filepath"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	"github.com/easysoft/zentaoatf/res"
)

func ReadRes(path string) (ret []byte, err error) {
	isRelease := commonUtils.IsRelease()

	if isRelease {
		ret, err = res.Asset(path)
	} else {
		pth, _ := filepath.Abs(path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			pth = filepath.Join(commConsts.ZtfDir, path) // in ide, set program args to testng project path
		}

		ret, err = os.ReadFile(pth)
	}

	return
}
