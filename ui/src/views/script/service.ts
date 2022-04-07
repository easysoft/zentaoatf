import request from '@/utils/request';
import {WsMsg} from './data.d';
import {removeEmptyField} from "@/utils/object";

const apiPath = 'scripts';
const apiPathFilters = 'filters';

export async function listFilterItems(filerType: string): Promise<any> {
    const params = {filerType: filerType}
    return request({
        url: `/${apiPathFilters}/listItems`,
        params
    });
}

export async function list(params): Promise<any> {
    return request({
        url: `/${apiPath}/list`,
        params
    });
}

export async function get(path: string, workspaceId: number): Promise<any> {
    const params = {path: path, workspaceId: workspaceId}

    return request({
        url: `/${apiPath}/get`,
        params
    });
}

export async function loadChildren(dir: string, workspaceId: number): Promise<any> {
    const params = {dir: dir, workspaceId: workspaceId}

    return request({
        url: `/${apiPath}/loadChildren`,
        params
    });
}

export async function extract(path: string, workspaceId: number): Promise<any> {
    const params = {path: path, workspaceId: workspaceId}

    return request({
        url: `/${apiPath}/extract`,
        params
    });
}

export async function create(params: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'POST',
        data: params,
    });
}

export async function update(id: number, params: any): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'PUT',
        data: params,
    });
}

export async function remove(id: number): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'delete',
    });
}

export async function updateCode(data: any): Promise<any> {
    return request({
        url: `/${apiPath}/updateCode`,
        method: 'PUT',
        data: data,
    });
}

export function genExecInfo(jsn: WsMsg) : WsMsg {
    jsn.msg = jsn.msg.replace(/^"+/,'').replace(/"+$/,'')
    return jsn
}

export async function syncFromZentao(params: any): Promise<any> {
    return request({
        url: `/${apiPath}/syncFromZentao`,
        method: 'POST',
        data: removeEmptyField(params),
    });
}

export async function syncToZentao(sets: any[]): Promise<any> {
    return request({
        url: `/${apiPath}/syncToZentao`,
        method: 'POST',
        data: sets
    });
}

export function genWorkspaceToScriptsMap(scripts: any[]): any[] {
    const workspaceIds = [] as number[]
    const mp = {}
    scripts.forEach((item) => {
        if (!mp[item.workspaceId]) {
            mp[item.workspaceId] = []
            workspaceIds.push(item.workspaceId)
        }

        mp[item.workspaceId].push(item.path)
    })

    const sets = [] as any[]
    workspaceIds.forEach((workspaceId) => {
        const set = {workspaceId: workspaceId, cases: mp[workspaceId]}
        sets.push(set)
    })

    return sets
}
