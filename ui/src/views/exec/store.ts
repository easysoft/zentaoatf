import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { Execution } from './data.d';
import { QueryResult, QueryParams } from '@/types/data.d';
import {
    query, remove, create, detail, update,
} from './service';

export interface StateType {
    queryResult: QueryResult;
    detailResult: Partial<Execution>;
}

export interface ModuleType extends StoreModuleType<StateType> {
    state: StateType;
    mutations: {
        setList: Mutation<StateType>;
        setItem: Mutation<StateType>;
    };
    actions: {
        queryExecution: Action<StateType, StateType>;
        deleteExecution: Action<StateType, StateType>;
        createExecution: Action<StateType, StateType>;
        getExecution: Action<StateType, StateType>;
        updateExecution: Action<StateType, StateType>;
    };
}
const initState: StateType = {
    queryResult: {
        data: [],
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
    name: 'ListExecution',
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
        async queryExecution({ commit }, params: QueryParams ) {
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
        async deleteExecution({ commit }, payload: number ) {
            try {
                await remove(payload);
                return true;
            } catch (error) {
                return false;
            }
        },
        async createExecution({ commit }, payload: Pick<Execution, "name" | "desc"> ) {
            try {
                await create(payload);
                return true;
            } catch (error) {
                return false;
            }
        },
        async getExecution({ commit }, payload: number ) {
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
        async updateExecution({ commit }, payload: Execution ) {
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
