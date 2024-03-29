import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import {WebSocket} from "@/services/websocket";
import settings from "@/config/settings";

import {
    list,
    get,
    extract,
    create,
    update,
    remove,
    loadChildren,
    updateCode,
    syncFromZentao,
    syncToZentao,
    paste, move,
    rename,
    scriptTreeAddAttr,
    getNodeMap,
} from './service';
import {ScriptFileNotExist} from "@/utils/const";
import { jsonStrDef } from '@/utils/dom';

export interface ScriptData {
    list: [];
    detail: any;
    treeDataMap: any;
    checkedNodes: [];

    currWorkspace: any
    queryParams: any;
    currentCodeChanged: boolean;
    watchPaths: string;
}

export interface ModuleType extends StoreModuleType<ScriptData> {
    state: ScriptData;
    mutations: {
        setList: Mutation<ScriptData>;
        setItem: Mutation<ScriptData>;
        setWorkspace: Mutation<ScriptData>;
        setQueryParams: Mutation<ScriptData>;
        setCheckedNodes: Mutation<ScriptData>;
        setCurrentCodeChanged: Mutation<ScriptData>;
        setWatchPath: Mutation<ScriptData>;
    };
    actions: {
        listScript: Action<ScriptData, ScriptData>;
        getScript: Action<ScriptData, ScriptData>;
        loadChildren: Action<ScriptData, ScriptData>;
        syncFromZentao: Action<ScriptData, ScriptData>;
        syncToZentao: Action<ScriptData, ScriptData>;
        extractScript: Action<ScriptData, ScriptData>;
        changeWorkspace: Action<ScriptData, ScriptData>;
        setCheckedNodes: Action<ScriptData, ScriptData>;

        createScript: Action<ScriptData, ScriptData>;
        updateScript: Action<ScriptData, ScriptData>;
        deleteScript: Action<ScriptData, ScriptData>;
        renameScript: Action<ScriptData, ScriptData>;
        pasteScript: Action<ScriptData, ScriptData>;
        moveScript: Action<ScriptData, ScriptData>;
        updateCode: Action<ScriptData, ScriptData>;
        updateCurrentCodeChanged: Action<ScriptData, ScriptData>;
    };
}
const initState: ScriptData = {
    list: [],
    detail: null,
    treeDataMap: {},

    currWorkspace: {id: 0, type: 'ztf'},
    queryParams: {},

    checkedNodes: [],
    currentCodeChanged: false,
    watchPaths: '',
};

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'Script',
    state: {
        ...initState
    },
    mutations: {
        setList(state, payload) {
            state.list = payload.length > 0 ? payload[0]['children'] : [];
            const treeDataMap = {};
            if(payload.length > 0) {
                getNodeMap(payload[0], treeDataMap);
            }
            state.treeDataMap = treeDataMap;
        },
        setItem(state, payload) {
            state.detail = payload;
        },
        setWorkspace(state, payload) {
            state.currWorkspace = payload;
        },
        setQueryParams(state, payload) {
            state.queryParams = payload;
        },

        setCheckedNodes(state, payload) {
            state.checkedNodes = payload;
        },
        setCurrentCodeChanged(state, payload) {
            state.currentCodeChanged = payload;
        },
        setWatchPath(state, payload) {
            state.watchPaths = payload;
        },
    },
    actions: {
        async listScript({ commit, state }, playload: any ) {
            const response: ResponseData = await list(playload);
            const data = response.data;
            data.id = data.path;
            data.children = scriptTreeAddAttr(data.children ? data.children : []);
            commit('setList', [data]);
            if(data.children != undefined){
                const watchPaths = [] as any;
                data.children.forEach(element => {
                    watchPaths.push({WorkspacePath: element.path})
                });
                if(state.watchPaths != JSON.stringify(watchPaths)){
                    console.log("watchPaths", state.watchPaths, JSON.stringify(watchPaths))
                    WebSocket.sentMsg(settings.webSocketRoom, JSON.stringify({act: 'watch',testSets:watchPaths}), "local")
                    commit('setWatchPath', JSON.stringify(watchPaths));
                }
            }
            commit('setQueryParams', playload);
            return true;
        },
        async loadChildren(_ctx, treeNode: any ) {
            console.log('load node children', treeNode.dataRef.workspaceType)
            if (treeNode.dataRef.workspaceType === 'ztf')
                return true

            loadChildren(treeNode.dataRef.path, treeNode.dataRef.workspaceId).then((json) => {
                treeNode.dataRef.children = json.data
                return true;
            })
        },

        async getScript({ commit}, script: any ) {
            if (!script || script.type !== 'file') {
                commit('setItem', null);
            } else if (script.path.indexOf('zentao-') === 0) {
                commit('setItem', {id: script.caseId, workspaceId: script.workspaceId, code: ScriptFileNotExist});
            } else {
                const response: ResponseData = await get(script.path, script.workspaceId);
                commit('setItem', response.data);
            }

            return true;
        },

        async syncFromZentao({ commit, dispatch, state }, payload: any ) {
            const resp = await syncFromZentao(payload)
            if (resp.code === 0) {
                await dispatch('listScript', state.queryParams)

                if (resp.code === 0 && resp.data != null && resp.data.length === 1) {
                    const getResp = await get(resp.data[0], payload.workspaceId);
                    commit('setItem', getResp.data);
                } else {
                    commit('setItem', null);
                }
            }

            return resp
        },

        async syncToZentao({ dispatch, state }, payload: any ) {
            const resp = await syncToZentao(payload)

            if (resp.code === 0) {
                await dispatch('listScript', state.queryParams)
            }

            return resp
        },

        async extractScript({ commit }, script: any ) {
            if (!script.path) return true

            const response: ResponseData = await extract(script.path, script.workspaceId)
            const { data } = response
            commit('setItem', data.script)

            return data.done
        },

        async createScript({ dispatch, state}, payload: any) {
            try {
                const jsn = await create(payload);
                const path = jsn.data

                await dispatch('listScript', state.queryParams)

                return path;
            } catch (error) {
                return ''
            }
        },
        async updateScript(_ctx, payload: any ) {
            try {
                const { id, ...params } = payload;
                await update(id, { ...params });
                return true;
            } catch (error) {
                return false;
            }
        },

        async updateCode({ dispatch, state }, payload: any ) {
            try {
                await updateCode(payload);
                dispatch('listScript', state.queryParams)
                return true;
            } catch (error) {
                return false;
            }
        },

        async pasteScript({ dispatch, state}, data: any ) {
            try {
                await paste(data);
                await dispatch('listScript', state.queryParams)

                return true;
            } catch (error) {
                return false;
            }
        },

        async moveScript({ dispatch, state}, data: any ) {
            try {
                await move(data);
                await dispatch('listScript', state.queryParams)

                return true;
            } catch (error) {
                return false;
            }
        },

        async deleteScript({ dispatch, state}, path: string ) {
            try {
                await remove(path);
                await dispatch('listScript', state.queryParams)

                return true;
            } catch (error) {
                return false;
            }
        },

        async renameScript({ dispatch, state}, data: any ) {
            try {
                await rename(data);
                await dispatch('listScript', state.queryParams)

                return true;
            } catch (error) {
                return false;
            }
        },

        async changeWorkspace({ commit }, payload: any ) {
            commit('setWorkspace', payload);
            return true;
        },

        async setCheckedNodes({ commit }, payload: any ) {
            commit('setCheckedNodes', payload);
            return true;
        },

        async updateCurrentCodeChanged({ commit }, payload: any ) {
            commit('setCurrentCodeChanged', payload);
            return true;
        },
    }
};

export default StoreModel;