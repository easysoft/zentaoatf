import request from '@/utils/request';

const apiPath = 'bugs';

export function submitBugToZentao(data: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'post',
        data: data,
    });
}
