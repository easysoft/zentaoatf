import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { QueryParams, QueryResult } from '@/types/data.d';
import {
    query, listByProduct, get, save, remove,
} from './service';

export interface WorkspaceData {
    queryResult: QueryResult;
    listResult: any[];
    detailResult: any;
}

export interface ModuleType extends StoreModuleType<WorkspaceData> {
    state: WorkspaceData;
    mutations: {
        setQueryResult: Mutation<WorkspaceData>;
        setListResult: Mutation<WorkspaceData>;
        setDetailResult: Mutation<WorkspaceData>;
    };
    actions: {
        query: Action<WorkspaceData, WorkspaceData>;
        list: Action<WorkspaceData, WorkspaceData>;
        get: Action<WorkspaceData, WorkspaceData>;
        save: Action<WorkspaceData, WorkspaceData>;
        delete: Action<WorkspaceData, WorkspaceData>;
    };
}
const initState: WorkspaceData = {
    queryResult: {
        result: [],
        pagination: {
            total: 0,
            page: 1,
            pageSize: 10,
            showSizeChanger: true,
            showQuickJumper: true,
        },
    },
    listResult: [],
    detailResult: {},
};

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'Workspace',
    state: {
        ...initState
    },
    mutations: {
        setQueryResult(state, payload) {
            state.queryResult = payload;
        },
        setListResult(state, payload) {
            state.listResult = payload;
        },
        setDetailResult(state, payload) {
            state.detailResult = payload;
        },
    },
    actions: {
        async query({ commit }, params: QueryParams ) {
            try {
                const response: ResponseData = await query(params);
                if (response.code != 0) {
                    return;
                }
                const data = response.data;
                commit('setQueryResult', data);

                return true;
            } catch (error) {
                return false;
            }
        },
        async list({ commit }, productId: number ) {
            try {
                const response: ResponseData = await listByProduct(productId);
                if (response.code != 0) {
                    return;
                }
                const data = response.data;
                commit('setListResult', data);

                return true;
            } catch (error) {
                return false;
            }
        },
        async get({ commit }, id: number ) {
            let data = {name:'', path: '', type: 'ztf'}
            if (id) {
                const response: ResponseData = await get(id);
                data = response.data;
            }
            commit('setDetailResult',data);
            return true;
        },
        async save({ commit }, payload) {
            try {
                await save(payload);
                return true;
            } catch (error) {
                return false;
            }
        },
        async delete({ commit }, id: number ) {
            try {
                await remove(id);
                await this.dispatch('Workspace/list', {})

                return true;
            } catch (error) {
                return false;
            }
        },
    }
};

export default StoreModel;
