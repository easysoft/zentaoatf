package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	commonTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/common"
	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	ztfTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/ztf"
	plwConf "github.com/easysoft/zentaoatf/cmd/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/cmd/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func RunScript(t provider.T) {
	t.ID("5743")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	scriptLocator := webpage.Locator(".tree-node-title>>text=1_string_match.php")
	scriptLocator.Click()
	webpage.WaitForTimeout(2000)

	selectLocalProxy(webpage)

	webpage.Click(".tabs-nav-toolbar>>[title=\"执行\"]")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	innerText := webpage.InnerText("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if !strings.Contains(innerText, "失败数：1(100.0%)") {
		t.Errorf("Exec 1_string_match.php fail")
		t.FailNow()
	}
	webpage.WaitForSelector("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('1_string_match.php')")

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result time in rightPane fail")
		t.FailNow()
	}
}

func CheckStep(t provider.T) {
	t.ID("7609")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	scriptLocator := webpage.Locator(".tree-node-title>>text=1_string_match.php")
	scriptLocator.Click()
	webpage.WaitForTimeout(2000)

	selectLocalProxy(webpage)

	webpage.Click(".tabs-nav-toolbar>>[title=\"执行\"]")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	innerText := webpage.InnerText("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if !strings.Contains(innerText, "失败数：1(100.0%)") {
		t.Errorf("Exec 1_string_match.php fail")
		t.FailNow()
	}
	webpage.WaitForSelector("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('1_string_match.php')")

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result time in rightPane fail")
		t.FailNow()
	}

	//check step
	webpage.Click("#show-detail")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('共有3个步骤，0个通过，3个失败。')")
}

func RunScriptByRightClick(t provider.T) {
	t.ID("5479")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	scriptLocator := webpage.Locator(".tree-node-title>>text=1_string_match.php")
	scriptLocator.RightClick()
	webpage.Click(".tree-context-menu>>text=执行")

	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	innerText := webpage.InnerText("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if !strings.Contains(innerText, "失败数：1(100.0%)") {
		t.Errorf("Exec 1_string_match.php fail")
		t.FailNow()
	}

	webpage.WaitForSelector("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('1_string_match.php')")

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:4] != resultTime[:4] {
		t.Errorf("Find result time in rightPane fail")
		t.FailNow()
	}
}

func RunNoInterpreterScript(t provider.T) {
	t.ID("5501")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	scriptLocator := webpage.Locator("text=1_string_match.rb")
	scriptLocator.Click()
	webpage.Click(".tabs-nav-toolbar>>[title=\"执行\"]")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('忽略1个未设置解析器的脚本')")
}

func RunSelectedScripts(t provider.T) {
	t.ID("5481")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	webpage.Click(`[title="批量选择"]`)
	scriptLocator := webpage.Locator(".tree-node-item:has-text('1_string_match.php')")
	scriptLocator = scriptLocator.Locator(".tree-node-check")
	scriptLocator.Click()
	scriptLocator = webpage.Locator(".tree-node-item:has-text('2_webpage_extract.php')")
	scriptLocator = scriptLocator.Locator(".tree-node-check")
	scriptLocator.Click()

	webpage.Click(".run-selected")

	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行2个用例，耗时')")
	innerText := webpage.InnerText("#log-list>>.msg-span>>:has-text('执行2个用例，耗时')")
	if !strings.Contains(innerText, "通过数：0(0.0%)，失败数：2(100.0%)") {
		t.Errorf("Exec 1_string_match.php,2_webpage_extract.php fail")
		t.FailNow()
	}

	webpage.WaitForSelector(fmt.Sprintf("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('%s(2)')", constTestHelper.WorkspaceName))

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	webpage.WaitForTimeout(1000)
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result time in rightPane fail")
		t.FailNow()
	}
}

