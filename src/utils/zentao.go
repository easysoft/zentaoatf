package utils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"strconv"
)

func GenSuperApiUri(module string, methd string, params map[string]string) string {
	paramStr := ""
	i := 0
	for pkey, pval := range params {
		if i > 0 {
			paramStr += "&"
		}
		paramStr += pkey + "=" + pval
	}

	uri := fmt.Sprintf("?m=api&f=getModel&module=%s&methodName=%s&params=%s", module, methd, paramStr)
	fmt.Println(uri)
	return uri
}

func GenApiUri(module string, methd string, param string) string {
	if RequestType == RequestTypePathInfo {
		return fmt.Sprintf("%s-%s-%s.json", module, methd, param)
	}

	return ""
}

func SetBugField(name string, optName string, filedValMap map[string]int) {
	var options []model.Option
	if name == "module" {
		options = ZendaoSettings.Modules
	} else if name == "category" {
		options = ZendaoSettings.Categories
	} else if name == "version" {
		options = ZendaoSettings.Versions
	} else if name == "priority" {
		options = ZendaoSettings.Priorities
	}

	for _, opt := range options {
		if opt.Name == optName {
			filedValMap[name] = opt.Id
		}
	}

	PrintToCmd(strconv.Itoa(filedValMap[name]))
}

func IsBugFieldDefault(optName string, options []model.Option) bool {
	for _, opt := range options {
		if opt.IsDefault && opt.Name == optName {
			return true
		}
	}

	return false
}
