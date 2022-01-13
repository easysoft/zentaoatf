/**
 * 自定义 request 网络请求工具,基于axios
 * @author LiQingSong
 */
import axios, { AxiosPromise, AxiosRequestConfig, AxiosResponse } from 'axios';
import { notification } from "ant-design-vue";
import router from '@/config/routes';
import settings from '@/config/settings';
import { getCache, setCache } from '@/utils/localCache';

export interface ResponseData {
    code: number;
    data?: any;
    msg?: string;
    token?: string;
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
 * 异常处理程序
 */
const errorHandler = (error: any) => {

    const { response, message } = error;
    if (message === 'CustomError') {
        // 自定义错误
        const { config, data } = response;
        const { url, baseURL} = config;
        const { code, msg } = data;
        const reqUrl = url.split("?")[0].replace(baseURL, '');
        const noVerifyBool = settings.ajaxResponseNoVerifyUrl.includes(reqUrl);
        if (!noVerifyBool) {
            notification.error({
              message: `提示`,
              description: customCodeMessage[code] || msg || 'Error',
            });
        }
    } else if (message === 'CancelToken') {
        // 取消请求 Token
        // eslint-disable-next-line no-console
        console.log(message);
    } else if (response && response.status) {
        const errorText = serverCodeMessage[response.status] || response.statusText;
        const { status, request } = response;
        notification.error({
            message: `请求错误 ${status}: ${request.responseURL}`,
            description: errorText,
        });
    } else if (!response) {
        notification.error({
            message: '请求失败',
            description: '无法连接到服务器！',
        });
    }

    return Promise.reject(error);
}

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
 * 请求前
 * 请求拦截器
 */
request.interceptors.request.use(
    async (config: AxiosRequestConfig & { cType?: boolean, baseURL?: string }) => {

        // 如果设置了cType 说明是自定义 添加 Content-Type类型 为自定义post 做铺垫
        if (config['cType']) {
            config.headers['Content-Type'] = 'application/x-www-form-urlencoded;charset=UTF-8';
        }

        // 加随机数清除缓存
        const projectPath = await getCache(settings.currProject);
        config.params = { ...config.params, currProject: projectPath, ts: Date.now() };

        console.log('=== request ===', config)
        return config;
    },
    /* error=> {} */ // 已在 export default catch
);

/**
 * 请求后
 * 响应拦截器
 */
request.interceptors.response.use(
    async (response: AxiosResponse) => {
        console.log('=== response ===', response)

        const res: ResponseData = response.data;
        const { code, token } = res;

        // 自定义状态码验证
        if (code !== 0) {
            return Promise.reject({
                response,
                message: 'CustomError',
            });
        }

        return response;
    },
    /* error => {} */ // 已在 export default catch
);

/** 
 * ajax 导出
 * 
 * Method: get
 *     Request Headers
 *         无 - Content-Type
 *     Query String Parameters
 *         name: name
 *         age: age
 * 
 * Method: post
 *     Request Headers
 *         Content-Type:application/json;charset=UTF-8
 *     Request Payload
 *         { name: name, age: age }
 *         Custom config parameters
 *             { cType: true }  Mandatory Settings Content-Type:application/json;charset=UTF-8
 * ......
 */
export default function(config: AxiosRequestConfig): AxiosPromise<any> {
    return request(config).then((response: AxiosResponse) => response.data).catch(error => errorHandler(error));
}
