import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { QueryParams, QueryResult } from '@/types/data.d';
import {
    list, get, remove,
} from './service';

export interface StateType {
    queryResult: QueryResult;
    detailResult: any;
}

export interface ModuleType extends StoreModuleType<StateType> {
    state: StateType;
    mutations: {
        setQueryResult: Mutation<StateType>;
        setDetailResult: Mutation<StateType>;
    };
    actions: {
        list: Action<StateType, StateType>;
        get: Action<StateType, StateType>;
        delete: Action<StateType, StateType>;
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
    name: 'Site',
    state: {
        ...initState
    },
    mutations: {
        setQueryResult(state, payload) {
            state.queryResult = payload;
        },
        setDetailResult(state, payload) {
            state.detailResult = payload;
        },
    },
    actions: {
        async list({ commit }, params: QueryParams ) {
            try {
                const response: ResponseData = await list(params);
                if (response.code != 0) return;
                const data = response.data;
                commit('setQueryResult', data);

                return true;
            } catch (error) {
                return false;
            }
        },
        async get({ commit }, payload: string ) {
            try {
                const response: ResponseData = await get(payload);
                const { data } = response;
                commit('setDetailResult',data);
                return true;
            } catch (error) {
                return false;
            }
        },
        async delete({ commit }, payload: string ) {
            try {
                await remove(payload);
                await this.dispatch('History/list', {})

                return true;
            } catch (error) {
                return false;
            }
        },
    }
};

export default StoreModel;
