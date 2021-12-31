import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { Script, QueryResult, QueryParams, PaginationConfig } from './data.d';
import {
    query, remove, create, detail, update,
} from './service';

export interface StateType {
    queryResult: QueryResult;
    detailResult: Partial<Script>;
}

export interface ModuleType extends StoreModuleType<StateType> {
    state: StateType;
    mutations: {
        setList: Mutation<StateType>;
        setItem: Mutation<StateType>;
    };
    actions: {
        queryScript: Action<StateType, StateType>;
        deleteScript: Action<StateType, StateType>;
        createScript: Action<StateType, StateType>;
        getScript: Action<StateType, StateType>;
        updateScript: Action<StateType, StateType>;
    };
}
const initState: StateType = {
    queryResult: {
        list: [],
        pagination: {
            total: 0,
            current: 1,
            pageSize: 10,
            showSizeChanger: true,
            showQuickJumper: true,
        },
    },
    detailResult: {},
};

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'ListScript',
    state: {
        ...initState
    },
    mutations: {
        setList(state, payload) {
            state.queryResult = payload;
        },
        setItem(state, payload) {
            state.detailResult = payload;
        },
    },
    actions: {
        async queryScript({ commit }, params: QueryParams ) {
            try {
                const response: ResponseData = await query(params);
                if (response.code != 0) return;

                const data = response.data;

                commit('setList',{
                    ...initState.queryResult,
                    list: data.result || [],
                    pagination: {
                        ...initState.queryResult.pagination,
                        current: params.page,
                        pageSize: params.pageSize,
                        total: data.total || 0,
                    },
                });
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
        async createScript({ commit }, payload: Pick<Script, "name" | "desc"> ) {
            try {
                await create(payload);
                return true;
            } catch (error) {
                return false;
            }
        },
        async getScript({ commit }, payload: number ) {
            try {
                const response: ResponseData = await detail(payload);
                const { data } = response;
                commit('setItem',{
                    ...initState.detailResult,
                    ...data,
                });
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
    }
};

export default StoreModel;
