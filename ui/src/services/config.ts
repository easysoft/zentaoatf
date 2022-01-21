import request from '@/utils/request';

const apiPath = 'config';

export async function saveConfig(params: any): Promise<any> {
    return request({
        url: `/${apiPath}/saveConfig`,
        method: 'POST',
        data: params,
    });
}

