import request from '@/utils/request';
import {LoginParamsType} from "@/views/user/login/data";

export async function queryCurrent(): Promise<any> {
    return request({
        url: '/users/profile',
        method: 'get'
    });
}

export async function queryMessage(): Promise<any> {
    return request({
        url: '/users/message'
    });
}