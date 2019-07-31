package utils

import "github.com/jroimartin/gocui"

const (
	PreferenceFile = "preferences.yaml"
	ConfigFile     = "conf.yaml"

	SuiteExt string = "suite"

	LanguageDefault = "en"
	LanguageEN      = "en"
	LanguageZH      = "zh"

	EnRes = "res/messages_en.json"
	ZhRes = "res/messages_zh.json"

	GenDir = "scripts/"

	LeftWidth = 36
	MinWidth  = 130
	MinHeight = 36

	CmdViewHeight = 10
)

var RunFromCui bool
var Cui *gocui.Gui
var MainViewHeight int
