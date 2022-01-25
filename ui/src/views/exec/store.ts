import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { Execution } from './data.d';
import { QueryResult, QueryParams } from '@/types/data.d';
import {
    list, get, remove,
} from './service';

export interface StateType {
    items: any[];
    item: Partial<Execution>;
}

export interface ModuleType extends StoreModuleType<StateType> {
    state: StateType;
    mutations: {
        setItems: Mutation<StateType>;
        setItem: Mutation<StateType>;
    };
    actions: {
        list: Action<StateType, StateType>;
        get: Action<StateType, StateType>;
        delete: Action<StateType, StateType>;
    };
}
const initState: StateType = {
    items: [],
    item: {},
};

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'History',
    state: {
        ...initState
    },
    mutations: {
        setItems(state, payload) {
            state.items = payload;
        },
        setItem(state, payload) {
            state.item = payload;
        },
    },
    actions: {
        async list({ commit }, params: QueryParams ) {
            try {
                const response: ResponseData = await list(params);
                if (response.code != 0) return;
                const data = response.data;
                commit('setItems', data);

                return true;
            } catch (error) {
                return false;
            }
        },
        async get({ commit }, payload: string ) {
            try {
                const response: ResponseData = await get(payload);
                const { data } = response;
                commit('setItem',{
                    ...initState.item,
                    ...data,
                });
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
