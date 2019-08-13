package ui

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	"github.com/easysoft/zentaoatf/src/model"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/date"
	print2 "github.com/easysoft/zentaoatf/src/utils/print"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/jroimartin/gocui"
	"strings"
	"time"
)

func InitImportPage() error {
	DestoryRightPages()

	maxX, _ := vari.Cui.Size()
	slideView, _ := vari.Cui.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + LabelWidth
	urlLabel := NewLabelWidget("urlLabel", left, 1, "ZentaoUrl")
	ViewMap["import"] = append(ViewMap["import"], urlLabel.Name())

	left = right + Space
	right = left + TextWidthFull
	urlInput := NewTextWidget("urlInput", left, 1, TextWidthFull, "client://ztpmp.ngtesting.org")
	ViewMap["import"] = append(ViewMap["import"], urlInput.Name())

	left = slideX + 2
	right = left + LabelWidth
	productLabel := NewLabelWidget("productLabel", left, 4, "ProductId")
	ViewMap["import"] = append(ViewMap["import"], productLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	productInput := NewTextWidget("productInput", left, 4, TextWidthHalf, "1")
	ViewMap["import"] = append(ViewMap["import"], productInput.Name())

	left = right + Space
	right = left + LabelWidth
	taskLabel := NewLabelWidget("taskLabel", left, 4, "or TaskId")
	ViewMap["import"] = append(ViewMap["import"], taskLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	taskInput := NewTextWidget("taskInput", left, 4, TextWidthHalf, "")
	ViewMap["import"] = append(ViewMap["import"], taskInput.Name())

	left = slideX + 2
	right = left + LabelWidth
	languageLabel := NewLabelWidget("languageLabel", left, 7, "Language")
	ViewMap["import"] = append(ViewMap["import"], languageLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	languageInput := NewTextWidget("languageInput", left, 7, TextWidthHalf, "python")
	ViewMap["import"] = append(ViewMap["import"], languageInput.Name())

	left = right + Space
	right = left + LabelWidth
	singleFileLabel := NewLabelWidget("singleFileLabel", left, 7, "SingleFile")
	ViewMap["import"] = append(ViewMap["import"], singleFileLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	singleFileInput := NewRadioWidget("singleFileInput", left, 7, true)
	ViewMap["import"] = append(ViewMap["import"], singleFileInput.Name())

	// zentaoService account and password
	y := 10
	left = slideX + 2
	right = left + LabelWidth
	accountLabel := NewLabelWidget("accountLabel", left, y, "Account")
	ViewMap["import"] = append(ViewMap["import"], accountLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	accountInput := NewTextWidget("accountInput", left, y, TextWidthHalf, "admin")
	ViewMap["import"] = append(ViewMap["import"], accountInput.Name())

	left = right + Space
	right = left + LabelWidth
	passwordLabel := NewLabelWidget("passwordLabel", left, y, "Password")
	ViewMap["import"] = append(ViewMap["import"], passwordLabel.Name())

	left = right + Space
	right = left + TextWidthHalf
	passwordInput := NewTextWidget("passwordInput", left, y, TextWidthHalf, "P2ssw0rd")
	ViewMap["import"] = append(ViewMap["import"], passwordInput.Name())

	// button
	y += 3
	buttonX := (maxX-constant.LeftWidth)/2 + constant.LeftWidth - ButtonWidth
	submitInput := NewButtonWidgetAutoWidth("submitInput", buttonX, 13, "Submit", ImportRequest)
	ViewMap["import"] = append(ViewMap["import"], submitInput.Name())

	keyBindsInput(ViewMap["import"])

	return nil
}

func ImportRequest(g *gocui.Gui, v *gocui.View) error {
	urlView, _ := g.View("urlInput")
	productView, _ := g.View("productInput")
	taskView, _ := g.View("taskInput")
	languageView, _ := g.View("languageInput")
	singleFileView, _ := g.View("singleFileInput")
	accountView, _ := g.View("accountInput")
	passwordView, _ := g.View("passwordInput")

	url := strings.TrimSpace(urlView.ViewBuffer())

	productId := strings.TrimSpace(productView.Buffer())
	taskId := strings.TrimSpace(taskView.Buffer())
	language := strings.TrimSpace(languageView.Buffer())
	account := strings.TrimSpace(accountView.Buffer())
	password := strings.TrimSpace(passwordView.Buffer())
	singleFileStr := strings.TrimSpace(singleFileView.Buffer())
	singleFile := ParseRadioVal(singleFileStr)

	var name string
	params := make(map[string]string)
	if productId != "" {
		params["entityType"] = "product"
		params["entityVal"] = productId

		product := zentaoService.GetProductInfo(url, productId)
		name = product.Name
	} else {
		params["entityType"] = "task"
		params["entityVal"] = taskId

		//taskJson := zentaoService.GetTaskInfo(url, taskId)
		//name, _ = taskJson.Get("name").String()
	}

	url = commonUtils.UpdateUrl(url)
	print2.PrintToCmd(fmt.Sprintf("#atf gen -u %s -t %s -v %s -l %s -s %t -a %s -p %s",
		url, params["entityType"], params["entityVal"], language, singleFile, account, password))

	zentaoService.Login(url, account, password)

	var cases []model.TestCase
	if productId != "" {
		cases = zentaoService.ListCaseByProduct(url, productId)
	} else {
		cases = zentaoService.ListCaseByTask(url, taskId)
	}

	count, err := action.Generate(cases, language, singleFile, account, password)
	if err == nil {
		configUtils.SaveConfig("", url, params["entityType"], params["entityVal"], language, singleFile,
			name, account, password)

		print2.PrintToCmd(fmt.Sprintf("success to generate %d test scripts in '%s' at %s",
			count, constant.ScriptDir, dateUtils.DateTimeStr(time.Now())))
	} else {
		print2.PrintToCmd(err.Error())
	}

	return nil
}

func DestoryImportPage() {
	for _, v := range ViewMap["import"] {
		vari.Cui.DeleteView(v)
		vari.Cui.DeleteKeybindings(v)
	}

	vari.Cui.DeleteKeybinding("", gocui.KeyTab, gocui.ModNone)
}
