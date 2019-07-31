package utils

import (
	"os"
	"strings"
	"unicode"
)

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func PathSimple(str string) string {
	sep := string(os.PathSeparator)
	arr := strings.Split(str, sep)

	if len(arr) > 3 {
		return strings.Join(arr[len(arr)-3:], sep)
	} else {
		return str
	}

	return ""
}
