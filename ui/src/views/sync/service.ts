import request from '@/utils/request';
import {SyncSettings} from "@/views/sync/data";
import {removeEmptyField} from "@/utils/object";

const apiPath = 'sync';

export async function syncFromZentao(settings: SyncSettings): Promise<any> {
    return request({
        url: `/${apiPath}/syncFromZentao`,
        method: 'POST',
    });
}

export async function syncToZentao(): Promise<any> {
    const params = {}
    return request({
        url: `/${apiPath}/syncToZentao`,
        method: 'POST',
        params
    });
}
