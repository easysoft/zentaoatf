import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import { Execution } from './data.d';
import { QueryResult, QueryParams } from '@/types/data.d';
import {
    list, remove, create, detail, update,
} from './service';

export interface StateType {
    listResult: any[];
    detailResult: Partial<Execution>;
}

export interface ModuleType extends StoreModuleType<StateType> {
    state: StateType;
    mutations: {
        setList: Mutation<StateType>;
        setItem: Mutation<StateType>;
    };
    actions: {
        listExecution: Action<StateType, StateType>;
        deleteExecution: Action<StateType, StateType>;
        getExecution: Action<StateType, StateType>;
    };
}
const initState: StateType = {
    listResult: [],
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
            state.listResult = payload;
        },
        setItem(state, payload) {
            state.detailResult = payload;
        },
    },
    actions: {
        async listExecution({ commit }, params: QueryParams ) {
            try {
                const response: ResponseData = await list(params);
                if (response.code != 0) return;
                const data = response.data;
                commit('setList', data);

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
        async deleteExecution({ commit }, payload: string ) {
            try {
                await remove(payload);
                await this.dispatch('ListExecution/listExecution', {})

                return true;
            } catch (error) {
                return false;
            }
        },
    }
};

export default StoreModel;
