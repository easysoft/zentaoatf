package resUtils

import (
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	"github.com/easysoft/zentaoatf/internal/pkg/lib/common"
	"github.com/easysoft/zentaoatf/res"
	"io/ioutil"
	"os"
	"path/filepath"
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

		ret, err = ioutil.ReadFile(pth)
	}

	return
}
