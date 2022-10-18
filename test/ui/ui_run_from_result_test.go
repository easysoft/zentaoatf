package main

import (
	"strings"
	"testing"

	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

func RunReExecFailCase(t provider.T) {
	t.ID("5491")
	t.AddParentSuite("测试结果页面执行脚本")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	webpage.Click(".result-action .btn:has-text('重新执行失败用例')")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	locator := webpage.Locator("#log-list>>code:has-text('执行1个用例，耗时')")
	innerText := locator.InnerText()
	if !strings.Contains(innerText, "0(0.0%) 通过，1(100.0%) 失败") {
		t.Errorf("Exec failed case fail")
		t.FailNow()
	}
	webpage.WaitForTimeout(2000)
	resultTitle := webpage.InnerText("#rightPane .result-list-item .list-item-title")
	if resultTitle != "1_string_match.php" {
		t.Errorf("Find result in rightPane fail")
		t.FailNow()
	}
	timeElement := locator.Locator(".time>>span")
	if resultTitle != "1_string_match.php" {
		t.Errorf("Find log time element in logPane fail")
		t.FailNow()
	}
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail")
		t.FailNow()
	}
}

func RunReExecAllCase(t provider.T) {
	t.ID("5750")
	t.AddParentSuite("测试结果页面执行脚本")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.ExpandWorspace(webpage)
	ztfTestHelper.RunScript(webpage, "1_string_match.php")
	webpage.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	webpage.Click(".result-action .btn:has-text('重新执行所有用例')")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	locator := webpage.Locator("#log-list>>code:has-text('执行1个用例，耗时')")
	innerText := locator.InnerText()
	if !strings.Contains(innerText, "0(0.0%) 通过，1(100.0%) 失败") {
		t.Errorf("Exec failed case fail")
		t.FailNow()
	}
	webpage.WaitForTimeout(2000)
	resultTitle := webpage.InnerText("#rightPane .result-list-item .list-item-title")
	if resultTitle != "1_string_match.php" {
		t.Errorf("Find result title in rightPane fail")
		t.FailNow()
	}
	timeElement := locator.Locator(".time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result time in rightPane fail")
		t.FailNow()
	}
}
func TestUiRunFromResult(t *testing.T) {
	runner.Run(t, "客户端-结果中重新执行所有脚本", RunReExecAllCase)
	runner.Run(t, "客户端-结果中重新执行失败脚本", RunReExecFailCase)
}
