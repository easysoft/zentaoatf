package main

import (
	"path/filepath"
	"strings"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	zentaoTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/zentao/ui"
	ztfTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/ztf"
	plwConf "github.com/easysoft/zentaoatf/cmd/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/cmd/test/ui/helper"
	commandConfig "github.com/easysoft/zentaoatf/internal/command/config"
	shellUtils "github.com/easysoft/zentaoatf/pkg/lib/shell"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func createUnitWorkspace(t provider.T, workspaceName, workspacePath, unitType string, webpage plwHelper.Webpage) {
	webpage.Click(`[title="新建工作目录"]`)
	webpage.WaitForSelector("#workspaceFormModal")

	locator := webpage.Locator("#workspaceFormModal input")
	locator.FillNth(0, workspaceName)
	locator.FillNth(1, workspacePath)
	locator = webpage.Locator("#workspaceFormModal select")
	locator.SelectNth(0, playwright.SelectOptionValues{Values: &[]string{unitType}})

	webpage.Click("#workspaceFormModal>>.modal-action>>span:has-text(\"确定\")")

	webpage.WaitForSelector("#workspaceFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	webpage.WaitForTimeout(1000)

	webpage.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "testng工作目录"})
}

func RunTestNG(t provider.T) {
	commandConfig.InitLog()

	testngDir := filepath.Join(constTestHelper.RootPath, "cmd", "test", "demo", "ci_test_testng")
	workspaceName := "testng工作目录"

	commonTestHelper.CloneGit("https://gitee.com/ngtesting/ci_test_testng.git", testngDir)
	shellUtils.ExeShellWithOutputInDir("mvn clean package test", testngDir)

	t.ID("5432")
	commonTestHelper.ReplaceLabel(t, "客户端-执行单元测试")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)
	defer func() {
		ztfTestHelper.DeleteWorkspace(webpage, workspaceName)
	}()

	plwConf.DisableErr()
	webpage.WaitForSelectorTimeout(".tree-node-title:has-text('testng工作目录')", 3000)
	locator := webpage.Locator(".tree-node-title:has-text('testng工作目录')")
	c := locator.Count()
	if c == 0 {
		createUnitWorkspace(t, workspaceName, testngDir, "testng", webpage)
	}
	plwConf.EnableErr()

	locator = webpage.Locator(".tree-node>>.tree-node-title", playwright.PageLocatorOptions{HasText: workspaceName})
	c = locator.Count()
	if c == 0 {
		t.Errorf("Find workspace fail")
		t.FailNow()
	}

	locator.RightClick()
	webpage.Click(".tree-context-menu>>text=执行")
	webpage.WaitForSelectorTimeout("#tabsPane >> text=执行", 3000)

	locator = webpage.Locator("#tabsPane>>.form-item:has-text('测试命令')>>input")
	locator.FillNth(0, "mvn clean package test")
	webpage.Check(`#tabsPane >> input[type="checkbox"]`)
	locator = webpage.Locator("#tabsPane:has-text('禅道测试单标题')>>input")
	locator.FillNth(2, "test unit")
	webpage.Click("#tabsPane >> text=执行")

	plwConf.DisableErr()
	err := webpage.WaitForSelector("#log-list")
	if err != nil {
		webpage.Click("#tabsPane >> text=执行")
	}
	plwConf.EnableErr()

	webpage.WaitForSelectorTimeout("#log-list>>.msg-span>>:has-text('执行3个用例，耗时')", 20000)
	locator = webpage.Locator("#log-list>>code:has-text('执行3个用例，耗时')")
	innerText := locator.InnerText()

	if !strings.Contains(innerText, "通过数：2(66.0%)，失败数：1(33.0%)") {
		t.Errorf("Exec testng fail, result:" + innerText)
		t.FailNow()
	}
	webpage.WaitForSelector("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('testng工作目录(3)')")

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text>>nth=0")
	if logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail")
		t.FailNow()
	}

	isSuccess := zentaoTestHelper.CheckUnitTestResult()
	if !isSuccess {
		webpage.ScreenShot()
		t.Errorf("Exec testng unit fail")
		t.FailNow()
	}
}

func TestUiRunUnit(t *testing.T) {
	runner.Run(t, "客户端-执行TestNG单元测试", RunTestNG)
}
