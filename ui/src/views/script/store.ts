import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { Script } from './data.d';
import {
    list, get, extract, create, update, remove
} from './service';

export interface ScriptData {
    list: [];
    detail: Partial<Script>;
}

export interface ModuleType extends StoreModuleType<ScriptData> {
    state: ScriptData;
    mutations: {
        setList: Mutation<ScriptData>;
        setItem: Mutation<ScriptData>;
    };
    actions: {
        listScript: Action<ScriptData, ScriptData>;
        getScript: Action<ScriptData, ScriptData>;
        extractScript: Action<ScriptData, ScriptData>;

        createScript: Action<ScriptData, ScriptData>;
        updateScript: Action<ScriptData, ScriptData>;
        deleteScript: Action<ScriptData, ScriptData>;
    };
}
const initState: ScriptData = {
    list: [],
    detail: {},
};

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'script',
    state: {
        ...initState
    },
    mutations: {
        setList(state, payload) {
            state.list = payload;
        },
        setItem(state, payload) {
            state.detail = payload;
        },
    },
    actions: {
        async listScript({ commit }, playload: any ) {
            const response: ResponseData = await list(playload);
            const { data } = response;
            commit('setList', [data]);
            return true;
        },

        async getScript({ commit }, script: any ) {
            if (script.isDir) {
                commit('setItem', {});
                return true;
            }

            const response: ResponseData = await get(script.path);
            const { data } = response;
            commit('setItem',data);
            return true;
        },
        async extractScript({ commit }, script: any ) {
            if (script.isDir) return true

            const response: ResponseData = await extract(script.path)
            const { data } = response
            commit('setItem', data)
            return true
        },

        async createScript({ commit }, payload: Pick<Script, "name" | "desc"> ) {
            try {
                await create(payload);
                return true;
            } catch (error) {
                return false;
            }
        },
        async updateScript({ commit }, payload: Script ) {
            try {
                const { id, ...params } = payload;
                await update(id, { ...params });
                return true;
            } catch (error) {
                return false;
            }
        },

        async deleteScript({ commit }, payload: number ) {
            try {
                await remove(payload);
                return true;
            } catch (error) {
                return false;
            }
        },
    }
};

export default StoreModel;
