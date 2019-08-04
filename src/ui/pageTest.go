package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/script"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"strings"
)

var CurrAsset string

var tabs []string
var contentViews []string

func InitTestPage() error {
	// left
	caseFiles, suitesFiles := script.LoadTestAssets()
	dir := utils.Prefer.WorkDir + utils.ScriptDir

	content := "Test Suite:" + "\n"
	for _, suitePath := range suitesFiles {
		suiteName := strings.Replace(suitePath, dir, "", -1)
		content += "  " + suiteName + "\n"
	}

	content += "Test Scripts:" + "\n"
	for _, casePath := range caseFiles {
		caseName := strings.Replace(casePath, dir, "", -1)
		content += "  " + caseName + "\n"
	}
	utils.PrintToSide(content)

	// right
	setViewScroll("side")
	setViewLineSelected("side", selectAssetEvent)

	return nil
}

func selectAssetEvent(g *gocui.Gui, v *gocui.View) error {
	clearPanelCaseResult()

	var line string
	var err error

	_, cy := v.Cursor()
	if line, err = v.Line(cy); err != nil {
		return nil
	}
	line = strings.TrimSpace(line)
	if strings.Index(line, ".") < 0 {
		utils.PrintToMainNoScroll("")
		return nil
	}
	CurrAsset = utils.ScriptDir + line

	// show
	if len(tabs) == 0 {
		HideHelp()
		showTab()
	}

	defaultTab, _ := g.View("tabContentView")
	showContent(g, defaultTab)

	return nil
}

func showTab() error {
	g := utils.Cui
	x := utils.LeftWidth + 1
	tabContentView := NewLabelWidgetAutoWidth("tabContentView", x, 0, "Content")
	ViewMap["testing"] = append(ViewMap["testing"], tabContentView.Name())
	tabs = append(tabs, tabContentView.Name())
	if err := g.SetKeybinding("tabContentView", gocui.MouseLeft, gocui.ModNone, showContent); err != nil {
		return nil
	}

	tabResultView := NewLabelWidgetAutoWidth("tabResultView", x+12, 0, "Results")
	ViewMap["testing"] = append(ViewMap["testing"], tabResultView.Name())
	tabs = append(tabs, tabResultView.Name())
	if err := g.SetKeybinding("tabResultView", gocui.MouseLeft, gocui.ModNone, showRun); err != nil {
		return nil
	}

	return nil
}

func showContent(g *gocui.Gui, v *gocui.View) error {
	DestoryRunPanel()
	HighlightTab(v.Name(), tabs)

	panelFileContent, _ := g.View("panelFileContent")
	if panelFileContent != nil {
		panelFileContent.Clear()
	} else {
		maxX, _ := g.Size()
		panelFileContent = NewPanelWidget(utils.CuiRunOutputView, utils.LeftWidth, 2,
			maxX-utils.LeftWidth-1, utils.MainViewHeight, "")
		ViewMap["testing"] = append(ViewMap["testing"], panelFileContent.Name())
		contentViews = append(contentViews, panelFileContent.Name())
		setViewScroll(panelFileContent.Name())

		runButton := NewButtonWidgetAutoWidth("runButton", maxX-10, 0, "[Run]", run)
		runButton.Frame = false
		contentViews = append(contentViews, runButton.Name())
	}

	panelFileContent.Clear()
	panelFileContent.SetOrigin(0, 0)
	content := utils.ReadFile(CurrAsset)
	fmt.Fprintln(panelFileContent, content)

	return nil
}

func init() {

}

func run(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}

	utils.PrintToCmd(fmt.Sprintf("#atf run -d %s -f %s", utils.Prefer.WorkDir, CurrAsset))
	output, _ := g.View(utils.CuiRunOutputView)
	output.Clear()
	action.Run(utils.Prefer.WorkDir, []string{CurrAsset}, "")

	return nil
}

func DestoryTestPage() {
	utils.Cui.DeleteKeybindings("side")
	for _, v := range ViewMap["testing"] {
		utils.Cui.DeleteView(v)
		utils.Cui.DeleteKeybindings(v)
	}
	tabs = []string{}
}

func DestoryContentPanel() {
	for _, v := range contentViews {
		utils.Cui.DeleteView(v)
		utils.Cui.DeleteKeybindings(v)
	}
}
