package vari

import (
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/jroimartin/gocui"
)

var (
	Prefer          model.Preference
	ZentaoBugFileds model.ZentaoBugFileds
	RunMode         constant.RunMode
	RunDir          string
	RunFromCui      bool
	Cui             *gocui.Gui
	MainViewHeight  int

	SessionVar  string
	SessionId   string
	RequestType string
	RequestFix  string

	CurrScriptFile string // scripts/tc-001.py
	CurrResultDate string // 2019-08-15T173802
	CurrCaseId     int    // 2019-08-15T173802
)
