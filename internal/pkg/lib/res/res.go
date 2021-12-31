package resUtils

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	"github.com/aaronchen2k/deeptest/res"
	"io/ioutil"
)

func ReadRes(path string) (ret []byte, err error) {
	isRelease := commonUtils.IsRelease()
	if isRelease {
		ret, err = res.Asset(path)
	} else {
		ret, err = ioutil.ReadFile(path)
	}

	return
}
