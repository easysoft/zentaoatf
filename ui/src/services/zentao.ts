import request from '@/utils/request';

const apiPath = 'zentao';
const apiPathBug = 'bug';

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

export function getBugSteps(data: any): Promise<any> {
    return request({
        url: `/${apiPathBug}/getBugSteps`,
        method: 'POST',
        data: data,
    });
}

export async function getDataForBugSubmition(data: any): Promise<any> {
    return request({
        url: `/${apiPath}/getDataForBugSubmition`,
        method: 'POST',
        data: data,
    });
}