import {QueryParams} from "@/types/data";
import request from "@/utils/request";

const apiPath = 'interpreters';

export async function getLangSettings() {
    const json: any = await request({
        url: '/interpreters/getLangSettings',
        method: 'GET',
    })

    if (json.code === 0) {
        const data = json.data
        return data
    }

    return {languages: [], languageMap: {}}
}

export async function getLangInterpreter(lang) {
    if (lang === '') return {path: '', info: ''}

    const params = {language: lang}
    const json: any = await request({
        url: '/interpreters/getLangInterpreter',
        method: 'GET',
        params,
    })

    if (json.code === 0) {
        const data = json.data
        return data
    }

    return {path: '', info: ''}
}

export async function listInterpreter(params?: QueryParams): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}

export async function getInterpreter(seq: number): Promise<any> {
    return request({
        url: `/${apiPath}/${seq}`
    });
}

export async function saveInterpreter(params: any): Promise<any> {
    return request({
        url: `/${apiPath}` + (params.id ? `/${params.id}` : ''),
        method: params.id? 'PUT': 'POST',
        data: params,
    });
}

export async function removeInterpreter(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
    });
}