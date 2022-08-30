package zentaoHelper

import (
	"errors"
	"fmt"
	"strings"

	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
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

	msg := i118Utils.Sprintf("fail_to_request_zentao", strings.Join(arr, ", "))

	err = errors.New(msg)
	logUtils.Infof(color.RedString(err.Error()))

	return
}
