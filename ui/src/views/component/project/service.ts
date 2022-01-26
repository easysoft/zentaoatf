import request from '@/utils/request';

const apiPath = 'file';

export async function listUserHome(parentDir): Promise<any> {
    const params = {parentDir: parentDir}

    return request({
        url: `/${apiPath}/listDir`,
        method: 'get',
        params
    });
}
