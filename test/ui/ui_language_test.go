package main

import (
	"testing"

	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

func SwitchLanguage(t provider.T) {
	t.ID("5464")
	t.AddParentSuite("设置界面语言")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.Click(`input[type="radio"]>>nth=1`)
	locator := webpage.Locator(".t-card-toolbar div>>nth=2")
	interpreterTitle := locator.InnerText()
	if interpreterTitle != "Remote Server" {
		t.Error("Switch language fail, find interpreter fail")
		t.FailNow()
	}
	locator = webpage.Locator(".t-card-toolbar button")
	CreateInterpreterTitle := locator.InnerText()
	if CreateInterpreterTitle != "Create Remote Server" {
		t.Error("Switch language fail, find create remote server btn fail")
		t.FailNow()
	}
	locator = webpage.Locator("#settingModal .modal-title")
	modalTitle := locator.InnerText()
	if modalTitle != "Settings" {
		t.Error("Switch language fail, find modalTitle fail")
		t.FailNow()
	}
	webpage.Click(`input[type="radio"]>>nth=0`)
}

func TestUiLanguage(t *testing.T) {
	runner.Run(t, "设置界面语言", SwitchLanguage)
}
