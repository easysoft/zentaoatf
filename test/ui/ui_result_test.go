package main

import (
	"strings"
	"testing"

	ztfTest "github.com/easysoft/zentaoatf/test/helper/ztf"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var resultBrowser playwright.Browser

func Detail(t provider.T) {
	t.ID("5489")
	t.AddParentSuite("测试结果")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	resultBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web resultBrowser: %v", err)
		t.FailNow()
	}
	page, err := resultBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = resultBrowser.Close(); err != nil {
			t.Errorf("The resultBrowser cannot be closed: %v", err)
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

	ztfTest.SelectSite(page)
	ztfTest.ExpandProduct(page)
	ztfTest.RunScript(page, "1_string_match.php")
	err = page.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	if err != nil {
		t.Errorf("Click first result fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".result-action .btn:has-text('提交结果到禅道')")
	if err != nil {
		t.Errorf("Wait result tabpage fail: %v", err)
		t.FailNow()
	}
	locator, err := page.Locator(".page-result .single small")
	if err != nil {
		t.Errorf("Find result percent locator fail: %v", err)
		t.FailNow()
	}
	result, err := locator.InnerText()
	if err != nil {
		t.Errorf("Find result percent fail: %v", err)
		t.FailNow()
	}
	if result != "通过 0.00%" {
		t.Error("Detail result error")
		t.FailNow()
	}
	locator, err = page.Locator(".result-step-checkpoint code")
	if err != nil {
		t.Errorf("Find result expect locator fail: %v", err)
		t.FailNow()
	}
	expectVal, err := locator.InnerText()
	if err != nil {
		t.Errorf("Find result expect fail: %v", err)
		t.FailNow()
	}
	if strings.TrimSpace(expectVal) != "~c:!=2~" {
		t.Error("Detail expect error")
		t.FailNow()
	}
	locator, err = page.Locator(".result-step-checkpoint code>>nth=1")
	if err != nil {
		t.Errorf("Find result actual locator fail: %v", err)
		t.FailNow()
	}
	actualVal, err := locator.InnerText()
	if err != nil {
		t.Errorf("Find result actual fail: %v", err)
		t.FailNow()
	}
	if strings.TrimSpace(actualVal) != "2" {
		t.Error("Detail actual error")
		t.FailNow()
	}
}

func SubmitResult(t provider.T) {
	t.ID("5499")
	t.AddParentSuite("测试结果")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	resultBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web resultBrowser: %v", err)
		t.FailNow()
	}
	page, err := resultBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = resultBrowser.Close(); err != nil {
			t.Errorf("The resultBrowser cannot be closed: %v", err)
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
	err = page.Click(".result-action .btn:has-text('提交结果到禅道')")
	if err != nil {
		t.Errorf("Click re-exec failed case btn fail: %v", err)
	}
	_, err = page.WaitForSelector("#syncToZentaoModal")
	if err != nil {
		t.Errorf("Wait syncToZentaoModal fail: %v", err)
		t.FailNow()
	}
	titleInput, err := page.Locator("#syncToZentaoModal>>.form-item:has-text('或输入新测试单标题')>>input")
	if err != nil {
		t.Errorf("Find title input fail: %v", err)
		t.FailNow()
	}
	err = titleInput.Fill("单元测试测试单")
	if err != nil {
		t.Errorf("Fil title input fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#syncToZentaoModal>>text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#syncToZentaoModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	if err != nil {
		t.Errorf("Wait syncToZentaoModal hide fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "提交成功"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Submit result to zentao fail: %v", err)
		t.FailNow()
	}
}

func SubmitBug(t provider.T) {
	t.ID("5500")
	t.AddParentSuite("测试结果")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	resultBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web resultBrowser: %v", err)
		t.FailNow()
	}
	page, err := resultBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = resultBrowser.Close(); err != nil {
			t.Errorf("The resultBrowser cannot be closed: %v", err)
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
	err = page.Click(".page-result .btn:has-text('提交缺陷到禅道')")
	if err != nil {
		t.Errorf("Click re-exec failed case btn fail: %v", err)
	}
	_, err = page.WaitForSelector("#submitBugModal")
	if err != nil {
		t.Errorf("Wait submitBugModal fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#submitBugModal>>text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#submitBugModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	if err != nil {
		t.Errorf("Wait submitBugModal hide fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "提交成功"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Submit bug to zentao fail: %v", err)
		t.FailNow()
	}
}

func SubmitBugTwoStep(t provider.T) {
	t.ID("5500")
	t.AddParentSuite("测试结果")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	resultBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web resultBrowser: %v", err)
		t.FailNow()
	}
	page, err := resultBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	defer func() {
		if err = resultBrowser.Close(); err != nil {
			t.Errorf("The resultBrowser cannot be closed: %v", err)
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
	ztfTest.RunScript(page, "1_string_match.php")
	err = page.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	if err != nil {
		t.Errorf("Click first result fail: %v", err)
	}
	err = page.Click(".page-result .btn:has-text('提交缺陷到禅道')")
	if err != nil {
		t.Errorf("Click re-exec failed case btn fail: %v", err)
	}
	_, err = page.WaitForSelector("#submitBugModal")
	if err != nil {
		t.Errorf("Wait submitBugModal fail: %v", err)
		t.FailNow()
	}
	page.Click("#cbox0")
	err = page.Click("#submitBugModal>>text=确定")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#submitBugModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	if err != nil {
		t.Errorf("Wait submitBugModal hide fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "提交成功"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Submit bug to zentao fail: %v", err)
		t.FailNow()
	}
}

func TestUiResult(t *testing.T) {
	runner.Run(t, "客户端-查看测试结果详情", Detail)
	runner.Run(t, "客户端-提交禅道用例脚本测试结果", SubmitResult)
	runner.Run(t, "客户端-提交禅道失败用例为缺陷", SubmitBug)
	runner.Run(t, "客户端-提交禅道部分失败用例为缺陷", SubmitBugTwoStep)
}
