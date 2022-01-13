import request from '@/utils/request';
import {Config} from "@/views/config/data";
import {SyncSettings} from "@/views/sync/data";
import {removeEmptyField} from "@/utils/object";

const apiPath = 'sync';

export async function syncFromZentao(params: SyncSettings): Promise<any> {
    return request({
        url: `/${apiPath}/syncFromZentao`,
        method: 'POST',
        data: removeEmptyField(params),
    });
}

export async function syncToZentao(): Promise<any> {
    return request({
        url: `/${apiPath}/syncToZentao`,
        method: 'POST',
    });
}
