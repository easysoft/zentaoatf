import request from '@/utils/request';
import {WsMsg} from "@/types/data";

const apiPath = 'results';

export async function list(params: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}

export async function get(params: any): Promise<any> {
    return request({
        url: `/${apiPath}/${params.workspaceId}/${params.seq}`,

    });
}

export async function getLastest(params: any): Promise<any> {
    return request({
        url: `/${apiPath}/latest`,
        method: 'get',
        params,
    });
}

export async function remove(params: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'delete',
        params
    });
}

export function submitResultToZentao(data: any): Promise<any> {
    return request({
        url: `${apiPath}`,
        method: 'post',
        data: data,
    });
}

export function genExecInfo(jsn: WsMsg, i: number): string {
    let msg = jsn.msg.replace(/^"+/,'').replace(/"+$/,'')
    msg = SetWidth(i + '. ', 40) + `<span>${msg}</span>`;

    let sty = ''
    if (jsn.category === 'exec') {
        sty = 'color: #68BB8D;'
    } else if (jsn.category === 'error') {
        sty = 'color: #FC2C25;'
    }

    msg = `<div style="${sty}"> ${msg} </div>`

    return msg
}

export function SetWidth(content: string, width: number): string{
    return `<span style="display: inline-block; width: ${width}px; text-align: right; padding-right: 6px;">${content}</span>`;
}

export function mvLog(data: any): Promise<any> {
    return request({
        url: `${apiPath}/mvLog`,
        method: 'post',
        data: data,
    });
}