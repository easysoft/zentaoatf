import { Mutation, Action } from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { ResponseData } from '@/utils/request';
import {queryProject} from "@/services/project";
import {ComputedRef} from "vue";

export interface ProjectData {
  projects: any[]
  currProject: any
  scriptTree: any[]
  scriptTreeOpenKeys: string[]
}

export interface ModuleType extends StoreModuleType<ProjectData> {
  state: ProjectData;
  mutations: {
    saveProjects: Mutation<ProjectData>;
    saveOpenKeys: Mutation<ProjectData>;
  };
  actions: {
    fetchProject: Action<ProjectData, ProjectData>;
    genOpenKeys: Action<ProjectData, ProjectData>
  };
}

const initState: ProjectData = {
  projects: [],
  currProject: {},
  scriptTree: [],
  scriptTreeOpenKeys: [],
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
      state.scriptTree = [payload.scriptTree];

      state.scriptTreeOpenKeys = []
      state.scriptTreeOpenKeys.push(payload.scriptTree.path)
    },

    saveOpenKeys(state, isExpand) {
      state.scriptTreeOpenKeys = []
      console.log('saveOpenKeys', isExpand)
      if (isExpand) {
        getOpenKeys(state.scriptTree[0], state.scriptTreeOpenKeys)
      }
    }
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
    async genOpenKeys({ commit }, isExpand) {
      try {
        console.log('genOpenKeys', isExpand)
        commit('saveOpenKeys', isExpand);

        return true;
      } catch (error) {
        return false;
      }
    }
  }
}

const getOpenKeys = (node, keys) => {
  if (!node) return

  keys.push(node.path)
  if (node.children) {
    node.children.forEach((item, index) => {
      getOpenKeys(item, keys)
    })
  }
}

export default StoreModel;
  