package script

import (
	"strconv"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var runStatisticPage playwright.Page

func RunFailStatistic(t provider.T) {
	t.ID("5487")
	t.AddParentSuite("执行脚本")

	if _, err := runStatisticPage.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	_, err := runStatisticPage.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}
	locator, err := runStatisticPage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
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
	elements, _ := runStatisticPage.QuerySelectorAll(".statistic>>span")
	runTimes, _ := elements[0].InnerText()
	failTimes, _ := elements[2].InnerText()
	runTimesInt, _ := strconv.Atoi(runTimes)
	failTimesInt, _ := strconv.Atoi(failTimes)
	err = runStatisticPage.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	if err != nil {
		t.Errorf("Click run fail: %v", err)
		t.FailNow()
	}
	_, err = runStatisticPage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec result fail: %v", err)
		t.FailNow()
	}
	runStatisticPage.WaitForTimeout(200)
	elements, _ = runStatisticPage.QuerySelectorAll(".statistic>>span")
	runTimes2, _ := elements[0].InnerText()
	failTimes2, _ := elements[2].InnerText()
	runTimes2Int, _ := strconv.Atoi(runTimes2)
	failTimes2Int, _ := strconv.Atoi(failTimes2)
	if runTimes2Int-runTimesInt != 1 || failTimes2Int-failTimesInt != 1 {
		t.Error("statistic error")
		t.FailNow()
	}
}

func RunSuccessStatistic(t provider.T) {
	t.ID("5487")
	t.AddParentSuite("执行脚本")

	if _, err := runStatisticPage.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	_, err := runStatisticPage.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}
	locator, err := runStatisticPage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
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
	scriptLocator, err := locator.Locator("text=3_http_interface_call.php")
	if err != nil {
		t.Errorf("Find 3_http_interface_call.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click()
	if err != nil {
		t.Errorf("Click script fail: %v", err)
		t.FailNow()
	}
	elements, _ := runStatisticPage.QuerySelectorAll(".statistic>>span")
	runTimes, _ := elements[0].InnerText()
	succTimes, _ := elements[1].InnerText()
	runTimesInt, _ := strconv.Atoi(runTimes)
	succTimesInt, _ := strconv.Atoi(succTimes)
	err = runStatisticPage.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	if err != nil {
		t.Errorf("Click run fail: %v", err)
		t.FailNow()
	}
	_, err = runStatisticPage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if err != nil {
		t.Errorf("Wait exec result fail: %v", err)
		t.FailNow()
	}
	runStatisticPage.WaitForTimeout(200)
	elements, _ = runStatisticPage.QuerySelectorAll(".statistic>>span")
	runTimes2, _ := elements[0].InnerText()
	succTimes2, _ := elements[1].InnerText()
	runTimes2Int, _ := strconv.Atoi(runTimes2)
	succTimes2Int, _ := strconv.Atoi(succTimes2)
	if runTimes2Int-runTimesInt != 1 || succTimes2Int-succTimesInt != 1 {
		t.Error("statistic error")
		t.FailNow()
	}
}

func RunBugStatistic(t provider.T) {
	t.ID("5487")
	t.AddParentSuite("执行脚本")

	if _, err := runStatisticPage.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	_, err := runStatisticPage.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		t.FailNow()
	}
	locator, err := runStatisticPage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
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
	runStatisticPage.WaitForTimeout(200)
	elements, _ := runStatisticPage.QuerySelectorAll(".statistic>>span")
	bugTimes, _ := elements[3].InnerText()
	bugTimesInt, _ := strconv.Atoi(bugTimes)
	err = runStatisticPage.Click(".statistic>>span>>nth=3")
	if err != nil {
		t.Errorf("Click bug btn fail: %v", err)
		t.FailNow()
	}
	elements, _ = runStatisticPage.QuerySelectorAll("#bugsModal>>tr")

	bugTimes2Int := len(elements)
	if bugTimes2Int-1 != bugTimesInt {
		t.Error("statistic error")
		t.FailNow()
	}
}

func TestUiRunStatistic(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
		return
	}
	headless := false
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		t.FailNow()
		return
	}
	runStatisticPage, err = browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
		return
	}
	if _, err = runStatisticPage.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
		return
	}
	runner.Run(t, "客户端-确认执行统计成功数据", RunSuccessStatistic)
	runner.Run(t, "客户端-确认执行统计失败数据", RunFailStatistic)
	runner.Run(t, "客户端-确认执行统计bug数据", RunBugStatistic)
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		t.FailNow()
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
		return
	}
}
