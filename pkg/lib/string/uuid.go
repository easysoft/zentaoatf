package stringUtils

import (
	"strings"

	uuid "github.com/satori/go.uuid"
)

func NewUuid() string {
	uid := uuid.NewV4().String()
	return strings.Replace(uid, "-", "", -1)
}

func NewUuidWithSep() string {
	uid := uuid.NewV4().String()
	return uid
}
