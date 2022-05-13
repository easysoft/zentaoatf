import { Getter, Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";

/**
 * Page type（Test script, Test result, settings and sites）
 * 标签页页面类型（测试脚步代码、测试结果、设置、站点管理）
 */
export type PageType = 'script' | 'result' | 'settings' | 'sites';

/**
 * Single page tab
 * 单独的页面标签页定义
 */
export interface PageTab {
    /**
     * Page unique identifier
     * ID，用于标识唯一的标签页页面
     */
    id: string;

    /**
     * Page title
     * 页面标题
     */
    title: string;

    /**
     * Page icon
     * 页面图标
     */
    icon?: string;

    /**
     * Page type
     * 标签页类型
     */
    type: PageType;

    /**
     * Whether or not has changes not saved
     * 是否有修改没有保存到磁盘
     */
    changed: boolean;

    /**
     * Whether is readonly
     * 是否只读
     */
    readonly: boolean;

    /**
     * Last Active time
     * 上次打开的时间
     */
    lastActiveTime: number;

    /**
     * Other data
     * 其他数据
     */
    data: Record<string, unknown>;
}

/**
 * Tabs module state data
 * 标签页模块状态数据
 */
export interface TabsData {
    /**
     * All tabs map
     * 所有已经打开的标签页表
     */
    map: {[id: string]: PageTab};

    /**
     * All tabs id list
     * 所有已经打开的标签页 ID 列表（维持显示顺序）
     */
    idList: string[];

    /**
     * Current active tab id
     * 当前激活的标签页 ID
     */
    activeID: string;
}

/**
 * Tabs module type
 * 标签页存储模块类型
 */
export interface TabsModuleType extends StoreModuleType<TabsData> {
    state: TabsData;
    getters: {
        count: Getter<TabsData, TabsData>,
        list: Getter<TabsData, TabsData>,
        currentTab: Getter<TabsData, TabsData>,
    };
    mutations: {
        activate: Mutation<TabsData>,
        activateLast: Mutation<TabsData>,
        addTab: Mutation<TabsData>,
        removeTab: Mutation<TabsData>,
        updateTab: Mutation<TabsData>,
    };
    actions: {
        open: Action<TabsData, TabsData>,
        close: Action<TabsData, TabsData>,
        hide: Action<TabsData, TabsData>,
        update: Action<TabsData, TabsData>,
    };
}

const initState: TabsData = {
    map: {},
    idList: [],
    activeID: '',
};

type PageTabInfo = Partial<PageTab> & {id: string};

const TabsModel: TabsModuleType = {
    namespaced: true,
    name: 'tabs',
    state: initState,
    getters: {
        count(state) {
            return state.idList.length;
        },
        list(state) {
            return state.idList.map(x => state.map[x]);
        },
        currentTab(state) {
            return state.map[state.activeID];
        }
    },
    mutations: {
        activate(state, payload: PageTabInfo) {
            const {id} = payload;
            const tab = state.map[id];
            if (tab) {
                state.activeID = id;
                Object.assign(tab, {lastActiveTime: Date.now()}, payload);
            }
        },
        activateLast(state, payload: {exceptID?: string}) {
            if (!state.idList.length) {
                return;
            }

            let lastTab: PageTab | undefined;
            Object.values(state.map).forEach(tab => {
                if (payload && payload.exceptID === tab.id) {
                    return;
                }
                if (!lastTab || lastTab.lastActiveTime < tab.lastActiveTime) {
                    lastTab = tab;
                }
            });

            if (lastTab) {
                state.activeID = lastTab.id;
                lastTab.lastActiveTime = Date.now();
            }
        },
        addTab(state, payload: PageTab) {
            const {id} = payload;
            state.map[id] = payload;
            if (!state.idList.includes(id)) {
                state.idList.push(id);
            }
            const lastActiveTab = state.map[state.activeID];
            if (payload.lastActiveTime && (!lastActiveTab || lastActiveTab.lastActiveTime < payload.lastActiveTime)) {
                state.activeID = payload.id;
            }
        },
        removeTab(state, payload: {id: string}) {
            const {id} = payload;
            delete state.map[id];
            const index = state.idList.findIndex(x => x === id);
            if (index >= 0) {
                state.idList.splice(index, 1);
            }
        },
        updateTab(state, payload: PageTab) {
            const {id} = payload;
            const tab = state.map[id];
            if (!tab) {
                return;
            }
            Object.assign(tab, payload);
        }
    },
    actions: {
        open(context, payload: PageTab) {
            const {state} = context;
            const {id} = payload;
            const tab = state.map[id];
            if (tab) {
                context.commit('activate', payload);
            } else {
                payload.lastActiveTime = Date.now();
                context.commit('addTab', payload);
            }
        },
        close(context, payload: {id: string}) {
            const {id} = payload;
            context.commit('removeTab', payload);
            context.commit('activateLast');
        },
        hide(context, payload: {id: string}) {
            context.commit('activateLast', {exceptID: payload.id});
        },
        update(context, payload: PageTab) {
            const {state} = context;
            const {id} = payload;
            context.commit('updateTab', payload);
            }
    }
};

export default TabsModel;
