import request from '@/utils/request';
import { QueryParams } from '@/types/data.d';

const apiPath = 'workspaces';

export async function query(params?: QueryParams): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}

export async function listByProduct(productId: number): Promise<any> {
    const params = {productId: productId}

    return request({
        url: `/${apiPath}/listByProduct`,
        method: 'get',
        params,
    });
}

export async function get(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`
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
