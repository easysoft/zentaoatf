import request from '@/utils/request';
import { QueryParams } from '@/types/data.d';

const apiPath = 'sites';

export async function list(params?: QueryParams): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}

export async function get(seq: string): Promise<any> {
    return request({
        url: `/${apiPath}/${seq}`
    });
}

export async function remove(seq: string): Promise<any> {
    return request({
        url: `/${apiPath}/${seq}`,
        method: 'delete',
    });
}
