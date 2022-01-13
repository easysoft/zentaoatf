import request from '@/utils/request';
import {Config} from "@/views/config/data";

const apiPath = 'zentao';

export async function queryLang(): Promise<any> {
    return request({
        url: `/${apiPath}/listLang`,
        method: 'GET',
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
