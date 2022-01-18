package resUtils

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/res"
	"io/ioutil"
	"log"
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

	msg := fmt.Sprintf("isRelease=%t, path=%s, dir=%s", isRelease, path, dir)
	if logUtils.LoggerConsole != nil {
		logUtils.Info(msg)
	} else {
		log.Println(msg)
	}

	return
}
