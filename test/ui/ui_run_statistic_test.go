package main

import (
	"strconv"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

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

func TestUiRunToolbar(t *testing.T) {
	runner.Run(t, "客户端-确认执行统计成功数据", RunSuccessStatistic)
	runner.Run(t, "客户端-确认执行统计失败数据", RunFailStatistic)
	runner.Run(t, "客户端-确认执行统计bug数据", RunBugStatistic)
}
