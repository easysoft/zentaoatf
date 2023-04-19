import http from 'k6/http';
import exec from 'k6/execution';

import { check, sleep, group } from "k6";

export const options = {
    // 并发用户数及其加载方式
    stages: [
        { duration: "2s", target: 3 }, // 1秒内，加载3个虚拟用户
        { duration: "1s", target: 0 }, // 1秒内，销毁所有虚拟用户
    ],

    // 设置性能有关指标的阀值
    thresholds: {
        // '登录请求'分组中所有请求的响应时间，90%小于6000毫秒
        'http_req_duration{group:::登录请求}': ['p(90) < 6000'],

        // '登录请求'分组的整体执行时间，平均值小于6000毫秒
        'group_duration{group:::登录请求}': ['avg < 6000'],
    },
};

export default function () {
    // 执行脚本
    group('登录请求', function () { // 单元测试套件
        let resp = http.get('https://httpbin.org/get?p1=1');

        // 期待响应状态码
        let expectRespStatus = 200

        // 通过设置错误的期待响应状态码，模拟三分之一的迭代失败
        if (+`${exec.vu.idInTest}` % 3 === 0) {
            expectRespStatus = 222
        }
        console.log(`in iteration ${exec.vu.idInTest}, expectRespStatus=${expectRespStatus}`)

        // 验证器
        const validator = (r) => resp.status == expectRespStatus

        // 断言：用例ID, 用例名称，验证点名称，被验证数据，验证器
        assert(0, '微信扫码登录', '验证跳转到个人仪表盘', resp, validator);

        // 提交到结果到禅道里存在的用例上
        assert(1, '邮箱登录', '验证到达首页', resp, (r) => resp.status == 200);

        sleep(1);
    });

    group('用户管理', function () { // 单元测试套件
        const resp = http.post('https://httpbin.org/post', JSON.stringify({
            foo: 'abc',
            bar: '123',
        }), {
            tags: { foo: 'bar' },
            headers: {
                'Content-Type': 'application/json',
            },
        });

        // 验证器
        const validator = (data) => {
            const status = data.status
            const dur = data.timings.duration
            // console.log('===', status, dur)

            const pass = status == 200 && dur < 3000
            return pass
        }

        // 断言：用例ID, 用例名称，验证点名称，被验证数据，验证器
        assert(0, '重置密码', '验证用户收到密码重置右键', resp, validator);

        sleep(1);
    });
}

export function setup() {
    console.log('--- setup')
}

export function teardown(data) {
    console.log('--- teardown')
}

function assert (caseId, caseName, checkpoint, data, validator) {
    const name = `${caseId} - ${caseName}`

    check(data, { [name]: validator }, { id: caseId, name: caseName, checkpoint: checkpoint });
}