package page

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/ui"
	"github.com/easysoft/zentaoatf/src/ui/widget"
	"github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	"github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"strings"
)

func InitImportPage() error {
	DestoryRightPages()

	//conf := configUtils.ReadCurrConfig()

	productId := ""
	taskId := ""
	//if conf.LangType == "task" {
	//	taskId = conf.EntityVal
	//} else {
	//	productId = conf.EntityVal
	//}

	maxX, _ := vari.Cui.Size()
	slideView, _ := vari.Cui.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + widget.LabelWidth
	urlLabel := widget.NewLabelWidget("urlLabel", left, 1, i118Utils.I118Prt.Sprintf("zentao_url"))
	ui.ViewMap["import"] = append(ui.ViewMap["import"], urlLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthFull
	urlInput := widget.NewTextWidget("urlInput", left, 1, widget.TextWidthFull, "conf.Url")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], urlInput.Name())

	left = slideX + 2
	right = left + widget.LabelWidth
	productLabel := widget.NewLabelWidget("productLabel", left, 4, i118Utils.I118Prt.Sprintf("product_id"))
	ui.ViewMap["import"] = append(ui.ViewMap["import"], productLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	productInput := widget.NewTextWidget("productInput", left, 4, widget.TextWidthHalf, productId)
	ui.ViewMap["import"] = append(ui.ViewMap["import"], productInput.Name())

	left = right + ui.Space
	right = left + widget.LabelWidth
	taskLabel := widget.NewLabelWidget("taskLabel", left, 4,
		i118Utils.I118Prt.Sprintf("task_id")+""+i118Utils.I118Prt.Sprintf("task_id"))
	ui.ViewMap["import"] = append(ui.ViewMap["import"], taskLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	taskInput := widget.NewTextWidget("taskInput", left, 4, widget.TextWidthHalf, taskId)
	ui.ViewMap["import"] = append(ui.ViewMap["import"], taskInput.Name())

	left = slideX + 2
	right = left + widget.LabelWidth
	languageLabel := widget.NewLabelWidget("languageLabel", left, 7, i118Utils.I118Prt.Sprintf("language"))
	ui.ViewMap["import"] = append(ui.ViewMap["import"], languageLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	languageInput := widget.NewTextWidget("languageInput", left, 7, widget.TextWidthHalf, "conf.LangType")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], languageInput.Name())

	left = right + ui.Space
	right = left + widget.LabelWidth
	independentFileLabel := widget.NewLabelWidget("independentFileLabel", left, 7,
		i118Utils.I118Prt.Sprintf("independent_expect_result"))
	ui.ViewMap["import"] = append(ui.ViewMap["import"], independentFileLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	independentFileInput := widget.NewRadioWidget("independentFileInput", left, 7, true)
	ui.ViewMap["import"] = append(ui.ViewMap["import"], independentFileInput.Name())

	// zentaoService account and password
	y := 10
	left = slideX + 2
	right = left + widget.LabelWidth
	accountLabel := widget.NewLabelWidget("accountLabel", left, y, i118Utils.I118Prt.Sprintf("account"))
	ui.ViewMap["import"] = append(ui.ViewMap["import"], accountLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	accountInput := widget.NewTextWidget("accountInput", left, y, widget.TextWidthHalf, "conf.Account")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], accountInput.Name())

	left = right + ui.Space
	right = left + widget.LabelWidth
	passwordLabel := widget.NewLabelWidget("passwordLabel", left, y, i118Utils.I118Prt.Sprintf("password"))
	ui.ViewMap["import"] = append(ui.ViewMap["import"], passwordLabel.Name())

	left = right + ui.Space
	right = left + widget.TextWidthHalf
	passwordInput := widget.NewTextWidget("passwordInput", left, y, widget.TextWidthHalf, "conf.Password")
	ui.ViewMap["import"] = append(ui.ViewMap["import"], passwordInput.Name())

	// button
	y += 3
	buttonX := (maxX-constant.LeftWidth)/2 + constant.LeftWidth - widget.ButtonWidth
	submitInput := widget.NewButtonWidgetAutoWidth("submitInput", buttonX, 13,
		i118Utils.I118Prt.Sprintf("submit"), ImportRequest)
	ui.ViewMap["import"] = append(ui.ViewMap["import"], submitInput.Name())

	ui.BindEventForInputWidgets(ui.ViewMap["import"])

	return nil
}

func ImportRequest(g *gocui.Gui, v *gocui.View) error {
	urlView, _ := g.View("urlInput")
	productView, _ := g.View("productInput")
	taskView, _ := g.View("taskInput")
	languageView, _ := g.View("languageInput")
	independentFileView, _ := g.View("independentFileInput")
	accountView, _ := g.View("accountInput")
	passwordView, _ := g.View("passwordInput")

	url := strings.TrimSpace(urlView.ViewBuffer())

	productId := strings.TrimSpace(productView.Buffer())
	taskId := strings.TrimSpace(taskView.Buffer())
	language := strings.TrimSpace(languageView.Buffer())
	account := strings.TrimSpace(accountView.Buffer())
	password := strings.TrimSpace(passwordView.Buffer())
	independentFileStr := strings.TrimSpace(independentFileView.Buffer())
	independentFile := widget.ParseRadioVal(independentFileStr)

	var entityType string
	var entityVal string
	if productId != "" {
		entityType = "product"
		entityVal = productId
	} else {
		entityType = "task"
		entityVal = taskId
	}

	url = commonUtils.UpdateUrl(url)
	logUtils.PrintToCmd(fmt.Sprintf("#atf gen -u %s -t %s -v %s -l %s -i %t -a %s -p %s",
		url, entityType, entityVal, language, independentFile, account, password))

	//action.GenerateScript(url, entityType, entityVal, language, independentFile, account, password)

	return nil
}

func DestoryImportPage() {
	for _, v := range ui.ViewMap["import"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}

	vari.Cui.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
}
