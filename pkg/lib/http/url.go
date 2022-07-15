package httpUtils

import "strings"

func AddSepIfNeeded(url string) (ret string) {
	ret = url
	if strings.LastIndex(ret, "/") != len(ret)-1 {
		ret += "/"
	}
	return
}
