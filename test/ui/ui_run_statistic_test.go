package script

import (
	"strconv"
	"testing"

	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var runStatisticPage playwright.Page

func RunFailStatistic(t provider.T) {
	t.ID("5487")
	t.AddParentSuite("执行脚本")

	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.WaitForSelector(".tree-node")
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Click()
	scriptLocator := locator.Locator("text=1_string_match.php")
	scriptLocator.Click()
	elements := webpage.QuerySelectorAll(".statistic>>span")
	runTimes := elements.InnerText(0)
	failTimes := elements.InnerText(2)
	runTimesInt, _ := strconv.Atoi(runTimes)
	failTimesInt, _ := strconv.Atoi(failTimes)
	webpage.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	webpage.WaitForTimeout(200)
	elements = webpage.QuerySelectorAll(".statistic>>span")
	runTimes2 := elements.InnerText(0)
	failTimes2 := elements.InnerText(2)
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

	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.WaitForSelector(".tree-node")
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Click()
	scriptLocator := locator.Locator("text=3_http_interface_call.php")
	scriptLocator.Click()
	elements := webpage.QuerySelectorAll(".statistic>>span")
	runTimes := elements.InnerText(0)
	succTimes := elements.InnerText(1)
	runTimesInt, _ := strconv.Atoi(runTimes)
	succTimesInt, _ := strconv.Atoi(succTimes)
	webpage.Click(".tabs-nav-toolbar>>[title=\"Run\"]")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	webpage.WaitForTimeout(200)
	elements = webpage.QuerySelectorAll(".statistic>>span")
	runTimes2 := elements.InnerText(0)
	succTimes2 := elements.InnerText(1)
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

	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	webpage.WaitForSelector(".tree-node")
	locator := webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	locator.Click()
	scriptLocator := locator.Locator("text=1_string_match.php")
	scriptLocator.Click()
	webpage.WaitForTimeout(200)
	elements := webpage.QuerySelectorAll(".statistic>>span")
	bugTimes := elements.InnerText(3)
	bugTimesInt, _ := strconv.Atoi(bugTimes)
	webpage.Click(".statistic>>span>>nth=3")
	elements = webpage.QuerySelectorAll("#bugsModal>>tr")

	bugTimes2Int := len(elements.ElementHandles)
	if bugTimes2Int-1 != bugTimesInt {
		t.Error("statistic error")
		t.FailNow()
	}
}

func TestUiRunStatistic(t *testing.T) {
	runner.Run(t, "客户端-确认执行统计成功数据", RunSuccessStatistic)
	runner.Run(t, "客户端-确认执行统计失败数据", RunFailStatistic)
	runner.Run(t, "客户端-确认执行统计bug数据", RunBugStatistic)
}
