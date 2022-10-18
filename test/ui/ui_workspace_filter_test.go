package main

import (
	"testing"

	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

func FilterDir(t provider.T) {
	t.ID("5494")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.Click(`[title="筛选"]`)
	webpage.WaitForSelector("#filterModal")
	webpage.Click("#filterModal>>.list-item-title:has-text(\"单元测试工作目录\")")
	eleArr := webpage.QuerySelectorAll("#leftPane>>.tree>>.tree-node")
	if len(eleArr.ElementHandles) < 1 {
		t.Errorf("Filter valid fail")
		t.FailNow()
		return
	}
}
func FilterSuite(t provider.T) {
	t.ID("5495")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.Click(`[title="筛选"]`)
	webpage.WaitForSelector("#filterModal")
	webpage.WaitForTimeout(1000)
	webpage.Click("#filterModal>>.tab-nav:has-text(\"按套件\")")
	webpage.WaitForSelector("#filterModal>>.list-item-title:has-text(\"test_suite\")")
	webpage.Click("#filterModal>>.list-item-title:has-text(\"test_suite\")")
	webpage.WaitForTimeout(200)
	webpage.WaitForSelector(".toolbar:has-text(\"按套件\")")
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.WaitForTimeout(200)
	webpage.Locator(".tree-node>>text=1_string_match.php")
}
func ByModule(t provider.T) {
	t.ID("5493")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.Click("#displayByMenuToggle")
	webpage.WaitForTimeout(1000)
	webpage.Click(".dropdown-menu>>.list-item-content:has-text(\"按模块\")")

	webpage.Click(".tree-node-title:has-text(\"module1\")")

	// scriptLocator := webpage.Locator(".tree-node>>:has-text(\"check string matches pattern\")")
	// c := scriptLocator.Count()
	// if err != nil || c == 0 {
	// 	t.Errorf("Filter suite fail: %v")
	// 	t.FailNow()
	// }
}
func FilterTask(t provider.T) {
	t.ID("5496")
	t.AddParentSuite("管理禅道站点下工作目录")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.Click(`[title="筛选"]`)
	webpage.WaitForSelector("#filterModal")
	webpage.WaitForTimeout(1000)
	webpage.Click("#filterModal>>.tab-nav:has-text(\"按测试单\")")
	webpage.WaitForSelector("#filterModal>>.list-item-title:has-text(\"企业网站第一期测试任务\")")
	webpage.Click("#filterModal>>.list-item-title:has-text(\"企业网站第一期测试任务\")")
	webpage.WaitForTimeout(200)
	webpage.WaitForSelector(".toolbar:has-text(\"按测试单\")")
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.Locator(".tree-node>>text=1_string_match.php")
}

func TestUiWorkspaceFilter(t *testing.T) {
	runner.Run(t, "客户端-按目录过滤禅道用例脚本", FilterDir)
	runner.Run(t, "客户端-按套件过滤禅道用例脚本", FilterSuite)
	runner.Run(t, "客户端-按测试单过滤禅道用例脚本", FilterTask)
	runner.Run(t, "客户端-按模块展示禅道用例脚本", ByModule)
}
