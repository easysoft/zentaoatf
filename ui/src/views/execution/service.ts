import request from '@/utils/request';
import { Execution, QueryParams } from './data.d';

const apiPath = 'exec';

export async function query(params?: QueryParams): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}

export async function create(params: Omit<Execution, 'id'>): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'POST',
        data: params,
    });
}

export async function update(id: number, params: Omit<Execution, 'id'>): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'PUT',
        data: params,
    });
}

export async function remove(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
    });
}

export async function detail(id: number): Promise<any> {
    return request({url: `/executions/${id}`});
}
