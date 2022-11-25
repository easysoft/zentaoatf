package main

import (
	"strconv"
	"testing"

	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
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
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	scriptLocator := webpage.Locator(".tree-node-title>>text=1_string_match.php")
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
	if runTimes2Int-runTimesInt != 1 {
		webpage.ScreenShot()
		t.Errorf("statistic fail num error, total num expect: %v ,actual： %v", runTimesInt+1, runTimes2Int)
		t.FailNow()
	}
	if failTimes2Int-failTimesInt != 1 {
		webpage.ScreenShot()
		t.Errorf("statistic fail num error, fail num expect: %v ,actual： %v", failTimesInt+1, failTimes2Int)
		t.FailNow()
	}
}

func RunSuccessStatistic(t provider.T) {
	t.ID("5487")
	t.AddParentSuite("执行脚本")

	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	scriptLocator := webpage.Locator(".tree-node-title>>text=3_http_interface_call.php")
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
	if runTimes2Int-runTimesInt != 1 {
		webpage.ScreenShot()
		t.Errorf("statistic success num error, total num expect: %v ,actual： %v", runTimesInt+1, runTimes2Int)
		t.FailNow()
	}

	if succTimes2Int-succTimesInt != 1 {
		webpage.ScreenShot()
		t.Errorf("statistic success num error, success num expect: %v ,actual： %v", succTimesInt+1, succTimes2Int)
		t.FailNow()
	}
}

func RunBugStatistic(t provider.T) {
	t.ID("5487")
	t.AddParentSuite("执行脚本")

	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	scriptLocator := webpage.Locator(".tree-node-title>>text=1_string_match.php")
	scriptLocator.Click()
	webpage.WaitForTimeout(2000)
	elements := webpage.QuerySelectorAll(".statistic>>span")
	bugTimes := elements.InnerText(3)
	bugTimesInt, _ := strconv.Atoi(bugTimes)
	webpage.Click(".statistic>>span>>nth=3")
	elements = webpage.QuerySelectorAll("#bugsModal>>tr")

	bugTimes2Int := len(elements.ElementHandles)
	if bugTimes2Int-1 != bugTimesInt {
		webpage.ScreenShot()
		t.Errorf("statistic bug num error, bug num expect: %v ,actual： %v", bugTimesInt+1, bugTimes2Int)
		t.FailNow()
	}
}

func TestUiRunStatistic(t *testing.T) {
	runner.Run(t, "客户端-确认执行统计成功数据", RunSuccessStatistic)
	runner.Run(t, "客户端-确认执行统计失败数据", RunFailStatistic)
	runner.Run(t, "客户端-确认执行统计bug数据", RunBugStatistic)
}
