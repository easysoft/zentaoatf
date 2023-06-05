import http from 'k6/http';
import exec from 'k6/execution';

import { check, sleep, group } from "k6";

export const options = {
    // 设置并发用户数及其加载方式
    stages: [
        { duration: "2s", target: 3 }, // 2秒内，加载3个虚拟用户
        { duration: "1s", target: 0 }, // 1秒内，销毁所有虚拟用户
    ],

    // 设置性能有关指标的阀值,形如{id:1}标签的指标会影响指定用例的成败。
    thresholds: {
        // 验证编号为1的用例的响应时间，平均值小于1000毫秒
        'http_req_duration{id:1}': ['avg < 6000'],

        // '用户登录'分组中所有请求的响应时间，90%小于6000毫秒
        'http_req_duration{group:::登录请求}': ['p(90) < 6000'],

        // '用户登录'分组的整体执行时间，平均值小于6000毫秒
        'group_duration{group:::登录请求}': ['avg < 6000'],
    },
};

export default function () {
    // 单元测试套件
    group('用户登录', function () {
        let resp = http.get('https://httpbin.org/get?p1=1', {
            tags: { id: '1' }, // 标记用例编号为1，用于上述阀值统计
        });

        // 期待响应状态码
        let expectRespStatus = 200

        // 通过设置错误的期待结果，模拟三分之一的迭代失败
        // if (+`${exec.vu.idInTest}` % 3 === 0) {
        //     expectRespStatus = 222
        // }
        // console.log(`in iteration ${exec.vu.idInTest}, expectRespStatus=${expectRespStatus}`)

        // 验证器方法，可以验证响应的状态码、耗时、内容等
        const validator = (r) => resp.status == expectRespStatus

        // 注意：此处的检查点和前面定义的阀值'http_req_duration{id:1}'均会影响用例的成败
        // 断言方法（用例ID, 用例名称，验证点名称，被验证数据，验证器）
        assert(1, '微信扫码登录', '验证跳转到个人仪表盘', resp, validator);

        // 子套件'用户登录-登录次数限制'及其下用例代码
        group('失败次数限制', function () {
            assert(0, '登录失败连续3次，账号锁定', '显示账号已被锁定', resp, (r) => true);
            assert(0, '登录非连续性失败累计达到3次，账号不锁定', '可成功登录', resp, (r) => true);
            assert(0, '锁定账号15分钟后自动解禁', '解禁后可成功登录', resp, (r) => true);
        });

        // 当前套件'用户登录'下的另一个用例
        group('CASE', function () {
            assert(0, '邮箱登录', '验证到达首页', resp, (r) => resp.status == 200);
        });

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

            // 验证状态码和响应时间
            const pass = status == 200 && dur < 6000

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
// 配置ZTF执行时请保留该函数，否则thresholds阀值结果不会影响用例结果
// export function handleSummary(data) {
//     return {
//         'results/summary.json': JSON.stringify(data), //the default data object
//     };
// }

function assert (caseId, caseName, checkpoint, data, validator) {
    const name = `${caseId} - ${caseName}`
    const tags = { id: caseId, name: caseName, checkpoint: checkpoint}

    check(data, { [name]: validator }, tags);
}