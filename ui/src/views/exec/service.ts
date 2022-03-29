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

export function genExecInfo(jsn: WsMsg) : WsMsg {
    jsn.msg = jsn.msg.replace(/^"+/,'').replace(/"+$/,'')
    return jsn
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
