import request from '@/utils/request';
import {Config} from "@/views/config/data";

const apiPath = 'projects';

export async function queryProject(currProjectPath: string): Promise<any> {
    const params = {currProject: currProjectPath}

    return request({
        url: `/${apiPath}/listByUser`,
        method: 'GET',
        params,
    });
}

export async function saveConfig(params: Partial<Config>): Promise<any> {
    return request({
        url: `/${apiPath}/saveConfig`,
        method: 'POST',
        data: params,
    });
}