package resUtils

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	"github.com/aaronchen2k/deeptest/res"
	"io/ioutil"
	"os"
)

func ReadRes(path string) (ret []byte, err error) {
	isRelease := commonUtils.IsRelease()

	if isRelease {
		ret, err = res.Asset(path)
	} else {
		ret, err = ioutil.ReadFile(path)
	}

	dir, _ := os.Getwd()
	fmt.Printf("isRelease=%t, path=%s, ret=%#v, dir=%s", isRelease, path, ret, dir)

	return
}
