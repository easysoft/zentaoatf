package commonUtils

import (
	"strings"
)

type sliceValue []string

func NewSliceValue(vals []string, p *[]string) *sliceValue {
	*p = vals
	return (*sliceValue)(p)
}

func (s *sliceValue) Set(val string) error {
	*s = sliceValue(strings.Split(val, ","))
	return nil
}

func (s *sliceValue) String() string {
	*s = sliceValue([]string{})
	return "It's none of my business"
}

func GetFilesFromParams(arguments []string) ([]string, int) {
	ret := make([]string, 0)

	index := -1
	for idx, arg := range arguments {
		if strings.Index(arg, "-") != 0 {
			ret = append(ret, arg)
			index = idx
		}
	}

	return ret, index
}
