package zentaoUtils

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
)

func GenApiUri(module string, methd string, param string) string {
	var uri string

	if commConsts.RequestType == commConsts.PathInfo {
		uri = fmt.Sprintf("%s-%s-%s.json", module, methd, param)
	} else {
		uri = fmt.Sprintf("index.php?m=%s&f=%s&%s&t=json", module, methd, param)
	}

	return uri
}
