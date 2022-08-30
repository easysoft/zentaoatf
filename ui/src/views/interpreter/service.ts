import {QueryParams} from "@/types/data";
import request from "@/utils/request";

const apiPath = 'interpreters';

export async function getLangSettings(proxyPath = 'local') {
    const json: any = await request({
        url: '/interpreters/getLangSettings',
        method: 'GET',
        params: {proxyPath}
    })

    if (json.code === 0) {
        return json.data
    }

    return {languages: [], languageMap: {}}
}

export async function getLangInterpreter(lang, proxyPath = 'local') {
    const params = {language: lang, proxyPath}
    const json: any = await request({
        url: '/interpreters/getLangInterpreter',
        method: 'GET',
        params,
    })

    if (json.code === 0) {
        return json.data
    }

    return {path: '', info: ''}
}

export async function listInterpreter(proxyPath: string): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params: {proxyPath},
    });
}

export async function getInterpreter(seq: number, proxyPath = 'local'): Promise<any> {
    return request({
        url: `/${apiPath}/${seq}`,
        params:{proxyPath}
    });
}

export async function saveInterpreter(params: any, proxyPath = 'local'): Promise<any> {
    return request({
        url: `/${apiPath}` + (params.id ? `/${params.id}` : ''),
        method: params.id? 'PUT': 'POST',
        data: params,
        params: {proxyPath},
    });
}

export async function removeInterpreter(id: number, proxyPath = 'local'): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
        data:{proxyPath}
    });
}