import { Mutation, Action } from 'vuex';
import { StoreModuleType } from '@/utils/store';
import { ResponseData } from '@/utils/request';
import { listProxy } from '@/views/proxy/service';
import {getCurrProxyId, setCurrProxyId} from "@/utils/cache";

export interface ProxyData {
  proxies: any[];
  proxyMap: Record<number, string>;
  currProxy: any;
}

export interface ModuleType extends StoreModuleType<ProxyData> {
  state: ProxyData;
  mutations: {
    saveProxies: Mutation<any>;
    saveProxyMap: Mutation<any>;
    saveCurrProxy: Mutation<any>;
  };
  actions: {
    fetchProxies: Action<ProxyData, ProxyData>;
    selectProxy: Action<ProxyData, ProxyData>;
  };
}

const initState: ProxyData = {
  proxies: [],
  proxyMap: {},
  currProxy:{},
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
    saveCurrProxy(state, payload) {
        if(payload.length == 0){
            payload = state.proxies;
        }
        console.log('payload', payload);
        getCurrProxyId().then((currProxyId) => {
            let currProxy = {};
            if(currProxyId == undefined || currProxyId == 0){
                currProxy = {id:0, name: '', path: 'local'};
            }else{
                payload.forEach(proxy => {
                    if(proxy.id == currProxyId){
                        currProxy = proxy
                    }
                })
            }
            state.currProxy = currProxy;
        })
      },
  },
  actions: {
    async fetchProxies({ commit }, params) {
      try {
        const response: ResponseData = await listProxy(params);
        const { data } = response;
        commit('saveProxies', data);
        commit('saveProxyMap', data);
        commit('saveCurrProxy', data);

        return true;
      } catch (error) {
        return false;
      }
    },
    async selectProxy({ commit }, payload) {
        await setCurrProxyId(payload.currProxyId);
        commit('saveCurrProxy', []);
      },
  },
};

export default StoreModel;
