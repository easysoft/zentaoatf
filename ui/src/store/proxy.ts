import { Mutation, Action } from 'vuex';
import { StoreModuleType } from '@/utils/store';
import { ResponseData } from '@/utils/request';
import { listProxy } from '@/views/proxy/service';

export interface ProxyData {
  proxies: any[];
  proxyMap: Record<number, string>;
}

export interface ModuleType extends StoreModuleType<ProxyData> {
  state: ProxyData;
  mutations: {
    saveProxies: Mutation<any>;
    saveProxyMap: Mutation<any>;
  };
  actions: {
    fetchProxies: Action<ProxyData, ProxyData>;
  };
}

const initState: ProxyData = {
  proxies: [],
  proxyMap: {},
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
    saveProxyMap(state, payload) {
        console.log('payload', payload);
        const map = {};
        payload.forEach(item => {
            map[item.id] = item.path;
        })
        state.proxyMap = map;
      },
  },
  actions: {
    async fetchProxies({ commit }, params) {
      try {
        const response: ResponseData = await listProxy(params);
        const { data } = response;
        commit('saveProxies', data);
        commit('saveProxyMap', data);

        return true;
      } catch (error) {
        return false;
      }
    },
  },
};

export default StoreModel;
