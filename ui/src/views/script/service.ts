import request from '@/utils/request';
import { Script } from './data.d';

const apiPath = 'scripts';
const apiPathFilters = 'filters';

export async function listFilterItems(filerType: string): Promise<any> {
    const params = {filerType: filerType}
    return request({
        url: `/${apiPathFilters}/listItems`,
        params
    });
}

export async function get(path: string): Promise<any> {
    const params = {path: path}

    return request({
        url: `/${apiPath}/get`,
        params
    });
}

export async function extract(path: string): Promise<any> {
    const params = {path: path}

    return request({
        url: `/${apiPath}/extract`,
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
