package main

import (
	"testing"

	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var bugBrowser playwright.Browser

func ScriptBug(t provider.T) {
	t.ID("5747")
	t.AddParentSuite("查看bug列表")

	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	ztfTestHelper.RunScript(webpage, "1_string_match.php")
	ztfTestHelper.SubmitResult(webpage)
	webpage.Click(".tree-node-title:has-text('1_string_match.php')")
	webpage.WaitForTimeout(200)
	webpage.Click(".statistic>>span>>nth=3")
	elements := webpage.QuerySelectorAll("#bugsModal>>tr")

	bugTimesInt := len(elements.ElementHandles)
	if bugTimesInt < 2 {
		t.Error("statistic error")
		t.FailNow()
	}
}

func ScriptsBug(t provider.T) {
	t.ID("5748")
	t.AddParentSuite("查看bug列表")

	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.Click(`[title="批量选择"]`)
	webpage.Click(".tree-node-item:has-text('1_string_match.php')>>.tree-node-check")
	webpage.Click(".tree-node-item:has-text('2_webpage_extract.php')>>.tree-node-check")
	webpage.Click(`[title="禅道BUG"]`)
	elements := webpage.QuerySelectorAll("#bugsModal>>tr")

	bugTimesInt := len(elements.ElementHandles)
	if bugTimesInt < 2 {
		t.Error("statistic error")
		t.FailNow()
	}
}
func TestUiBug(t *testing.T) {
	runner.Run(t, "客户端-查看单个脚本bug列表", ScriptBug)
	runner.Run(t, "客户端-查看选中脚本bug列表", ScriptsBug)
}
