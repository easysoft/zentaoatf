import request from '@/utils/request';
import {Execution, WsMsg} from './data.d';
import { QueryParams } from '@/types/data.d';
import {logLevelMap} from "@/utils/const";

const apiPath = 'exec';

export async function list(params?: QueryParams): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'get',
        params,
    });
}

export async function get(seq: string): Promise<any> {
    return request({
        url: `/${apiPath}/${seq}`
    });
}

export async function remove(seq: string): Promise<any> {
    return request({
        url: `/${apiPath}/${seq}`,
        method: 'delete',
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

export function genExecInfo(jsn: WsMsg, i: number, logLevel: string): string {
    console.log('===', logLevelMap[logLevel].code, logLevelMap[jsn.category].code)
    if (logLevelMap[jsn.category].code < logLevelMap[logLevel].code) {
        return ''
    }

    let msg = jsn.msg.replace(/^"+/,'').replace(/"+$/,'')
    msg = SetWidth(i + '. ', 40) + `<span>${msg}</span>`;

    const sty = 'color: ' + logLevelMap[jsn.category].color + ';'
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
