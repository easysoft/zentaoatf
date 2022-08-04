import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { QueryParams, QueryResult } from '@/types/data.d';
import {
    list, get, remove,getStatistic,
} from './service';

export interface StateType {
    queryResult: QueryResult;
    detailResult: any;
    statistic: any;
}

export interface ModuleType extends StoreModuleType<StateType> {
    state: StateType;
    mutations: {
        setQueryResult: Mutation<StateType>;
        setStatistic: Mutation<StateType>;
        setDetailResult: Mutation<StateType>;
    };
    actions: {
        list: Action<StateType, StateType>;
        get: Action<StateType, StateType>;
        getStatistic: Action<StateType, StateType>;
        delete: Action<StateType, StateType>;
    };
}
const initState: StateType = {
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
    detailResult: {},
    statistic: {},

};

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'Result',
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
        setStatistic(state, payload) {
            state.statistic = payload;
        },
    },
    actions: {
        async list({ commit }, params: QueryParams ) {
            try {
                const response: ResponseData = await list(params);
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
        async get({ commit }, params: any ) {
            const response: ResponseData = await get(params);
            const data = response.data;

            commit('setDetailResult',data);

            return true;
        },
        async delete({ dispatch}, data: any ) {
            try {
                await remove(data);
                dispatch('list', {})

                return true;
            } catch (error) {
                return false;
            }
        },
        async getStatistic({ commit, dispatch }, script: any ) {
            const response: ResponseData = await getStatistic(script.path);
            commit('setStatistic', response.data);
            return true;
        },
    }
};

export default StoreModel;
