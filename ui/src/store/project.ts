import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import {queryProject} from "@/services/project";

export interface ProjectData {
  projects: any[]
  currProject: any
  currConfig: any
  scriptTree: any[]
}

export interface ModuleType extends StoreModuleType<ProjectData> {
  state: ProjectData;
  mutations: {
    saveProjects: Mutation<ProjectData>;
  };
  actions: {
    fetchProject: Action<ProjectData, ProjectData>;
  };
}

const initState: ProjectData = {
  projects: [],
  currProject: {},
  currConfig: {},
  scriptTree: [],
}

const StoreModel: ModuleType = {
  namespaced: true,
  name: 'project',
  state: {
    ...initState
  },
  mutations: {
    saveProjects(state, payload) {
      console.log('payload', payload)

      state.projects = [...payload.projects, {id: 0, name: '其他', path: ''}];
      state.currProject = payload.currProject;
      state.currConfig = payload.currConfig;
      state.scriptTree = [payload.scriptTree];
    },
  },
  actions: {
    async fetchProject({ commit }, currProjectPath) {
      try {
        const response: ResponseData = await queryProject(currProjectPath);
        const { data } = response;
        commit('saveProjects', data || 0);

        return true;
      } catch (error) {
        return false;
      }
    },
  }
}

export default StoreModel;
  