import request from '@/utils/request';
import {Config} from "@/views/config/data";
import {SyncSettings} from "@/views/sync/data";
import {removeEmptyField} from "@/utils/object";

const apiPath = 'exec';

export async function execCase(params: string[]): Promise<any> {
    return request({
        url: `/${apiPath}/execCase`,
        method: 'POST',
        params
    });
}
