import request from "@/utils/request";

const apiPath = 'proxies';

export async function listProxy(params?: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}

export async function getProxy(seq: number): Promise<any> {
    return request({
        url: `/${apiPath}/${seq}`
    });
}

export async function saveProxy(params: any): Promise<any> {
    return request({
        url: `/${apiPath}` + (params.id ? `/${params.id}` : ''),
        method: params.id? 'PUT': 'POST',
        data: params,
        params:{proxyPath: params.proxyPath}
    });
}

export async function removeProxy(id: number, proxyPath: string): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
        params: {proxyPath: proxyPath}
    });
}

export async function checkProxy(seq: number): Promise<any> {
    return request({
        url: `/${apiPath}/${seq}/check`
    });
}