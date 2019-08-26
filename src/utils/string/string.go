package stringUtils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
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

func U2s(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func BoolToPass(b bool) string {
	if b {
		return constant.PASS.String()
	} else {
		return constant.FAIL.String()
	}
}

func StructToStr(obj interface{}) string {
	val, _ := json.Marshal(obj)
	return string(val)
}

func FindInArr(str string, arr []string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}

	return false
}
