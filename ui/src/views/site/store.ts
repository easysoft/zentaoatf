import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { QueryParams, QueryResult } from '@/types/data.d';
import {
    list, get, save, remove,
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
        save: Action<StateType, StateType>;
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
        async get({ commit }, id: number ) {
            let data = {name:'', url: '', username: '', password: ''}
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
                await this.dispatch('Site/list', {})

                return true;
            } catch (error) {
                return false;
            }
        },
    }
};

export default StoreModel;
