package main

import (
	"os"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

var pw, err = os.Getwd()
var (
	workspacePath = pw + "\\demo\\php"
)

func CreateWorkspace(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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

	err = page.Click(`[title="新建工作目录"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#workspaceFormModal")
	locator, err = page.Locator("#workspaceFormModal input")
	if err != nil {
		t.Errorf("Find create workspace input fail: %v", err)
		t.FailNow()
	}
	titleInput, err := locator.Nth(0)
	if err != nil {
		t.Errorf("Find title input fail: %v", err)
		t.FailNow()
	}
	err = titleInput.Fill("单元测试工作目录")
	if err != nil {
		t.Errorf("Fil title input fail: %v", err)
		t.FailNow()
	}
	pathInput, err := locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		t.FailNow()
	}
	err = pathInput.Fill(workspacePath)
	if err != nil {
		t.Errorf("Fil address input fail: %v", err)
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
	_, err = typeInput.SelectOption(playwright.SelectOptionValues{Values: &[]string{"ztf"}})
	if err != nil {
		t.Errorf("Fil name input fail: %v", err)
		t.FailNow()
	}
	langInput, err := locator.Nth(1)
	if err != nil {
		t.Errorf("Find lang input fail: %v", err)
		t.FailNow()
	}
	_, err = langInput.SelectOption(playwright.SelectOptionValues{Values: &[]string{"php"}})
	if err != nil {
		t.Errorf("Fil lang input fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#workspaceFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	var waitTimeOut float64 = 5000
	_, err = page.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		t.Errorf("Wait created workspace result fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created workspace fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func SyncFromZentao(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	var waitTimeOut float64 = 5000
	_, err = page.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		CreateWorkspace(t)
		SyncFromZentao(t)
		return
	}
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		CreateWorkspace(t)
		SyncFromZentao(t)
		return
	}
	locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click node fail: %v", err)
		t.FailNow()
	}
	page.Click(".tree-context-menu>>text=从禅道同步")
	if err != nil {
		t.Errorf("Click sync from zentao fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal hide fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步"})
	c, err = locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Sync from zentao fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func SyncTwoCaseFromZentao(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		t.FailNow()
	}
	if _, err = page.Goto("http://127.0.0.1:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		t.FailNow()
	}
	var waitTimeOut float64 = 5000
	_, err = page.WaitForSelector(".tree-node", playwright.PageWaitForSelectorOptions{Timeout: &waitTimeOut})
	if err != nil {
		CreateWorkspace(t)
		SyncFromZentao(t)
		return
	}
	locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click node fail: %v", err)
		t.FailNow()
	}
	page.Click(".tree-context-menu>>text=从禅道同步")
	if err != nil {
		t.Errorf("Click sync from zentao fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal fail: %v", err)
		t.FailNow()
	}
	err = page.Click("text=编号标题类型状态结果 >> input[type=\"checkbox\"]")
	page.Click("text=1check string matches pattern功能测试正常通过 >> input[type=\"checkbox\"]")
	page.Click("text=2extract content from webpage功能测试正常 >> input[type=\"checkbox\"]")
	err = page.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal hide fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步2个用例"})

	c, err = locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Sync from zentao fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func SyncToZentao(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
		t.Errorf("Right click node fail: %v", err)
		t.FailNow()
	}
	page.Click(".tree-context-menu>>text=同步到禅道")
	if err != nil {
		t.Errorf("Click sync to zentao fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector(".toast-notification-close")
	if err != nil {
		t.Errorf("Wait toast-notification-close fail: %v", err)
		t.FailNow()
	}
	locator, err = page.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功同步"})
	c, err = locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Sync to zentao fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func Copy(t provider.T) {
	t.ID("5474")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	productLocator, err := locator.Locator(".tree-node-item:has-text('product1')")
	if err != nil {
		t.Errorf("Find product1 fail: %v", err)
		t.FailNow()
	}
	err = productLocator.Click()
	if err != nil {
		t.Errorf("Click product1 fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator("text=1_string_match.php")
	if err != nil {
		t.Errorf("Find 1_string_match.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-context-menu>>text=复制")
	if err != nil {
		t.Errorf("Click copy fail: %v", err)
		t.FailNow()
	}
	err = productLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click workspace fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-context-menu>>text=粘贴")
	if err != nil {
		t.Errorf("Click parse fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	scriptLocator, err = locator.Locator(".tree-node-item>>div:has-text('1_string_match.php')")
	c, err = scriptLocator.Count()
	if err != nil || c < 2 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func DeleteScript(t provider.T) {
	t.ID("5478")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	scriptLocator, err := locator.Locator("text=1.php")
	if err != nil {
		t.Errorf("Find 1.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-context-menu>>text=删除")
	if err != nil {
		t.Errorf("Click delete fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	scriptLocator, err = locator.Locator(".tree-node-item>>div:has-text('1.php')")
	c, err = scriptLocator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete workspace fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func DeleteDir(t provider.T) {
	t.ID("5477")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	productLocator, err := locator.Locator(".tree-node-item:has-text('product1')")
	if err != nil {
		t.Errorf("Find product1 fail: %v", err)
		t.FailNow()
	}
	err = productLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-context-menu>>text=删除")
	if err != nil {
		t.Errorf("Click delete fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	scriptLocator, err := locator.Locator(".tree-node-item>>div:has-text('product1')")
	c, err = scriptLocator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete workspace fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func DeleteWorkspace(t provider.T) {
	t.ID("5468")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	locator, err := page.Locator(".tree-node-item", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = locator.Hover()
	if err != nil {
		t.Errorf("The hover workspace fail: %v", err)
		t.FailNow()
	}
	err = page.Click(`[title="删除"]`)
	if err != nil {
		t.Errorf("The click delete fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	scriptLocator, err := page.Locator(".tree-node-title:has-text('单元测试工作目录')")
	c, err = scriptLocator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete workspace fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func Clip(t provider.T) {
	t.ID("5476")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	productLocator, err := locator.Locator(".tree-node-item:has-text('product1')")
	if err != nil {
		t.Errorf("Find product1 fail: %v", err)
		t.FailNow()
	}
	err = productLocator.Click()
	if err != nil {
		t.Errorf("Click product1 fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := locator.Locator("text=1.php")
	if err != nil {
		t.Errorf("Find 1.php fail: %v", err)
		t.FailNow()
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-context-menu>>text=剪切")
	if err != nil {
		t.Errorf("Click copy fail: %v", err)
		t.FailNow()
	}
	workspaceLocator, err := page.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	if err != nil {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}
	err = workspaceLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click workspace fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-context-menu>>text=粘贴")
	if err != nil {
		t.Errorf("Click parse fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(1000)
	scriptLocator, err = locator.Locator(".tree-node-item>>div:has-text('1.php')")
	c, err = scriptLocator.Count()
	if err != nil || c < 1 {
		t.Errorf("Find workspace fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}

func FilterDir(t provider.T) {
	t.ID("5494")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	err = page.Click(`[title="筛选"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#filterModal")
	if err != nil {
		t.Errorf("Wait filter modal fail: %v", err)
		t.FailNow()
	}

	err = page.Click("#filterModal>>div:has-text(\"单元测试工作目录\")")
	if err != nil {
		t.Errorf("The Click php filter fail: %v", err)
		t.FailNow()
	}
	eleArr, err := page.QuerySelectorAll("#leftPane .tree .tree-node")
	if len(eleArr) != 1 {
		t.Errorf("Filter valid fail: %v", err)
		t.FailNow()
	}
	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func FilterSuite(t provider.T) {
	t.ID("5495")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	err = page.Click(`[title="筛选"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#filterModal")
	if err != nil {
		t.Errorf("Wait filter modal fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(600)
	err = page.Click("#filterModal>>.tab-nav:has-text(\"按套件\")")
	if err != nil {
		t.Errorf("The Click by suite fail: %v", err)
		t.FailNow()
	}
	page.WaitForSelector("#filterModal>>.list-item-title:has-text(\"test_suite\")")
	err = page.Click("#filterModal>>.list-item-title:has-text(\"test_suite\")")
	if err != nil {
		t.Errorf("The Click test_suite filter fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(200)
	page.WaitForSelector(".toolbar:has-text(\"按套件\")")
	err = page.Click(".tree-node-title:has-text(\"单元测试工作目录\")")
	page.WaitForSelector(".tree")
	page.WaitForTimeout(200)
	scriptLocator, err := page.Locator(".tree>>text=1_string_match.php")
	c, err := scriptLocator.Count()
	if err != nil || c == 0 {
		t.Errorf("Filter suite fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func ByModule(t provider.T) {
	t.ID("5493")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	err = page.Click(".tree-node-title:has-text(\"单元测试工作目录\")")
	if err != nil {
		t.Errorf("The Click 单元测试工作目录 workspace fail: %v", err)
		t.FailNow()
	}
	err = page.Click("#displayByMenuToggle")
	if err != nil {
		t.Errorf("The Click byModule btn fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(100)
	err = page.Click(".dropdown-menu>>.list-item-content:has-text(\"按模块\")")
	if err != nil {
		t.Errorf("The Click by module fail: %v", err)
		t.FailNow()
	}
	err = page.Click(".tree-node-title:has-text(\"module1\")")
	if err != nil {
		t.Errorf("The Click module1 dir fail: %v", err)
		t.FailNow()
	}
	scriptLocator, err := page.Locator(".tree-node>>:has-text(\"check string matches pattern\")")
	c, err := scriptLocator.Count()
	if err != nil || c == 0 {
		t.Errorf("Filter suite fail: %v", err)
		t.FailNow()
	}
	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func FilterTask(t provider.T) {
	t.ID("5496")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	err = page.Click(`[title="筛选"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		t.FailNow()
	}
	_, err = page.WaitForSelector("#filterModal")
	if err != nil {
		t.Errorf("Wait filter modal fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(600)
	err = page.Click("#filterModal>>.tab-nav:has-text(\"按测试单\")")
	if err != nil {
		t.Errorf("The Click by suite fail: %v", err)
		t.FailNow()
	}
	page.WaitForSelector("#filterModal>>.list-item-title:has-text(\"test_task\")")
	err = page.Click("#filterModal>>.list-item-title:has-text(\"test_task\")")
	page.WaitForTimeout(200)
	if err != nil {
		t.Errorf("The Click test_task filter fail: %v", err)
		t.FailNow()
	}
	page.WaitForSelector(".toolbar:has-text(\"按测试单\")")
	err = page.Click(".tree-node-title:has-text(\"单元测试工作目录\")")
	scriptLocator, err := page.Locator(".tree>>text=1_string_match.php")
	c, err := scriptLocator.Count()
	if err != nil || c == 0 {
		t.Errorf("Filter suite fail: %v", err)
		t.FailNow()
	}

	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func Collapse(t provider.T) {
	t.ID("5472")
	t.AddParentSuite("管理禅道站点下工作目录")
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	headless := true
	var slowMo float64 = 100
	workspaceBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web workspaceBrowser: %v", err)
		t.FailNow()
	}
	page, err := workspaceBrowser.NewPage()
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
	err = page.Click(`#leftPane>>.toolbar>>[title="展开"]`)
	if err != nil {
		t.Errorf("Click expand workspace btn fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(100)
	locator, _ = page.Locator("#leftPane>>.tree-node-item>>text=1_string_match.php")
	count, _ := locator.Count()
	if count == 0 {
		t.Error("Expand workspace fail")
		t.FailNow()
	}
	err = page.Click(`#leftPane>>.toolbar>>[title="折叠"]`)
	if err != nil {
		t.Errorf("Click Collapse workspace btn fail: %v", err)
		t.FailNow()
	}
	page.WaitForTimeout(100)
	locator, _ = page.Locator("#leftPane>>.tree-node-item>>text=1_string_match.php")
	count, _ = locator.Count()
	if count > 0 {
		t.Error("Collapse workspace fail")
		t.FailNow()
	}
	if err = workspaceBrowser.Close(); err != nil {
		t.Errorf("The workspaceBrowser cannot be closed: %v", err)
		t.FailNow()
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		t.FailNow()
	}
}
func TestWorkspace(t *testing.T) {
	runner.Run(t, "客户端-从禅道同步部分用例", SyncTwoCaseFromZentao)
	runner.Run(t, "客户端-从禅道同步", SyncFromZentao)
	runner.Run(t, "客户端-同步到禅道", SyncToZentao)
	runner.Run(t, "客户端-复制粘贴树状脚本文件", Copy)
	runner.Run(t, "客户端-剪切粘贴树状脚本文件", Clip)
	runner.Run(t, "客户端-删除树状脚本文件", DeleteScript)
	runner.Run(t, "客户端-删除树状脚本文件夹", DeleteDir)
	runner.Run(t, "客户端-按目录过滤禅道用例脚本", FilterDir)
	runner.Run(t, "客户端-按套件过滤禅道用例脚本", FilterSuite)
	runner.Run(t, "客户端-按测试单过滤禅道用例脚本", FilterTask)
	runner.Run(t, "客户端-显示展开折叠脚本树状结构", Collapse)
	runner.Run(t, "客户端-按模块展示禅道用例脚本", ByModule)
	runner.Run(t, "客户端-删除禅道工作目录", DeleteWorkspace)
	runner.Run(t, "客户端-创建禅道工作目录", CreateWorkspace)
}
