package main

import (
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/test/helper/common"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func RunScript(t provider.T) {
	t.ID("5479")
	t.AddParentSuite("执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator("text=1_string_match.php")
	if err != nil {
		t.Errorf("Find 1_string_match.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click script fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	if err != nil {
		t.Errorf("Click run fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec result fail: %v", err)
		t.FailNow()
	}
	element, err := page.QuerySelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	innerText, err := element.InnerText()
	if err != nil {
		t.Errorf("Find result fail: %v", err)
		t.FailNow()
	}
	if !strings.Contains(innerText, "1(100.0%) 失败") {
		t.Errorf("Exec 1_string_match.php fail: %v", err)
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
	timeElement, err := page.QuerySelector("#log-list .item .time")
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
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

	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunSelectedScripts(t provider.T) {
	t.ID("5481")
	t.AddParentSuite("执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		t.FailNow()
	}
	err = page.Click(`[title="批量选择"]`)
	if err != nil {
		t.Errorf("The Click select btn fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator(".tree-node-item:has-text('1_string_match.php')")
	if err != nil {
		t.Errorf("Find 1_string_match.php fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err = scriptLocator.Locator(".tree-node-check")
	if err != nil {
		t.Errorf("Find 1_string_match.php checkbox fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click 1_string_match.php checkbox fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err = locator.Locator(".tree-node-item:has-text('2_webpage_extract.php')")
	if err != nil {
		t.Errorf("Find 2_webpage_extract.php fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err = scriptLocator.Locator(".tree-node-check")
	if err != nil {
		t.Errorf("Find 2_webpage_extract.php checkbox fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click 2_webpage_extract.php checkbox fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".run-selected")
	if err != nil {
		t.Errorf("Click run fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行2个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec result fail: %v", err)
		t.FailNow()
	}
	element, err := page.QuerySelector("#log-list>>.msg-span>>:has-text('执行2个用例，耗时')")
	innerText, err := element.InnerText()
	if err != nil {
		t.Errorf("Find result fail: %v", err)
		t.FailNow()
	}
	if !strings.Contains(innerText, "1(50.0%) 通过，1(50.0%) 失败") {
		t.Errorf("Exec 1_string_match.php,2_webpage_extract.php fail: %v", err)
		t.FailNow()
	}
	resultTitleElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-title")
	if err != nil {
		t.Errorf("Find log title in logPane fail: %v", err)
		t.FailNow()
	}
	resultTitle, err := resultTitleElement.InnerText()
	if err != nil || resultTitle != "单元测试工作目录(2)" {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}
	timeElement, err := page.QuerySelector("#log-list .item .time")
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
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
	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunOpenedAndLast(t provider.T) {
	t.ID("5484")
	t.AddParentSuite("执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	locator, err = page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil || c == 0 {
		t.Errorf("Click workspace fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-node-item:has-text('1_string_match.php')")
	if err != nil {
		t.Errorf("Click 1_string_match.php fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-node-item:has-text('2_webpage_extract.php')")
	if err != nil {
		t.Errorf("Click 2_webpage_extract.php fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#batchRunMenuToggle")
	if err != nil {
		t.Errorf("Click batchRunMenuToggle fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".list-item-content:has-text('执行打开文件')")
	if err != nil {
		t.Errorf("The Click Run opened scripts fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行2个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec opened scripts result fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#log-list>>code:has-text('执行2个用例，耗时')")
	if err != nil {
		t.Errorf("Find exec opened scripts result fail: %v", err)
		t.FailNow()
	}
	innerText, err := locator.InnerText()
	if err != nil {
		t.Errorf("Find exec opened scripts result fail: %v", err)
		t.FailNow()
	}
	if !strings.Contains(innerText, "1(50.0%) 通过，1(50.0%) 失败") {
		t.Errorf("Exec opened scripts fail: %v", err)
		t.FailNow()
	}

	resultTitleElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-title")
	if err != nil {
		t.Errorf("Find log title in logPane fail: %v", err)
	}
	resultTitle, err := resultTitleElement.InnerText()
	if err != nil || resultTitle != "单元测试工作目录(3)" {
		t.Errorf("Find result in rightPane fail: %v", err)
	}
	timeElement, err := locator.Locator(".time>>span")
	if err != nil {
		t.Errorf("Find log time element in logPane fail: %v", err)
		t.FailNow()
	}
	logTime, err := timeElement.InnerText()
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
	}
	resultTimeElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-trailing-text")
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
	}
	resultTime, err := resultTimeElement.InnerText()
	if err != nil || logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail: %v", err)
	}

	err = page.Click("#batchRunMenuToggle")
	if err != nil {
		t.Errorf("Click batchRunMenuToggle fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".list-item-content:has-text('执行上次')")
	if err != nil {
		t.Errorf("The Click Run last time fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行2个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec last time result fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#log-list>>code:has-text('执行2个用例，耗时')")
	innerText, err = locator.InnerText()
	if err != nil {
		t.Errorf("Find exec last time result fail: %v", err)
		t.FailNow()
	}
	if !strings.Contains(innerText, "1(50.0%) 通过，1(50.0%) 失败") {
		t.Errorf("Exec last time fail: %v", err)
		t.FailNow()
	}
	resultTitleElement, err = page.QuerySelector("#rightPane .result-list-item .list-item-title")
	if err != nil {
		t.Errorf("Find log title in logPane fail: %v", err)
		t.FailNow()
	}
	resultTitle, err = resultTitleElement.InnerText()
	if err != nil || resultTitle != "单元测试工作目录(2)" {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}
	timeElement, err = locator.Locator(".time>>span")
	if err != nil {
		t.Errorf("Find log time element in logPane fail: %v", err)
		t.FailNow()
	}
	logTime, err = timeElement.InnerText()
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
		t.FailNow()
	}
	resultTimeElement, err = page.QuerySelector("#rightPane .result-list-item .list-item-trailing-text")
	if err != nil {
		t.Errorf("Find log time in logPane fail: %v", err)
		t.FailNow()
	}
	resultTime, err = resultTimeElement.InnerText()
	if err != nil || logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}
	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunAll(t provider.T) {
	t.ID("5485")
	t.AddParentSuite("执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	locator, err = page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}

	err = page.Click("#batchRunMenuToggle")
	if err != nil {
		t.Errorf("Click batchRunMenuToggle fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".list-item-content:has-text('执行所有文件')")
	if err != nil {
		t.Errorf("The Click Run all scripts fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行3个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec all scripts result fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#log-list>>code:has-text('执行3个用例，耗时')")
	if err != nil {
		t.Errorf("Find exec all scripts result fail: %v", err)
		t.FailNow()
	}
	innerText, err := locator.InnerText()
	if err != nil {
		t.Errorf("Find exec all scripts result fail: %v", err)
		t.FailNow()
	}
	if !strings.Contains(innerText, "2(66.0%) 通过，1(33.0%) 失败") {
		t.Errorf("Exec all fail: %v", err)
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
	}
	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunReExecFailCase(t provider.T) {
	t.ID("5491")
	t.AddParentSuite("执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunReExecAllCase(t provider.T) {
	t.ID("5491")
	t.AddParentSuite("执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunFailStatistic(t provider.T) {
	t.ID("5487")
	t.AddParentSuite("执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator("text=1_string_match.php")
	if err != nil {
		t.Errorf("Find 1_string_match.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click script fail: %v", err)
		t.FailNow()
	}
	elements, _ := page.QuerySelectorAll(".statistic>>span")
	runTimes, _ := elements[0].InnerText()
	failTimes, _ := elements[2].InnerText()
	runTimesInt, _ := strconv.Atoi(runTimes)
	failTimesInt, _ := strconv.Atoi(failTimes)
	err = page.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	if err != nil {
		t.Errorf("Click run fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec result fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(200)
	elements, _ = page.QuerySelectorAll(".statistic>>span")
	runTimes2, _ := elements[0].InnerText()
	failTimes2, _ := elements[2].InnerText()
	runTimes2Int, _ := strconv.Atoi(runTimes2)
	failTimes2Int, _ := strconv.Atoi(failTimes2)
	if runTimes2Int-runTimesInt != 1 || failTimes2Int-failTimesInt != 1 {
		t.Error("statistic error")
		t.FailNow()
	}
	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunSuccessStatistic(t provider.T) {
	t.ID("5487")
	t.AddParentSuite("执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator("text=2_webpage_extract.php")
	if err != nil {
		t.Errorf("Find 2_webpage_extract.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click script fail: %v", err)
		t.FailNow()
	}
	elements, _ := page.QuerySelectorAll(".statistic>>span")
	runTimes, _ := elements[0].InnerText()
	succTimes, _ := elements[1].InnerText()
	runTimesInt, _ := strconv.Atoi(runTimes)
	succTimesInt, _ := strconv.Atoi(succTimes)
	err = page.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	if err != nil {
		t.Errorf("Click run fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec result fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(200)
	elements, _ = page.QuerySelectorAll(".statistic>>span")
	runTimes2, _ := elements[0].InnerText()
	succTimes2, _ := elements[1].InnerText()
	runTimes2Int, _ := strconv.Atoi(runTimes2)
	succTimes2Int, _ := strconv.Atoi(succTimes2)
	if runTimes2Int-runTimesInt != 1 || succTimes2Int-succTimesInt != 1 {
		t.Error("statistic error")
		t.FailNow()
	}
	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunBugStatistic(t provider.T) {
	t.ID("5487")
	t.AddParentSuite("执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator("text=1_string_match.php")
	if err != nil {
		t.Errorf("Find 1_string_match.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click script fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(200)
	elements, _ := page.QuerySelectorAll(".statistic>>span")
	bugTimes, _ := elements[3].InnerText()
	bugTimesInt, _ := strconv.Atoi(bugTimes)
	err = page.Click(".statistic>>span>>nth=3")
	if err != nil {
		t.Errorf("Click bug btn fail: %v", err)
		t.FailNow()
	}
	elements, _ = page.QuerySelectorAll("#bugsModal>>tr")

	bugTimes2Int := len(elements)
	if bugTimes2Int-1 != bugTimesInt {
		t.Error("statistic error")
		t.FailNow()
	}
	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunWorkspace(t provider.T) {
	t.ID("5482")
	t.AddParentSuite("右键执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
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
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click workspace fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-context-menu>>text=执行")
	if err != nil {
		t.Errorf("Click copy fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行3个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec workspace result fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#log-list>>code:has-text('执行3个用例，耗时')")
	if err != nil {
		t.Errorf("Find exec workspace log fail: %v", err)
		t.FailNow()
	}
	innerText, err := locator.InnerText()
	if err != nil {
		t.Errorf("Find exec workspace result fail: %v", err)
		t.FailNow()
	}
	if !strings.Contains(innerText, "2(66.0%) 通过，1(33.0%) 失败") {
		t.Errorf("Exec workspace fail: %v", err)
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
	timeElement, err := page.Locator("#log-list>>.case-item:has-text('3_http_interface_call')>>.time>>span")
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
		t.Errorf("Find log time in rightPane fail: %v", err)
		t.FailNow()
	}
	resultTime, err := resultTimeElement.InnerText()
	if err != nil || logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}

	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func RunUnit(t provider.T) {
	var pwd, err = os.Getwd()
	testngDir := pwd + "/demo/ci_test_testng"
	if runtime.GOOS == "windows" {
		testngDir = pwd + "\\demo\\ci_test_testng"
	}
	commonTestHelper.CloneGit("https://gitee.com/ngtesting/ci_test_testng.git", testngDir)
	t.ID("5432")
	t.Title("执行TestNG单元测试")
	t.ID("5482")
	t.AddParentSuite("右键执行脚本")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web runBrowser: %v", err)
		t.FailNow()
	}
	page, err := runBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
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
		t.Errorf("Wait for workspace list nav fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".list-item-title>>text=单元测试站点")
	if err != nil {
		t.Errorf("The Click workspace nav fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "testng工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		createWorkspace(t, testngDir, page)
	}
	locator, err = page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "testng工作目录"})
	c, err = locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click workspace fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-context-menu>>text=执行")
	if err != nil {
		t.Errorf("Click right run btn fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#tabsPane >> text=执行")
	if err != nil {
		t.Errorf("Click run btn fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#log-list>>.msg-span>>:has-text('执行3个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec testng result fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#log-list>>code:has-text('执行3个用例，耗时')")
	if err != nil {
		t.Errorf("Find exec testng log fail: %v", err)
		t.FailNow()
	}
	innerText, err := locator.InnerText()
	if err != nil {
		t.Errorf("Find exec testng result fail: %v", err)
		t.FailNow()
	}
	if !strings.Contains(innerText, "3(100.0%) 通过，0(0.0%) 失败") {
		t.Errorf("Exec testng fail: %v", err)
		t.FailNow()
	}
	resultTitleElement, err := page.QuerySelector("#rightPane .result-list-item .list-item-title")
	if err != nil {
		t.Errorf("Find log title in logPane fail: %v", err)
		t.FailNow()
	}
	resultTitle, err := resultTitleElement.InnerText()
	if err != nil || resultTitle != "testng工作目录(3)" {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}
	timeElement, err := page.Locator("#log-list>>code:has-text('执行3个用例，耗时')>>.time>>span")
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
		t.Errorf("Find log time in rightPane fail: %v", err)
		t.FailNow()
	}
	resultTime, err := resultTimeElement.InnerText()
	if err != nil || logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail: %v", err)
		t.FailNow()
	}

	if err = runBrowser.Close(); err != nil {
		t.Errorf("The runBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func createWorkspace(t provider.T, workspacePath string, page playwright.Page) {
	err := page.Click(`[title="新建工作目录"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#workspaceFormModal")
	locator, err := page.Locator("#workspaceFormModal input")
	if err != nil {
		t.Errorf("Find create workspace input fail: %v", err)
		t.FailNow()
	}
	titleInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find title input fail: %v", err)
		t.FailNow()
	}
	err = titleInput.Fill("testng工作目录")
	if err != nil {
		t.Errorf("Fill title input fail: %v", err)
		t.FailNow()
	}
	pathInput, err := locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		t.FailNow()
	}
	err = pathInput.Fill(workspacePath)
	if err != nil {
		t.Errorf("Fill address input fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator("#workspaceFormModal select")
	if err != nil {
		t.Errorf("Find create workspace select fail: %v", err)
		t.FailNow()
	}
	typeInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find name input fail: %v", err)
		t.FailNow()
	}
	_, err = typeInput.SelectOption(playwright.SelectOptionValues{Values: &[]string{"testng"}})
	if err != nil {
		t.Errorf("Fil name input fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#workspaceFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: "testng工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created workspace fail: %v", err)
		t.FailNow()
	}
}

func TestUiRun(t *testing.T) {
	//start docker
	// commonTestHelper.Run("12.3.3")
	// runner.Run(t, "客户端-执行单个脚本", RunScript)
	// runner.Run(t, "客户端-执行选中的脚本文件和文件夹", RunSelectedScripts)
	// runner.Run(t, "客户端-执行打开的脚本文件", RunOpenedAndLast)
	// runner.Run(t, "客户端-执行所有的脚本文件", RunAll)
	// runner.Run(t, "客户端-结果中重新执行所有脚本", RunReExecAllCase)
	// runner.Run(t, "客户端-结果中重新执行失败脚本", RunReExecFailCase)
	// runner.Run(t, "客户端-确认执行统计成功数据", RunSuccessStatistic)
	// runner.Run(t, "客户端-确认执行统计失败数据", RunFailStatistic)
	// runner.Run(t, "客户端-确认执行统计bug数据", RunBugStatistic)
	// runner.Run(t, "客户端-右键执行工作目录", RunWorkspace)
	runner.Run(t, "客户端-右键执行工作目录", RunUnit)
}
