package zentaoHelper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"

	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
)

func ZentaoLoginErr(errs ...interface{}) (err error) {
	arr := make([]string, 0)

	for _, item := range errs {
		arr = append(arr, fmt.Sprintf("%v", item))
	}

	msg := i118Utils.Sprintf("fail_to_login_with_err", strings.Join(arr, "; "))

	err = errors.New(msg)
	logUtils.Infof(color.RedString(err.Error()))

	return
}

func ZentaoRequestErr(errs ...interface{}) (err error) {
	arr := make([]string, 0)

	for _, item := range errs {
		arr = append(arr, fmt.Sprintf("%v", item))
	}

	msg := strings.Join(arr, ", ")
	if strings.Contains(msg, "403") {
		msg = i118Utils.Sprintf("fail_to_request_zentao_403")
	} else {
		msg = i118Utils.Sprintf("fail_to_request_zentao", msg)
	}
	err = errors.New(msg)
	logUtils.Infof(color.RedString(err.Error()))

	return
}
