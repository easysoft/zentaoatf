import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";

export interface ExecStatus {
  isRunning: string
}

export interface ModuleType extends StoreModuleType<ExecStatus> {
  state: ExecStatus;
  mutations: {
    updateRunning: Mutation<ExecStatus>;
  };
  actions: {
    setRunning: Action<ExecStatus, ExecStatus>;
  };
}

const initState: ExecStatus = {
  isRunning: 'false',
}

const StoreModel: ModuleType = {
  namespaced: true,
  name: 'Exec',
  state: {
    ...initState
  },
  mutations: {
    updateRunning(state, payload) {
      console.log('payload', payload)
      state.isRunning = payload;
    },
  },
  actions: {
    async setRunning({ commit }, running) {
      commit('updateRunning', running);
      return true;
    },
  }
}

export default StoreModel;