func RunOpenedAndLast(t provider.T) {
	t.ID("5484")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	webpage.Click(".tree-node-item:has-text('1_string_match.php')")
	webpage.Click(".tree-node-item:has-text('2_webpage_extract.php')")

	webpage.Click("#batchRunMenuToggle")
	webpage.Click(".list-item-content:has-text('执行打开文件')")

	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行2个用例，耗时')")
	locator := webpage.Locator("#log-list>>code:has-text('执行2个用例，耗时')")
	innerText := locator.InnerText()
	if !strings.Contains(innerText, "通过数：0(0.0%)，失败数：2(100.0%)") {
		t.Errorf("Exec opened scripts fail")
		t.FailNow()
	}

	webpage.WaitForSelector(fmt.Sprintf("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('%s(2)')", constTestHelper.WorkspaceName))

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail")
	}

	webpage.Click(`#bottomPane>>[title="清空"]`)

	webpage.Click("#batchRunMenuToggle")
	webpage.Click(".list-item-content:has-text('执行上次')")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行2个用例，耗时')")
	locator = webpage.Locator("#log-list>>code:has-text('执行2个用例，耗时')")
	innerText = locator.InnerText()
	if !strings.Contains(innerText, "通过数：0(0.0%)，失败数：2(100.0%)") {
		t.Errorf("Exec last time fail")
		t.FailNow()
	}
	webpage.WaitForSelector(fmt.Sprintf("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('%s(2)')", constTestHelper.WorkspaceName))

	timeElement = webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime = timeElement.InnerText()
	resultTime = webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail")
		t.FailNow()
	}
}

func RunAll(t provider.T) {
	t.ID("5485")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	os.RemoveAll(commonTestHelper.GetZtfProductPath())
	os.Remove(commonTestHelper.GetPhpWorkspacePath() + "1.php")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	runInfo := "执行5个用例，耗时"
	runRes := "通过数：1(20.0%)，失败数：4(80.0%)"
	resTitle := fmt.Sprintf("%s(5)", constTestHelper.WorkspaceName)
	if runtime.GOOS == "windows" {
		runInfo = "执行4个用例，耗时"
		runRes = "通过数：1(25.0%)，失败数：3(75.0%)"
		resTitle = fmt.Sprintf("%s(4)", constTestHelper.WorkspaceName)
	}

	webpage.Click("#batchRunMenuToggle")
	webpage.Click(".list-item-content:has-text('执行所有文件')")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('" + runInfo + "')")
	locator := webpage.Locator("#log-list>>code:has-text('" + runInfo + "')")
	innerText := locator.InnerText()
	if !strings.Contains(innerText, runRes) {
		t.Errorf("Exec all fail")
		t.FailNow()
	}

	webpage.WaitForSelector("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('" + resTitle + "')")

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail")
	}
}

func RunWorkspace(t provider.T) {
	t.ID("5482")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	webpage.WaitForSelector(".tree-node")
	locator := webpage.Locator(fmt.Sprintf(".tree-node-root>>.tree-node-title>> :scope:has-text('%s')", constTestHelper.WorkspaceName))
	c := locator.Count()
	if c == 0 {
		t.Errorf("Find workspace fail")
		t.FailNow()
	}

	runInfo := "执行4个用例，耗时"
	runRes := "通过数：1(25.0%)，失败数：3(75.0%)"
	resTitle := fmt.Sprintf("%s(4)", constTestHelper.WorkspaceName)
	if runtime.GOOS == "windows" {
		runInfo = "执行3个用例，耗时"
		runRes = "通过数：1(33.0%)，失败数：2(66.0%)"
		resTitle = fmt.Sprintf("%s(3)", constTestHelper.WorkspaceName)
	}

	locator.RightClick()
	webpage.Click(".tree-context-menu>>text=执行")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('" + runInfo + "')")
	locator = webpage.Locator("#log-list>>code:has-text('" + runInfo + "')")
	innerText := locator.InnerText()
	if !strings.Contains(innerText, runRes) {
		t.Errorf("Exec workspace fail")
		t.FailNow()
	}

	webpage.WaitForSelector("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('" + resTitle + "')")

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail")
		t.FailNow()
	}
}

