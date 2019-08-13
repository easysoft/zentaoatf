package vari

import (
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/jroimartin/gocui"
)

var (
	Prefer         model.Preference
	ZendaoSettings model.ZentaoSettings
	RunMode        constant.RunMode
	RunDir         string
	RunFromCui     bool
	Cui            *gocui.Gui
	MainViewHeight int

	SessionVar  string
	SessionId   string
	RequestType string
	RequestFix  string
)
