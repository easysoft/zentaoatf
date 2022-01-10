import request from '@/utils/request';

export async function queryProject(currProjectPath: string): Promise<any> {
    const params = {currProject: currProjectPath}

    return request({
        url: '/projects/listByUser',
        params,
    });
}