import request from '@/utils/request';
import {Config} from "@/views/config/data";
import {SyncSettings} from "@/views/sync/data";

const apiPath = 'sync';

export async function syncFromZentao(params: SyncSettings): Promise<any> {
    return request({
        url: `/${apiPath}/syncFromZentao`,
        method: 'POST',
        data: params,
    });
}

export async function syncToZentao(): Promise<any> {
    return request({
        url: `/${apiPath}/syncToZentao`,
        method: 'POST',
    });
}
