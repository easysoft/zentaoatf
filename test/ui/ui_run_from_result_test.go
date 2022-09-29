package main

import (
	"strings"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var runFromResultBrowser playwright.Browser

func RunReExecFailCase(t provider.T) {
	t.ID("5491")
	t.AddParentSuite("测试结果页面执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runFromResultBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runFromResultBrowser: %v", err)
		t.FailNow()
	}
	page, err := runFromResultBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = runFromResultBrowser.Close(); err != nil {
			t.Errorf("The workspaceBrowser cannot be closed: %v", err)
			t.FailNow()
			return
		}
		if err = pw.Stop(); err != nil {
			t.Errorf("The playwright cannot be stopped: %v", err)
			t.FailNow()
			return
		}
	}()
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}

	locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for site list nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".list-item-title>>text=单元测试站点")
	if err != nil {
		t.Errorf("The Click site nav fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	if err != nil {
		t.Errorf("Click first result fail: %v", err)
	}
	err = page.Click(".result-action .btn:has-text('重新执行失败用例')")
	if err != nil {
		t.Errorf("Click re-exec failed case btn fail: %v", err)
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec script result fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#log-list>>code:has-text('执行1个用例，耗时')")
	if err != nil {
		t.Errorf("Find exec script log fail: %v", err)
		t.FailNow()
	}
	innerText, err := locator.InnerText()
	if err != nil {
		t.Errorf("Find exec script result fail: %v", err)
		t.FailNow()
	}
	if !strings.Contains(innerText, "0(0.0%) 通过，1(100.0%) 失败") {
		t.Errorf("Exec failed case fail: %v", err)
		t.FailNow()
	}
	resultTitleElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-title")
	if err != nil {
		t.Errorf("Find log title in logPane fail: %v", err)
		t.FailNow()
	}
	resultTitle, err := resultTitleElement.InnerText()
	if err != nil || resultTitle != "1_string_match.php" {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}
	timeElement, err := locator.Locator(".time>>span")
	if err != nil || resultTitle != "1_string_match.php" {
		t.Errorf("Find log time element in logPane fail: %v", err)
		t.FailNow()
	}
	logTime, err := timeElement.InnerText()
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
		t.FailNow()
	}
	resultTimeElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-trailing-text")
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
		t.FailNow()
	}
	resultTime, err := resultTimeElement.InnerText()
	if err != nil || logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}
}

func RunReExecAllCase(t provider.T) {
	t.ID("5750")
	t.AddParentSuite("测试结果页面执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runFromResultBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runFromResultBrowser: %v", err)
		t.FailNow()
	}
	page, err := runFromResultBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = runFromResultBrowser.Close(); err != nil {
			t.Errorf("The workspaceBrowser cannot be closed: %v", err)
			t.FailNow()
			return
		}
		if err = pw.Stop(); err != nil {
			t.Errorf("The playwright cannot be stopped: %v", err)
			t.FailNow()
			return
		}
	}()
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}

	locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for site list nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".list-item-title>>text=单元测试站点")
	if err != nil {
		t.Errorf("The Click site nav fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	if err != nil {
		t.Errorf("Click first result fail: %v", err)
	}
	err = page.Click(".result-action .btn:has-text('重新执行所有用例')")
	if err != nil {
		t.Errorf("Click re-exec failed case btn fail: %v", err)
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行3个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec script result fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#log-list>>code:has-text('执行3个用例，耗时')")
	if err != nil {
		t.Errorf("Find exec script log fail: %v", err)
		t.FailNow()
	}
	innerText, err := locator.InnerText()
	if err != nil {
		t.Errorf("Find exec script result fail: %v", err)
		t.FailNow()
	}
	if !strings.Contains(innerText, "2(66.0%) 通过，1(33.0%) 失败") {
		t.Errorf("Exec failed case fail: %v", err)
		t.FailNow()
	}
	resultTitleElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-title")
	if err != nil {
		t.Errorf("Find log title in logPane fail: %v", err)
		t.FailNow()
	}
	resultTitle, err := resultTitleElement.InnerText()
	if err != nil || resultTitle != "单元测试工作目录(3)" {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}
	timeElement, err := locator.Locator(".time>>span")
	if err != nil {
		t.Errorf("Find log time element in logPane fail: %v", err)
		t.FailNow()
	}
	logTime, err := timeElement.InnerText()
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
		t.FailNow()
	}
	resultTimeElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-trailing-text")
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
		t.FailNow()
	}
	resultTime, err := resultTimeElement.InnerText()
	if err != nil || logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}
}
func TestUiRunFromResult(t *testing.T) {
	runner.Run(t, "客户端-结果中重新执行所有脚本", RunReExecAllCase)
	runner.Run(t, "客户端-结果中重新执行失败脚本", RunReExecFailCase)
}
