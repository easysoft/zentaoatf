/**
 * 自定义 request 网络请求工具,基于axios
 * @author LiQingSong
 */
import axios, { AxiosPromise, AxiosRequestConfig, AxiosResponse } from 'axios';
import settings from '@/config/settings';
import { getCache, setCache } from '@/utils/localCache';
import i18n from "@/config/i18n";
import {getCurrProductIdBySite, getCurrSiteId} from "@/utils/cache";
import bus from "@/utils/eventBus";

export interface ResponseData {
    code: number;
    data?: any;
    msg?: string;
    token?: string;
}
export interface ResultErr {
    httpCode: number;
    resultCode: number;
    resultMsg: string;
}

const customCodeMessage: {[key: number]: string} = {
  10002: '当前用户登入信息已失效，请重新登入再操作', // 未登陆
};

const serverCodeMessage: {[key: number]: string} = {
  200: '服务器成功返回请求的数据',
  400: 'Bad Request',
  401: 'Unauthorized',
  403: 'Forbidden',
  404: 'Not Found',
  500: '服务器发生错误，请检查服务器(Internal Server Error)',
  502: '网关错误(Bad Gateway)',
  503: '服务不可用，服务器暂时过载或维护(Service Unavailable)',
  504: '网关超时(Gateway Timeout)',
};

/**
 * 配置request请求时的默认参数
 */
const request = axios.create({
    baseURL: process.env.VUE_APP_APIHOST, // url = api url + request url
    withCredentials: true, // 当跨域请求时发送cookie
    timeout: 0 // 请求超时时间,5000(单位毫秒) / 0 不做限制
});

// 全局设置 - post请求头
// request.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded;charset=UTF-8';

/**
 * 请求拦截器
 */
request.interceptors.request.use(
    async (config: AxiosRequestConfig & { cType?: boolean, baseURL?: string }) => {

        // 如果设置了cType 说明是自定义 添加 Content-Type类型 为自定义post 做铺垫
        if (config['cType']) {
            config.headers['Content-Type'] = 'application/x-www-form-urlencoded;charset=UTF-8';
        }

        config.params = { ...config.params, ts: Date.now() };
        if (!config.params[settings.currSiteId]) {
            const currSiteId = await getCurrSiteId();
            config.params = { ...config.params, [settings.currSiteId]: currSiteId, lang: i18n.global.locale.value };
        }
        if (!config.params[settings.currProductId]) {
            const currProductId = await getCurrProductIdBySite(config.params[settings.currSiteId])
            config.params = { ...config.params, [settings.currProductId]: currProductId, lang: i18n.global.locale.value };
        }

        console.log(`currSiteId=${config.params[settings.currSiteId]}, currProductId=${config.params[settings.currProductId]}`)

        // if (!config.params[settings.currWorkspace]) {
        //     const workspacePath = await getCache(settings.currWorkspace);
        //     config.params = { ...config.params, [settings.currWorkspace]: workspacePath, lang: i18n.global.locale.value };
        // }

        console.log('=== request ===', config.url, config)
        return config;
    },
    /* error=> {} */ // 已在 export default catch
);

/**
 * 响应拦截器
 */
request.interceptors.response.use(
    async (axiosResponse: AxiosResponse) => {
        console.log('=== response ===', axiosResponse.config.url, axiosResponse)

        const data: ResponseData = axiosResponse.data;
        const { code, msg } = data;

        // 自定义状态码验证
        if (code !== 0) {
            return Promise.reject(axiosResponse);
        }

        return axiosResponse;
    }
);

/**
 * 异常处理
 */
const errorHandler = (resp: any) => {
    console.log(resp)
    if (!resp) resp = {status: 500}

    if (resp.status !== 200) {
        bus.emit(settings.eventNotify, {httpCode: resp.status})
        return
    }

    const result ={httpCode: resp.status, resultCode: resp.data.code, resultMsg: resp.data.msg} as ResultErr
    bus.emit(settings.eventNotify, result)

    return Promise.reject({})
}

export default function(config: AxiosRequestConfig): AxiosPromise<any> {
    return request(config).then((response: AxiosResponse) => response.data).catch(error => errorHandler(error));
}
