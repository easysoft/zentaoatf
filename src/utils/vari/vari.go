package vari

import (
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/jroimartin/gocui"
)

var (
	Config         model.Config
	Cui            *gocui.Gui
	MainViewHeight int

	RunMode    constant.RunMode
	RunDir     string
	RunFromCui bool

	SessionVar  string
	SessionId   string
	RequestType string
	RequestFix  string

	CurrScriptFile string // scripts/tc-001.py
	CurrResultDate string // 2019-08-15T173802
	CurrCaseId     int    // 2019-08-15T173802

	ReportDir       string
	ScreenWidth     int
	ZentaoBugFileds model.ZentaoBugFileds
)
