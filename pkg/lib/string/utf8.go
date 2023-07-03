package stringUtils

import (
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func Convert2Utf8IfNeeded(data string) string {
	if !utf8.Valid([]byte(data)) && IsGBK([]byte(data)) {
		newLine, _ := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(data))
		data = string(newLine)
	}

	return data
}
