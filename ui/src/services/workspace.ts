import request from '@/utils/request';

const apiPath = 'workspaces';

export async function queryWorkspace(currWorkspacePath: string): Promise<any> {
    const params = {currWorkspace: currWorkspacePath}

    return request({
        url: `/${apiPath}/getByUser`,
        method: 'GET',
        params,
    });
}
export async function deleteWorkspace(currWorkspacePath: string): Promise<any> {
    const params = {path: currWorkspacePath}

    return request({
        url: `/${apiPath}`,
        method: 'DELETE',
        params,
    });
}

export async function createWorkspace(data: string): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'post',
        data: data,
    });
}
