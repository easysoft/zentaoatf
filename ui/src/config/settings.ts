import { RoutesDataItem } from "@/utils/routes";

/**
 * 站点配置
 */
export interface SettingsType {
    /**
     * 站点名称
     */
    siteTitle: string;

    /**
     * 顶部菜单开启
     */
    topNavEnable: boolean;

    /**
     * 头部固定开启
     */
    headFixed: boolean;

    /**
     * 站点本地存储Token的Key值
     */
    siteTokenKey: string;

    /**
     * 站点本地存储当前工作目录的Key值
     */
    currWorkspace: string;

    /**
     * Ajax请求头发送Token 的 Key值
     */
    ajaxHeadersTokenKey: string;

    /**
     * Ajax返回值不参加统一验证的api地址
     */
    ajaxResponseNoVerifyUrl: string[];

    /**
     * iconfont.cn 项目在线生成的 js 地址
     */
    iconfontUrl: string[];
}

const settings: SettingsType = {
    siteTitle: 'ZTF',
    topNavEnable: true,
    headFixed: true,
    siteTokenKey: 'admin_antd_vue_token',
    currWorkspace: 'currWorkspace',
    // ajaxHeadersTokenKey: 'x-token',
    ajaxHeadersTokenKey: 'Authorization',
    ajaxResponseNoVerifyUrl: [
        '/account/login', // 用户登录
        '/account/info', // 获取用户信息
    ],
    iconfontUrl: [],
};

export default settings;
  