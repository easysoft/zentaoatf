import request from '@/utils/request';
import {Config} from "@/views/config/data";
import {SyncSettings} from "@/views/sync/data";
import {removeEmptyField} from "@/utils/object";

const apiPath = 'exec';

export async function execCase(scriptPaths: string[]): Promise<any> {
    const data = {cases: scriptPaths}
    return request({
        url: `/${apiPath}/execCase`,
        method: 'POST',
        data: data,
    });
}
