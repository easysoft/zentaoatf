import request from '@/utils/request';

const apiPath = 'zentao';

export async function queryLang(): Promise<any> {
    return request({
        url: `/${apiPath}/listLang`,
        method: 'GET',
    });
}


export async function getProfile(): Promise<any> {
    return request({
        url: `/${apiPath}/getProfile`,
        method: 'GET',
    });
}

export async function querySiteAndProduct(params): Promise<any> {
    return request({
        url: `/${apiPath}/listSiteAndProduct`,
        method: 'get',
        params,
    });
}

export async function queryProduct(): Promise<any> {
    return request({
        url: `/${apiPath}/listProduct`,
        method: 'GET',
    });
}

export async function queryModule(productId: string): Promise<any> {
    const params = {productId: productId}

    return request({
        url: `/${apiPath}/listModule`,
        method: 'GET',
        params,
    });
}

export async function querySuite(productId: string): Promise<any> {
    const params = {productId: productId}

    return request({
        url: `/${apiPath}/listSuite`,
        method: 'GET',
        params,
    });
}

export async function queryTask(productId: string): Promise<any> {
    const params = {productId: productId}

    return request({
        url: `/${apiPath}/listTask`,
        method: 'GET',
        params,
    });
}

export async function queryBugFields(): Promise<any> {
    return request({
        url: `/${apiPath}/listBugFields`,
        method: 'GET',
    });
}