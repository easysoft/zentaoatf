import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import {queryProject} from "@/services/project";
import {ComputedRef} from "vue";

export interface ProjectData {
  projects: any[];
  currProject: any;
  scriptTree: any[];
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

      state.projects = payload.projects;
      state.currProject = payload.currProject;
      state.scriptTree = [payload.scriptTree];
    }
  },
  actions: {
    async fetchProject({ commit }) {
      try {
        const response: ResponseData = await queryProject();
        const { data } = response;
        commit('saveProjects', data || 0);

        return true;
      } catch (error) {
        return false;
      }
    }
  }
}

export default StoreModel;
  