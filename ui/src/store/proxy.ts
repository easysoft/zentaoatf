import { Mutation, Action } from 'vuex';
import { StoreModuleType } from '@/utils/store';
import { ResponseData } from '@/utils/request';
import { listProxy } from '@/views/proxy/service';

export interface ProxyData {
  proxies: any[];
}

export interface ModuleType extends StoreModuleType<ProxyData> {
  state: ProxyData;
  mutations: {
    saveProxies: Mutation<any>;
  };
  actions: {
    fetchProxies: Action<ProxyData, ProxyData>;
  };
}

const initState: ProxyData = {
  proxies: [],
};

const StoreModel: ModuleType = {
  namespaced: true,
  name: 'proxy',
  state: {
    ...initState,
  },
  mutations: {
    saveProxies(state, payload) {
      console.log('payload', payload);
      state.proxies = payload;
    },
  },
  actions: {
    async fetchProxies({ commit }, params) {
      try {
        const response: ResponseData = await listProxy(params);
        const { data } = response;
        commit('saveProxies', data || 0);

        return true;
      } catch (error) {
        return false;
      }
    },
  },
};

export default StoreModel;
