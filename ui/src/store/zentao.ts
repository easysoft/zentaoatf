import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import {queryLang, querySiteAndProduct, getProfile, queryProduct, queryModule, querySuite, queryTask} from "../services/zentao";
import {setCache} from "@/utils/localCache";
import settings from "@/config/settings";

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

        saveSitesAndProducts: Mutation<any>;
        saveProducts: Mutation<any>;
        saveModules: Mutation<any>;
        saveSuites: Mutation<any>;
        saveTasks: Mutation<any>;
    };
    actions: {
        fetchLangs: Action<ZentaoData, ZentaoData>;
        getProfile: Action<ZentaoData, ZentaoData>;

        fetchSitesAndProducts: Action<ZentaoData, ZentaoData>;
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
    name: 'zentao',
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
        saveSitesAndProducts(state, payload) {
            state.sites = payload.sites;
            state.products = payload.products;

            setCache(settings.currSiteId, payload.currSite.id);
            setCache(settings.currProductId, payload.currProduct.id);
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
            try {
                const response: ResponseData = await getProfile();
                const { data } = response;
                // data.avatar = ''
                commit('saveProfile', data)

                return true;
            } catch (error) {
                return false;
            }
        },
        async fetchSitesAndProducts({ commit }) {
            const response: ResponseData = await querySiteAndProduct();
            const { data } = response;
            commit('saveSitesAndProducts', data)

            return true;
        },
        async fetchProducts({ commit }) {
            const response: ResponseData = await queryProduct();
            const { data } = response;
            commit('saveProducts', data)

            return true;
        },
        async fetchModules({ commit }, productId) {
            try {
                const response: ResponseData = await queryModule(productId);
                const { data } = response;
                commit('saveModules', data || 0);

                return true;
            } catch (error) {
                return false;
            }
        },
        async fetchSuites({ commit }, productId) {
            try {
                const response: ResponseData = await querySuite(productId);
                const { data } = response;
                commit('saveSuites', data || 0);

                return true;
            } catch (error) {
                return false;
            }
        },
        async fetchTasks({ commit }, productId) {
            try {
                const response: ResponseData = await queryTask(productId);
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
