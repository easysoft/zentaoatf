import request from '@/utils/request';

const apiPath = 'projects';

export async function queryProject(currProjectPath: string): Promise<any> {
    const params = {currProject: currProjectPath}

    return request({
        url: `/${apiPath}/getByUser`,
        method: 'GET',
        params,
    });
}

export async function createProject(data: string): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'post',
        data: data,
    });
}
