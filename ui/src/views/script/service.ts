import request from '@/utils/request';
import {WsMsg} from "@/types/data";
import {removeEmptyField} from "@/utils/object";
import {momentTime} from "@/utils/datetime";
import {testToolMap} from "@/utils/testing";

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

export async function create(data: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'POST',
        data,
    });
}

export async function update(id: number, params: any): Promise<any> {
    return request({
        url: `/${apiPath}/${id}`,
        method: 'PUT',
        data: params,
    });
}
export async function updateNameReq(path: string, name: string): Promise<any> {
    const data = {path, name}

    return request({
        url: `/${apiPath}/updateName`,
        method: 'PUT',
        data,
    });
}

export async function move(data: any): Promise<any> {
    return request({
        url: `/${apiPath}/move`,
        method: 'put',
        data,
    });
}

export async function remove(path: string): Promise<any> {
    const data = {path: path}

    return request({
        url: `/${apiPath}`,
        method: 'delete',
        data,
    });
}

export async function updateCode(data: any): Promise<any> {
    return request({
        url: `/${apiPath}/updateCode`,
        method: 'PUT',
        data: data,
    });
}

export function genExecInfo(item: WsMsg, count: number) : WsMsg {
    if (item.info) item.info.key = item.info.key ? item.info.key + '-' + count : undefined

    if (item.info && item.info.status)  {
        item.msg = setFirstLineColor(item.msg, item.info.status)
    }

    item.msg = item.msg.replace(/^"+/,'').replace(/"+$/,'')
        .replaceAll('\n','<br />')
        .replaceAll('[','&nbsp;&nbsp;&nbsp;[')
    if (item.msg) item.time = momentTime(new Date());
    return item
}

function setFirstLineColor(msg, status) {
    const arr = msg.split('<br/>')
    arr[0] = `<span class="result-${status}">${arr[0]}</span>`
    return arr.join('<br/>')
}

export function getCaseIdsFromReport(workspace, seq, scope) {
    const params = {workspaceId: workspace, seq: seq, scope: scope}
    return request({
        url: `/${apiPath}/getCaseIdsFromReport`,
        method: 'get',
        params,
    });
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

export const getNodeMap = (node, mp): void => {
    if (!node) return

    mp[node.path] = node
    if (node.children) {
        node.children.forEach(c => {
            getNodeMap(c, mp)
        })
    }
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

export function getSyncFromInfoFromMenu(key: string, node: any): any {
    let moduleId = 0
    let caseId = 0
    let scriptPath = ''

    if (key.indexOf('zentao-') === 0) {
        const arr = key.split('-')
        if (arr[1] === 'module') {
            moduleId = parseInt(arr[2])
        } else if (arr[1] === 'case') {
            caseId = parseInt(arr[2])
        }
    } else {
        moduleId = node.moduleId
        caseId = node.caseId

        scriptPath = key
    }

    return {
        moduleId,
        caseId,
        scriptPath,
    }
}

export const getFileNodesUnderParent = (node): string[] => {
    console.log('getFileNodesUnderParent')

    const nodeMap = {}
    getNodeMap(node, nodeMap)

    const arr = [] as string[]
    Object.keys(nodeMap).forEach((k, v) => {
        const node = nodeMap[k]
        if (node.type === 'file') {
            node.childrem = null
            arr.push(node)
        }
    })

    return arr
}

export function getSyncToInfoFromMenu(key: string, node: any): any {

    return
}

export function scriptTreeAddAttr(treeData) {
    if(treeData == undefined){
        return treeData;
    }
    treeData = treeData.map((item, index) => {
        item.id = item.path;
        item.checkable = item.workspaceType == 'ztf' ? true : false;
        if (item.isLeaf) {
            item.toolbarItems = [
                // { hint: 'create_file', icon: 'file-add', key: 'createFile'},
            ];
        } else {
            item.toolbarItems = [
                {hint: 'create_workspace', icon: 'folder-add', key: 'createDir'},
                { hint: 'create_file', icon: 'file-add', key: 'createFile'},
            ];
            if(item.workspaceType != 'ztf'){
                item.toolbarItems.push({hint: testToolMap[item.workspaceType], icon: 'play', hintI18n: 'test', key: 'runTest'})
            }
        }

        if (item.type === "workspace") {
            item.toolbarItems.push({hint:'delete', icon:'delete', key: 'deleteWorkspace'})
        }
        if (item.children != undefined && item.children.length > 0) {
            item.children = scriptTreeAddAttr(item.children)
        }
        return item;
    })

    return treeData
}
