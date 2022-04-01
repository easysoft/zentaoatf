import request from '@/utils/request';
import {WsMsg} from "@/views/result/data";

const apiPath = 'reports';

export function getCaseIdsFromReport(workspace, seq, scope) {
    const params = {workspaceId: workspace, seq: seq, scope: scope}
    return request({
        url: `/${apiPath}/getCaseIdsFromReport`,
        method: 'get',
        params,
    });
}
