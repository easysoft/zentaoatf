import request from '@/utils/request';

const apiPath = 'bugs';

export function submitBugToZentao(data: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'post',
        data: data,
    });
}

export function prepareBugData(data: any): Promise<any> {
    return request({
        url: `/${apiPath}/prepareBugData`,
        method: 'POST',
        data: data,
    });
}

export function loadBugs(): Promise<any> {
    return request({
        url: `${apiPath}`,
        method: 'get',
    });
}