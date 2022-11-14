# coding=utf-8

import ctypes
import time
import pytest
import allure


@allure.feature('测试API')
class TestCApi(object):
    # b155282bc7840b1afa4df2576096f649
    # b155282bc7840b1afa4df2576096f649
    @allure.issue('444444')
    @allure.link('sssss')
    @allure.story('测试初始化')
    @allure.testcase('5555555')
    @allure.id('11')
    @pytest.mark.parametrize('para', ['asd', 'rtrtrt'], ids=['1', '2'])
    def test_c_api_init(self, para):
        # allure.getLifecycle()
        # lifecycle.updateTestCase(testResult -> testResult.setTestCaseId("1"));
        a = 1 + 1
        print(f"日志个数为{a}")
