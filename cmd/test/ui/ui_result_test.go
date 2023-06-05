package main

import (
	"math"
	"strings"
	"testing"
	"time"

	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	ztfTest "github.com/easysoft/zentaoatf/cmd/test/helper/ztf"
	ztfTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/ztf"
	plwHelper "github.com/easysoft/zentaoatf/cmd/test/ui/helper"
	dateUtils "github.com/easysoft/zentaoatf/pkg/lib/date"
	apiTest "github.com/easysoft/zentaoatf/test/helper/zentao/api"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func Detail(t provider.T) {
	t.ID("5489")
	t.AddParentSuite("测试结果")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTest.SelectSite(webpage)
	ztfTest.ExpandWorspace(webpage)

	ztfTest.RunScript(webpage, "1_string_match.php")

	webpage.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	webpage.WaitForSelector(".result-action .btn:has-text('提交结果到禅道')")

	locator := webpage.Locator(".page-result .single small")
	result := locator.InnerText()
	if result != "通过 0.00%" {
		t.Error("Detail result error")
		t.FailNow()
	}

	locator = webpage.Locator(".result-step-checkpoint code")
	expectVal := locator.InnerText()
	if strings.TrimSpace(expectVal) != "~c:!=2~" {
		t.Error("Detail expect error")
		t.FailNow()
	}

	locator = webpage.Locator(".result-step-checkpoint code>>nth=1")
	actualVal := locator.InnerText()
	if strings.TrimSpace(actualVal) != "2" {
		t.Error("Detail actual error")
		t.FailNow()
	}
}

func SubmitResult(t provider.T) {
	t.ID("5499")
	t.AddParentSuite("测试结果")

	caseInfo := apiTest.GetCaseResult(1)
	lastId := caseInfo["Id"].(int64)

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)

	webpage.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	webpage.Click(".result-action .btn:has-text('提交结果到禅道')")

	webpage.WaitForSelector("#syncToZentaoModal")
	titleInput := webpage.Locator("#syncToZentaoModal>>.form-item:has-text('或输入新测试单标题')>>input")

	titleInput.Fill("单元测试测试单")
	webpage.Click("#syncToZentaoModal>>text=确定")

	webpage.WaitForSelector("#syncToZentaoModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	webpage.Locator(".toast-notification-container:has-text('提交成功')")

	//check zentao
	caseInfo = apiTest.GetCaseResult(1)
	resultTime := dateUtils.TimeStrToTimestamp(caseInfo["Date"].(string))
	t.Require().Greater(caseInfo["Id"].(int64), lastId)
	t.Require().Equal("fail", caseInfo["CaseResult"])
	t.Require().LessOrEqual(math.Abs(float64(resultTime-time.Now().Unix())), float64(10))
}

func SubmitBug(t provider.T) {
	t.ID("5500")
	t.AddParentSuite("测试结果")

	lastId := apiTest.GetLastBugId()

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)

	webpage.Click("#rightPane .result-list-item .list-item-title>>nth=0")
	webpage.Click(".page-result .btn:has-text('提交缺陷到禅道')")
	webpage.WaitForSelector("#submitBugModal")
	webpage.Click("#submitBugModal>>text=确定")

	webpage.WaitForSelector("#submitBugModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "提交成功"})

	newLastId := apiTest.GetLastBugId()
	t.Require().Equal(lastId+1, newLastId)
}

func SubmitBugTwoStep(t provider.T) {
	t.ID("5500")
	t.AddParentSuite("测试结果")

	lastId := apiTest.GetLastBugId()

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)

	ztfTest.RunScript(webpage, "1_string_match.php")
	webpage.Click("#rightPane .result-list-item .list-item-title>>nth=0")

	webpage.Click(".page-result .btn:has-text('提交缺陷到禅道')")

	webpage.WaitForSelector("#submitBugModal")
	webpage.Click("#cbox0")
	webpage.Click("#submitBugModal>>text=确定")

	webpage.WaitForSelector("#submitBugModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "提交成功"})

	newLastId := apiTest.GetLastBugId()
	t.Require().Equal(lastId+1, newLastId)
}

func TestUiResult(t *testing.T) {
	runner.Run(t, "客户端-查看测试结果详情", Detail)
	runner.Run(t, "客户端-提交禅道用例脚本测试结果", SubmitResult)
	runner.Run(t, "客户端-提交禅道失败用例为缺陷", SubmitBug)
	runner.Run(t, "客户端-提交禅道部分失败用例为缺陷", SubmitBugTwoStep)
}
