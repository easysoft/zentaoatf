
import { Mutation/* , Action */ , Getter} from 'vuex';
import { StoreModuleType } from "@/utils/store";
import { TabNavItem } from '@/utils/routes';
import settings from '@/config/settings';
import router from '@/config/routes';

export interface StateType {
  // 左侧展开收起
  collapsed: boolean;
  // 顶部菜单开启
  topNavEnable: boolean;
  // 头部固定开启
  headFixed: boolean;

  logPaneOriginSize: number;
  
  logPaneMaximized: boolean;
}

export interface ModuleType extends StoreModuleType<StateType> {
  state: StateType;
  getters: {
    logPaneSize: Getter<StateType, StateType>;
    editorPaneSize: Getter<StateType, StateType>;
  };
  mutations: {
    changeLayoutCollapsed: Mutation<StateType>;
    setTopNavEnable: Mutation<StateType>;
    setHeadFixed: Mutation<StateType>;
    setLogPaneResized:Mutation<StateType>;
  };
  actions: {
  };
}

const initState: StateType = {
  collapsed: false,
  topNavEnable: settings.topNavEnable,
  headFixed: settings.headFixed,
  logPaneOriginSize: 70, // settings.logPaneSize,
  logPaneMaximized: false
};

const StoreModel: ModuleType = {
  namespaced: true,
  name: 'global',
  state: {
    ...initState
  },
  getters: {
    logPaneSize(state) {
      if (state.logPaneMaximized) {
        return 100;
      }

      return state.logPaneOriginSize;
    },
    editorPaneSize(state) {
      if (state.logPaneMaximized) {
        return 0;
      }

      return 100 - state.logPaneOriginSize;
    }
  },
  mutations: {
    changeLayoutCollapsed(state, payload) {
      state.collapsed = payload;
    },
    setTopNavEnable(state, payload) {
      state.topNavEnable = payload;
    },
    setHeadFixed(state, payload) {
      state.headFixed = payload;
    },
    setLogPaneResized(state) {
      state.logPaneMaximized = !state.logPaneMaximized
    }
  },
  actions: {}
}



export default StoreModel;
