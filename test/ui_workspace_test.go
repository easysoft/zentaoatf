package main

import (
	"fmt"
	"testing"

	playwright "github.com/playwright-community/playwright-go"
)

var (
	workspacePath = "/root/ztf/lang/javascript"
)

func CreateWorkspace(t *testing.T) {
	fmt.Println(workspacePath)
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}

	Locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for workspace list nav fail: %v", err)
		return
	}
	err = page.Click(".list-item-title>>text=testSite")
	if err != nil {
		t.Errorf("The Click workspace nav fail: %v", err)
		return
	}

	err = page.Click(`[title="新建工作目录"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#workspaceFormModal")
	Locator, err = page.Locator("#workspaceFormModal input")
	if err != nil {
		t.Errorf("Find create workspace input fail: %v", err)
		return
	}
	titleInput, err := Locator.Nth(0)
	if err != nil {
		t.Errorf("Find title input fail: %v", err)
		return
	}
	err = titleInput.Fill("单元测试工作目录")
	if err != nil {
		t.Errorf("Fil title input fail: %v", err)
		return
	}
	pathInput, err := Locator.Nth(1)
	if err != nil {
		t.Errorf("Find address input fail: %v", err)
		return
	}
	err = pathInput.Fill(workspacePath)
	if err != nil {
		t.Errorf("Fil address input fail: %v", err)
		return
	}
	Locator, err = page.Locator("#workspaceFormModal select")
	if err != nil {
		t.Errorf("Find create workspace select fail: %v", err)
		return
	}
	typeInput, err := Locator.Nth(0)
	if err != nil {
		t.Errorf("Find name input fail: %v", err)
		return
	}
	_, err = typeInput.SelectOption(playwright.SelectOptionValues{Values: &[]string{"ztf"}})
	if err != nil {
		t.Errorf("Fil name input fail: %v", err)
		return
	}
	langInput, err := Locator.Nth(1)
	if err != nil {
		t.Errorf("Find lang input fail: %v", err)
		return
	}
	_, err = langInput.SelectOption(playwright.SelectOptionValues{Values: &[]string{"javascript"}})
	if err != nil {
		t.Errorf("Fil lang input fail: %v", err)
		return
	}
	err = page.Click("#workspaceFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		return
	}
	Locator, err = page.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find created workspace fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_create.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}

func SyncFromZentao(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait treeNode fail: %v", err)
		return
	}
	Locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	Locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click node fail: %v", err)
		return
	}
	page.Click(".tree-context-menu>>text=从禅道同步")
	if err != nil {
		t.Errorf("Click sync from zentao fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#syncFromZentaoFormModal .z-tbody-checkbox")
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal fail: %v", err)
		return
	}
	err = page.Click("#syncFromZentaoFormModal>>.modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#syncFromZentaoFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateHidden})
	if err != nil {
		t.Errorf("Wait syncFromZentaoFormModal hide fail: %v", err)
		return
	}
	Locator, err = page.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功从禅道同步"})
	c, err = Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Sync from zentao fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_create.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func SyncToZentao(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		return
	}
	Locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	err = Locator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click node fail: %v", err)
		return
	}
	page.Click(".tree-context-menu>>text=同步到禅道")
	if err != nil {
		t.Errorf("Click sync to zentao fail: %v", err)
		return
	}
	_, err = page.WaitForSelector(".toast-notification-close")
	if err != nil {
		t.Errorf("Wait toast-notification-close fail: %v", err)
		return
	}
	Locator, err = page.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "成功同步"})
	c, err = Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Sync to zentao fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_create.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func Copy(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		return
	}
	Locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		return
	}
	productLocator, err := Locator.Locator(".tree-node-item:has-text('product1')")
	if err != nil {
		t.Errorf("Find product1 fail: %v", err)
		return
	}
	err = productLocator.Click()
	if err != nil {
		t.Errorf("Click product1 fail: %v", err)
		return
	}
	scriptLocator, err := Locator.Locator("text=1_string_match.js")
	if err != nil {
		t.Errorf("Find 1_string_match.js fail: %v", err)
		return
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		return
	}
	err = page.Click(".tree-context-menu>>text=复制")
	if err != nil {
		t.Errorf("Click copy fail: %v", err)
		return
	}
	// workspaceLocator, err := page.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	// if err != nil {
	// 	t.Errorf("Find workspace fail: %v", err)
	// 	return
	// }
	err = productLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click workspace fail: %v", err)
		return
	}
	err = page.Click(".tree-context-menu>>text=粘贴")
	if err != nil {
		t.Errorf("Click parse fail: %v", err)
		return
	}
	page.WaitForTimeout(1000)
	scriptLocator, err = Locator.Locator(".tree-node-item>>div:has-text('1_string_match.js')")
	c, err = scriptLocator.Count()
	if err != nil || c < 2 {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_copy.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func DeleteScript(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		return
	}
	Locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		return
	}
	scriptLocator, err := Locator.Locator("text=1.js")
	if err != nil {
		t.Errorf("Find 1.js fail: %v", err)
		return
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		return
	}
	err = page.Click(".tree-context-menu>>text=删除")
	if err != nil {
		t.Errorf("Click delete fail: %v", err)
		return
	}
	err = page.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		return
	}
	page.WaitForTimeout(1000)
	scriptLocator, err = Locator.Locator(".tree-node-item>>div:has-text('1.js')")
	c, err = scriptLocator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete workspace fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_delete.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func DeleteDir(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		return
	}
	Locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		return
	}
	productLocator, err := Locator.Locator(".tree-node-item:has-text('product1')")
	if err != nil {
		t.Errorf("Find product1 fail: %v", err)
		return
	}
	err = productLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		return
	}
	err = page.Click(".tree-context-menu>>text=删除")
	if err != nil {
		t.Errorf("Click delete fail: %v", err)
		return
	}
	err = page.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		return
	}
	page.WaitForTimeout(1000)
	scriptLocator, err := Locator.Locator(".tree-node-item>>div:has-text('product1')")
	c, err = scriptLocator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete workspace fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_delete_dir.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}

func DeleteWorkspace(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		return
	}
	Locator, err := page.Locator(".tree-node-item", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	err = Locator.Hover()
	if err != nil {
		t.Errorf("The hover workspace fail: %v", err)
		return
	}
	err = page.Click(`[title="删除"]`)
	if err != nil {
		t.Errorf("The click delete fail: %v", err)
		return
	}
	err = page.Click(".modal-action>>span:has-text(\"确定\")")
	if err != nil {
		t.Errorf("The Click submit form fail: %v", err)
		return
	}
	page.WaitForTimeout(1000)
	scriptLocator, err := page.Locator(".tree-node-title:has-text('单元测试工作目录')")
	c, err = scriptLocator.Count()
	if err != nil || c > 0 {
		t.Errorf("Delete workspace fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_delete_dir.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func Clip(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := true
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}
	_, err = page.WaitForSelector(".tree-node")
	if err != nil {
		t.Errorf("Wait tree-node fail: %v", err)
		return
	}
	Locator, err := page.Locator(".tree-node", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	c, err := Locator.Count()
	if err != nil || c == 0 {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("Click node fail: %v", err)
		return
	}
	productLocator, err := Locator.Locator(".tree-node-item:has-text('product1')")
	if err != nil {
		t.Errorf("Find product1 fail: %v", err)
		return
	}
	err = productLocator.Click()
	if err != nil {
		t.Errorf("Click product1 fail: %v", err)
		return
	}
	scriptLocator, err := Locator.Locator("text=1.js")
	if err != nil {
		t.Errorf("Find 1.js fail: %v", err)
		return
	}
	err = scriptLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click script fail: %v", err)
		return
	}
	err = page.Click(".tree-context-menu>>text=剪切")
	if err != nil {
		t.Errorf("Click copy fail: %v", err)
		return
	}
	workspaceLocator, err := page.Locator(".tree-node-title", playwright.PageLocatorOptions{HasText: "单元测试工作目录"})
	if err != nil {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	err = workspaceLocator.Click(playwright.PageClickOptions{Button: playwright.MouseButtonRight})
	if err != nil {
		t.Errorf("Right click workspace fail: %v", err)
		return
	}
	err = page.Click(".tree-context-menu>>text=粘贴")
	if err != nil {
		t.Errorf("Click parse fail: %v", err)
		return
	}
	page.WaitForTimeout(1000)
	scriptLocator, err = Locator.Locator(".tree-node-item>>div:has-text('1.js')")
	c, err = scriptLocator.Count()
	if err != nil || c < 1 {
		t.Errorf("Find workspace fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_copy.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}

func FilterDir(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := false
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}

	Locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for workspace list nav fail: %v", err)
		return
	}
	err = page.Click(".list-item-title>>text=testSite")
	if err != nil {
		t.Errorf("The Click workspace nav fail: %v", err)
		return
	}
	err = page.Click(`[title="筛选"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#filterModal")
	if err != nil {
		t.Errorf("Wait filter modal fail: %v", err)
		return
	}

	err = page.Click("#filterModal>>div:has-text(\"php_test\")")
	if err != nil {
		t.Errorf("The Click php filter fail: %v", err)
		return
	}
	eleArr, err := page.QuerySelectorAll("#leftPane .tree .tree-node")
	if len(eleArr) != 1 {
		t.Errorf("Filter valid fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_create.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func FilterSuite(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := false
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}

	Locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for workspace list nav fail: %v", err)
		return
	}
	err = page.Click(".list-item-title>>text=testSite")
	if err != nil {
		t.Errorf("The Click workspace nav fail: %v", err)
		return
	}
	err = page.Click(`[title="筛选"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#filterModal")
	if err != nil {
		t.Errorf("Wait filter modal fail: %v", err)
		return
	}
	page.WaitForTimeout(600)
	err = page.Click("#filterModal>>.tab-nav:has-text(\"按套件\")")
	if err != nil {
		t.Errorf("The Click by suite fail: %v", err)
		return
	}
	page.WaitForSelector("#filterModal>>.list-item-title:has-text(\"test_suite\")")
	err = page.Click("#filterModal>>.list-item-title:has-text(\"test_suite\")")
	if err != nil {
		t.Errorf("The Click test_suite filter fail: %v", err)
		return
	}
	page.WaitForSelector(".toolbar:has-text(\"按套件\")")
	err = page.Click(".tree-node-title:has-text(\"php_test\")")
	scriptLocator, err := page.Locator(".tree>>text=1_string_match.php")
	c, err := scriptLocator.Count()
	if err != nil || c == 0 {
		t.Errorf("Filter suite fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_create.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func FilterTask(t *testing.T) {
	pw, err := playwright.Run()
	if err != nil {
		t.Error(err)
		return
	}
	headless := false
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		t.Errorf("Fail to launch the web browser: %v", err)
		return
	}
	page, err := browser.NewPage()
	if err != nil {
		t.Errorf("Create the new page fail: %v", err)
		return
	}
	if _, err = page.Goto("http://app.me:8000/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	}); err != nil {
		t.Errorf("The specific URL is missing: %v", err)
		return
	}

	Locator, err := page.Locator("#siteMenuToggle")
	if err != nil {
		t.Errorf("The siteMenuToggle is missing: %v", err)
		return
	}
	err = Locator.Click()
	if err != nil {
		t.Errorf("The Click is fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#navbar .list-item")
	if err != nil {
		t.Errorf("Wait for workspace list nav fail: %v", err)
		return
	}
	err = page.Click(".list-item-title>>text=testSite")
	if err != nil {
		t.Errorf("The Click workspace nav fail: %v", err)
		return
	}
	err = page.Click(`[title="筛选"]`)
	if err != nil {
		t.Errorf("The Click create workspace fail: %v", err)
		return
	}
	_, err = page.WaitForSelector("#filterModal")
	if err != nil {
		t.Errorf("Wait filter modal fail: %v", err)
		return
	}
	page.WaitForTimeout(600)
	err = page.Click("#filterModal>>.tab-nav:has-text(\"按任务\")")
	if err != nil {
		t.Errorf("The Click by suite fail: %v", err)
		return
	}
	page.WaitForSelector("#filterModal>>.list-item-title:has-text(\"test_task\")")
	err = page.Click("#filterModal>>.list-item-title:has-text(\"test_task\")")
	if err != nil {
		t.Errorf("The Click test_suite filter fail: %v", err)
		return
	}
	page.WaitForSelector(".toolbar:has-text(\"按任务\")")
	err = page.Click(".tree-node-title:has-text(\"php_test\")")
	scriptLocator, err := page.Locator(".tree>>text=1_string_match.php")
	c, err := scriptLocator.Count()
	if err != nil || c == 0 {
		t.Errorf("Filter suite fail: %v", err)
		return
	}
	if _, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String("workspace_create.png"),
	}); err != nil {
		t.Errorf("screenshot cannot be created: %v", err)
		return
	}
	if err = browser.Close(); err != nil {
		t.Errorf("The browser cannot be closed: %v", err)
		return
	}
	if err = pw.Stop(); err != nil {
		t.Errorf("The playwright cannot be stopped: %v", err)
		return
	}
}
func TestWorkspace(t *testing.T) {
	t.Run("CreateWorkspace", CreateWorkspace)
	t.Run("SyncFromZentao", SyncFromZentao)
	t.Run("SyncToZentao", SyncToZentao)
	t.Run("Copy", Copy)
	t.Run("Clip", Clip)
	t.Run("DeleteScript", DeleteScript)
	t.Run("DeleteDir", DeleteDir)
	t.Run("DeleteWorkspace", DeleteWorkspace)
	t.Run("FilterDir", FilterDir)
	t.Run("FilterSuite", FilterSuite)
	t.Run("FilterTask", FilterTask)
}
