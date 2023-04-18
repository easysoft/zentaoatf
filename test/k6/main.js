import http from 'k6/http';
import exec from 'k6/execution';

import { check, sleep, group } from "k6";

export const options = {
    // 并发用户数及其加载方式
    stages: [
        { duration: "1s", target: 3 },
        { duration: "1s", target: 0 },
    ],

    // 设置性能有关指标的阀值
    thresholds: {
        // '登录请求'分组中所有请求的响应时间，90%小于1000毫秒
        'http_req_duration{group:::登录请求}': ['p(90) < 1000'],

        // '登录请求'分组的整体执行时间，平均值小于1000毫秒
        'group_duration{group:::登录请求}': ['avg < 1000'],
    },
};

export default function () {
    // 执行脚本
    group('登录请求', function () { // 模块
        let resp = http.get('https://httpbin.org/get?p1=1');

        // 验证点
        const validator = (r) => resp.status == 2001

        // 断言
        assert(1, '微信扫码登录', '验证跳转到个人仪表盘', resp, validator);
        assert(2, '邮箱登录', '验证跳转到个人仪表盘', resp, validator); // 测试结果分析用

        sleep(1);
    });

    group('用户管理', function () {
        const resp = http.post('https://httpbin.org/post', JSON.stringify({
            foo: 'abc',
            bar: '123',
        }), {
            tags: { foo: 'bar' },
            headers: {
                'Content-Type': 'application/json',
            },
        });

        // 验证点
        const validator = (data) => {
            const status = data.status
            const dur = data.timings.duration
            // console.log('===', status, dur)

            const pass = status == 200 && dur < 3000
            return pass
        }

        // 断言
        assert(2, '重置密码', '验证用户收到密码重置右键', resp, validator);

        sleep(1);
    });
}

function assert (caseId, caseName, checkpoint, data, validator) {
    const name = `${caseId} - ${caseName}`

    check(data, { [name]: validator }, { id: caseId, name: caseName, checkpoint: checkpoint });
}