import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import {queryWorkspace, deleteWorkspace} from "@/services/workspace";
import {setCache} from "@/utils/localCache";
import settings from '@/config/settings';
import {saveConfig} from "@/services/config";

export interface WorkspaceData {
  workspaces: any[]
  currWorkspace: any
  currConfig: any
  scriptTree: any[]
}

export interface ModuleType extends StoreModuleType<WorkspaceData> {
  state: WorkspaceData;
  mutations: {
    saveWorkspaces: Mutation<WorkspaceData>;
  };
  actions: {
    fetchWorkspace: Action<WorkspaceData, WorkspaceData>;
    removeWorkspace: Action<WorkspaceData, WorkspaceData>;
    saveConfig: Action<WorkspaceData, WorkspaceData>;
  };
}

const initState: WorkspaceData = {
  workspaces: [],
  currWorkspace: {},
  currConfig: {},
  scriptTree: [],
}

const StoreModel: ModuleType = {
  namespaced: true,
  name: 'workspace',
  state: {
    ...initState
  },
  mutations: {
    saveWorkspaces(state, payload) {
      console.log('payload', payload)

      setCache(settings.currWorkspace, payload.currWorkspace.path);

      state.workspaces = [...payload.workspaces];
      state.currWorkspace = payload.currWorkspace;
      state.currConfig = payload.currConfig;
      state.scriptTree = [payload.scriptTree];
    },
  },
  actions: {
    async fetchWorkspace({ commit }, currWorkspacePath) {
      try {
        const response: ResponseData = await queryWorkspace(currWorkspacePath);
        const { data } = response;
        commit('saveWorkspaces', data || {});

        return true;
      } catch (error) {
        return false;
      }
    },
    async removeWorkspace({ commit }, selectedWorkspacePath) {
      try {
        await deleteWorkspace(selectedWorkspacePath);

        const response: ResponseData = await queryWorkspace('');
        const { data } = response;
        commit('saveWorkspaces', data || {});

        return true;
      } catch (error) {
        return false;
      }
    },
    async saveConfig({ commit }, config) {
      const resp: ResponseData = await saveConfig(config);
      console.log('&&&')
      const { data } = resp;
      commit('saveWorkspaces', data);

      return resp;
    },
  }
}

export default StoreModel;
  