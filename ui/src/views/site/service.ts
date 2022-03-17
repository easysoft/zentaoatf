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

export async function get(seq: number): Promise<any> {
    return request({
        url: `/${apiPath}/${seq}`
    });
}

export async function save(params: any): Promise<any> {
    return request({
        url: `/${apiPath}` + (params.id ? `/${params.id}` : ''),
        method: params.id? 'PUT': 'POST',
        data: params,
    });
}

export async function remove(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
    });
}
