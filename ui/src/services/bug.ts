import request from '@/utils/request';
import {WsMsg} from "@/views/result/data";

const apiPath = 'bugs';

export function submitBugToZentao(data: any): Promise<any> {
    return request({
        url: `/${apiPath}`,
        method: 'post',
        data: data,
    });
}
