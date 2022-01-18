import request from '@/utils/request';
import { Script } from './data.d';

const apiPath = 'scripts';

export async function get(path: string): Promise<any> {
    const params = {path: path}

    return request({
        url: `/scripts/get`,
        params
    });
}

export async function create(params: Partial<Script>): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'POST',
        data: params,
    });
}

export async function update(id: number, params: Partial<Script>): Promise<any> {
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
