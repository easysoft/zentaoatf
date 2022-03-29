import request from '@/utils/request';
import {WsMsg} from "@/views/result/data";

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

export async function remove(params: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'delete',
        params
    });
}

export function getCaseIdsFromReport(report: any, scope: string): string[] {
    const ret = new Array<string>()

    report.funcResult.forEach(item => {
        const path = item.path
        const status = item.status
        const selected = scope === 'all' || scope === status
        if (path && selected) ret.push(path)
    })

    return ret
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

export function submitResultToZentao(data: any): Promise<any> {
    return request({
        url: `/result`,
        method: 'post',
        data: data,
    });
}
export function submitBugToZentao(data: any): Promise<any> {
    return request({
        url: `/bug`,
        method: 'post',
        data: data,
    });
}
