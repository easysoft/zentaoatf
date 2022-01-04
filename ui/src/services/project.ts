import request from '@/utils/request';

export async function queryProject(): Promise<any> {
    return request({
        url: '/projects/listByUser'
    });
}