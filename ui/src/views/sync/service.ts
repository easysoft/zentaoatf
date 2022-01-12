import request from '@/utils/request';

const apiPath = 'zentao';

export async function queryProduct(): Promise<any> {
    const params = {}

    return request({
        url: `/${apiPath}/listProduct`,
        method: 'GET',
        params,
    });
}

export async function queryModule(productId): Promise<any> {
    const params = {productId: productId}

    return request({
        url: `/${apiPath}/listModule`,
        method: 'GET',
        params,
    });
}

export async function querySuite(productId): Promise<any> {
    const params = {productId: productId}

    return request({
        url: `/${apiPath}/listSuite`,
        method: 'GET',
        params,
    });
}

export async function queryTask(productId): Promise<any> {
    const params = {productId: productId}

    return request({
        url: `/${apiPath}/listTask`,
        method: 'GET',
        params,
    });
}
