package stringUtils

import (
	"strconv"
	"strings"
)

func UnescapeUnicode(raw []byte) (ret string) {
	temp := strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1)
	ret, _ = strconv.Unquote(temp)

	return
}
