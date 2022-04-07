import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import {queryLang, querySiteAndProduct, getProfile, queryProduct, queryModule, querySuite, queryTask} from "../services/zentao";

import {setCurrProductIdBySite, setCurrSiteId} from "@/utils/cache";

export interface ZentaoData {
    langs: any[]

    profile: any

    sites: any[]
    products: any[]
    currSite: any
    currProduct: any

    modules: any[]
    suites: any[]
    tasks: any[]
}

export interface ModuleType extends StoreModuleType<ZentaoData> {
    state: ZentaoData;
    mutations: {
        saveLangs: Mutation<any>;

        saveProfile: Mutation<any>;

        saveSitesAndProduct: Mutation<any>;
        saveProducts: Mutation<any>;
        saveModules: Mutation<any>;
        saveSuites: Mutation<any>;
        saveTasks: Mutation<any>;
    };
    actions: {
        fetchLangs: Action<ZentaoData, ZentaoData>;
        getProfile: Action<ZentaoData, ZentaoData>;

        fetchSitesAndProduct: Action<ZentaoData, ZentaoData>;
        fetchProducts: Action<ZentaoData, ZentaoData>;
        fetchModules: Action<ZentaoData, ZentaoData>;
        fetchSuites: Action<ZentaoData, ZentaoData>;
        fetchTasks: Action<ZentaoData, ZentaoData>;
    };
}

const initState: ZentaoData = {
    langs: [],

    profile: {},

    sites: [],
    products: [],
    currSite: {},
    currProduct: {},

    modules: [],
    suites: [],
    tasks: [],
}

const StoreModel: ModuleType = {
    namespaced: true,
    name: 'Zentao',
    state: {
        ...initState
    },
    mutations: {
        saveLangs(state, payload) {
            console.log('payload', payload)
            state.langs = payload
        },
        saveProfile(state, payload) {
            console.log('payload', payload)
            state.profile = payload
        },
        async saveSitesAndProduct(state, payload) {
            console.log('saveSitesAndProduct', payload)
            if (!payload.currSite || !payload.currProduct) return

            state.sites = payload.sites;
            state.products = payload.products;

            // cache current site and product
            await setCurrSiteId(payload.currSite.id);
            await setCurrProductIdBySite(payload.currSite.id, payload.currProduct.id);

            // must after saving to cache since, since will file a event to load new data by new value in other pages.
            state.currSite = payload.currSite;
            state.currProduct = payload.currProduct;
        },
        saveProducts(state, payload) {
            console.log('payload', payload)
            state.products = payload
            state.modules = []
            state.suites = []
            state.tasks = []
        },
        saveModules(state, payload) {
            console.log('payload', payload)
            state.modules = payload
        },
        saveSuites(state, payload) {
            console.log('payload', payload)
            state.suites = payload
        },
        saveTasks(state, payload) {
            console.log('payload', payload)
            state.tasks = payload
        },
    },
    actions: {
        async fetchLangs({ commit }) {
            try {
                const response: ResponseData = await queryLang();
                const { data } = response;
                commit('saveLangs', data)

                return true;
            } catch (error) {
                return false;
            }
        },

        async getProfile({ commit }) {
            const response: ResponseData = await getProfile();
            const { data } = response;
            // data.avatar = ''
            commit('saveProfile', data)

            return true;
        },

        async fetchSitesAndProduct({ commit }, payload) {
            const response: ResponseData = await querySiteAndProduct(payload);
            const { data } = response;

            commit('saveSitesAndProduct', data)

            return true;
        },
        async fetchProducts({ commit }) {
            const response: ResponseData = await queryProduct();
            const { data } = response;
            commit('saveProducts', data)

            return true;
        },

        async fetchModules({ commit }, params) {
            try {
                const response: ResponseData = await queryModule(params);
                const { data } = response;
                commit('saveModules', data || 0);

                return true;
            } catch (error) {
                return false;
            }
        },
        async fetchSuites({ commit }, params) {
            try {
                const response: ResponseData = await querySuite(params);
                const { data } = response;
                commit('saveSuites', data || 0);

                return true;
            } catch (error) {
                return false;
            }
        },
        async fetchTasks({ commit }, params) {
            try {
                const response: ResponseData = await queryTask(params);
                const { data } = response;
                commit('saveTasks', data || 0);

                return true;
            } catch (error) {
                return false;
            }
        },
    }
}

export default StoreModel;
