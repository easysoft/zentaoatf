import request from '@/utils/request';
import { Config } from './data.d';

const apiPath = 'config';

export async function saveConfig(params: Partial<Config>): Promise<any> {
    return request({
        url: `/${apiPath}/saveConfig`,
        method: 'POST',
        data: params,
    });
}

