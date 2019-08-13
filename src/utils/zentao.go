package utils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/model"
	"strconv"
)

func GenSuperApiUri(module string, methd string, params [][]string) string {
	var sep string
	if RequestType == RequestTypePathInfo {
		sep = ","
	} else {
		sep = "&"
	}

	paramStr := ""
	i := 0
	for _, p := range params {
		if i > 0 {
			paramStr += sep
		}
		paramStr += p[0] + "=" + p[1]
		i++
	}

	var uri string
	if RequestType == RequestTypePathInfo {
		uri = fmt.Sprintf("api-getmodel-%s-%s-%s.json", module, methd, paramStr)
	} else {
		uri = fmt.Sprintf("?m=api&f=getmodel&module=%s&methodName=%s&params=%s", module, methd, paramStr)
	}

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
