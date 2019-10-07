package action

import (
	"fmt"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
	"strings"
)

func AddPath() {
	ztfPath := strings.TrimRight(vari.ZtfDir, string(os.PathSeparator))
	println(ztfPath)

	pathEnv := os.Getenv("PATH")
	fmt.Println(pathEnv)

	if strings.Index(pathEnv, ztfPath) == -1 {
		if commonUtils.IsWin() {

		} else {
			os.Setenv("PATH", pathEnv+":"+ztfPath)
		}
	}

	pathEnv = os.Getenv("PATH")
	fmt.Println(pathEnv)

	logUtils.PrintToStdOut(i118Utils.I118Prt.Sprintf("success_add_to_path"), -1)
}
