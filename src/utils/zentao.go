package utils

import (
	"github.com/easysoft/zentaoatf/src/model"
	"strconv"
)

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
