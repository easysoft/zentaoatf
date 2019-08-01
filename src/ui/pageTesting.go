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

func InitTestingPage(g *gocui.Gui) error {
	caseFiles, suitesFiles := loadTestAssets()
	dir := utils.Prefer.WorkDir + utils.GenDir

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
	utils.PrintToSide(g, content)

	setViewScroll(g, "side")
	setViewLineSelected(g, "side", selectLineEvent)

	return nil
}

func selectLineEvent(g *gocui.Gui, v *gocui.View) error {
	var line string
	var err error

	_, cy := v.Cursor()
	if line, err = v.Line(cy); err != nil {
		return nil
	}
	line = strings.TrimSpace(line)
	//utils.PrintToCmd(g, line)
	if strings.Index(line, ".") < 0 {
		utils.PrintToMainNoScroll(g, "")
		return nil
	}

	showAsset(g, line)

	return nil
}
func showAsset(g *gocui.Gui, file string) {
	HideHelp(g)
	CurrAsset = utils.Prefer.WorkDir + utils.GenDir + file

	showRunButton(g)
	content := utils.ReadFile(CurrAsset)
	utils.PrintToMainNoScroll(g, content)
}

func showRunButton(g *gocui.Gui) error {
	maxX, _ := g.Size()

	runButton := NewButtonWidgetAutoWidth(g, "runButton", maxX-10, 1, "Run", run)
	ViewMap["testing"] = append(ViewMap["testing"], runButton.Name())

	return nil
}

func run(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("main"); err != nil {
		return err
	}

	utils.PrintToCmd(g, fmt.Sprintf("#atf run -d %s -f %s", utils.Prefer.WorkDir, CurrAsset))
	utils.PrintToMain(g, "")
	action.Run(utils.Prefer.WorkDir, []string{CurrAsset}, "")

	return nil
}

func loadTestAssets() ([]string, []string) {
	config := utils.ReadCurrConfig()
	ext := script.GetLangMap()[config.LangType]["extName"]

	caseFiles, _ := utils.GetAllFiles(utils.Prefer.WorkDir+utils.GenDir, ext)
	suitesFiles, _ := utils.GetAllFiles(utils.Prefer.WorkDir+utils.GenDir, "suite")

	return caseFiles, suitesFiles
}

func init() {

}

func DestoryTestingPage(g *gocui.Gui) {
	g.DeleteKeybindings("side")
	for _, v := range ViewMap["testing"] {
		g.DeleteView(v)
		g.DeleteKeybindings(v)
	}
}