func RunDir(t provider.T) {
	t.ID("5480")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	ztfTestHelper.ExpandWorspace(webpage)

	webpage.RightClick(".tree-node-children>>.tree-node>>:has-text('testdir')")
	webpage.Click(".tree-context-menu>>text=执行")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	innerText := webpage.InnerText("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if !strings.Contains(innerText, "失败数：1(100.0%)") {
		t.Errorf("Exec 1_string_match.php fail")
		t.FailNow()
	}

	webpage.WaitForSelector("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('1_string_match.php')")

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail")
		t.FailNow()
	}
}

func RunUseProxy(t provider.T) {
	t.ID("5746")
	commonTestHelper.ReplaceLabel(t, "客户端-执行脚本")

	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage, "")
	CreateProxyAndInterpreter(webpage, t)
	ztfTestHelper.ExpandWorspace(webpage)

	scriptLocator := webpage.Locator(".tree-node-title>>text=1_string_match.php")
	scriptLocator.Click()

	webpage.Click("#proxyMenuToggle")
	webpage.Click(".list-item-title:has-text('测试执行节点')")
	webpage.Click(".tabs-nav-toolbar>>[title=\"执行\"]")
	webpage.WaitForSelector("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	innerText := webpage.InnerText("#log-list>>.msg-span>>:has-text('执行1个用例，耗时')")
	if !strings.Contains(innerText, "失败数：1(100.0%)") {
		t.Errorf("Exec 1_string_match.php fail")
		t.FailNow()
	}

	webpage.WaitForSelector("#rightPane>>.result-list-item>>nth=0>>.list-item-title:has-text('1_string_match.php')")

	timeElement := webpage.Locator("#log-list>>code:has-text('开始任务')>>.time>>span")
	logTime := timeElement.InnerText()
	resultTime := webpage.InnerText("#rightPane .result-list-item .list-item-trailing-text")
	if logTime[:5] != resultTime {
		t.Errorf("Find result in rightPane fail")
		t.FailNow()
	}
	selectLocalProxy(webpage)
}

func CreateProxyAndInterpreter(webpage plwHelper.Webpage, t provider.T) {
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.WaitForSelector("#proxyTable>>.z-tbody-td>>:has-text('本地节点')", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateAttached})

	plwConf.DisableErr()
	locator := webpage.Locator("#proxyTable>>.z-tbody-td>>:has-text('测试执行节点')")
	c := locator.Count()
	if c > 0 {
		webpage.Click("#settingModal>>.modal-close")
		plwConf.EnableErr()
		return
	}
	plwConf.EnableErr()

	webpage.Click("#serverTable>>button:has-text('新建执行节点')")
	locator = webpage.Locator("#proxyFormModal input")
	locator.FillNth(0, "测试执行节点")
	webpage.WaitForTimeout(200)
	locator.FillNth(1, "http://127.0.0.1:8085")
	webpage.Click("#proxyFormModal>>text=确定")

	webpage.WaitForSelector("#proxyFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	webpage.WaitForTimeout(1000)

	webpage.Click("#settingModal>>.modal-close")
}

func selectLocalProxy(webpage plwHelper.Webpage) {
	proxy := webpage.InnerText(".flex-align-center>>#proxyMenuToggle")
	if strings.Contains(proxy, "本地") {
		return
	}

	webpage.Click("#proxyMenuToggle")
	webpage.Click(".list-item-title:has-text('本地节点')")
}
func TestUiRun(t *testing.T) {
	runner.Run(t, "客户端-脚本执行步骤和通过数日志", CheckStep)
	runner.Run(t, "客户端-执行单个脚本", RunScript)
	runner.Run(t, "客户端-右键执行单个脚本", RunScriptByRightClick)
	if runtime.GOOS == "windows" {
		runner.Run(t, "客户端-忽略执行未设置解析器的脚本", RunNoInterpreterScript)
	}
	runner.Run(t, "客户端-执行选中的脚本文件和文件夹", RunSelectedScripts)
	runner.Run(t, "客户端-执行打开的脚本文件", RunOpenedAndLast)
	runner.Run(t, "客户端-执行所有的脚本文件", RunAll)
	runner.Run(t, "客户端-右键执行工作目录", RunWorkspace)
	runner.Run(t, "客户端-右键执行文件夹", RunDir)
	runner.Run(t, "客户端-使用代理执行单个脚本", RunUseProxy)
}
