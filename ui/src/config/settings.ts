import {getCmdHistories, getInitStatus, getOpenedScripts} from "@/utils/cache";

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

    eventExec: string;
    eventNotify: string;
    eventWebSocketConnStatus: string,
    eventWebSocketMsg: string,
    webSocketRoom: string,
    electronMsg: string,
    electronMsgReplay: string,

    initStatus: string;
    currSiteId: string;
    currProductId: string;
    currProductIdBySite: string;
    currWorkspace: string;
    displayBy: string;
    scriptFilters: string;
    expandedKeys: string;
    cmdHistories: string;
    openedScripts: string;

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

    eventExec: 'eventExec',
    eventNotify: 'eventNotify',
    eventWebSocketConnStatus: 'eventWebSocketStatus',
    eventWebSocketMsg: 'eventWebSocketMsg',
    webSocketRoom: 'webSocketRoom',
    electronMsg: 'electronMsg',
    electronMsgReplay: 'electronMsgReplay',

    initStatus: 'initStatus',
    currSiteId: 'currSiteId',
    currProductId: 'currProductId',
    currProductIdBySite: 'currProductIdBySite',
    currWorkspace: 'currWorkspace',
    scriptFilters: 'scriptFilters',
    displayBy: 'displayBy',
    expandedKeys: 'expandedKeys',
    cmdHistories: 'cmdHistories',
    openedScripts: 'openedScripts',

    ajaxHeadersTokenKey: 'Authorization',
    ajaxResponseNoVerifyUrl: [
        '/account/login', // 用户登录
        '/account/info', // 获取用户信息
    ],
    iconfontUrl: [],
};

export default settings;
  