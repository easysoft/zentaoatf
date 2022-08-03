import {Mutation, Action} from 'vuex';
import { StoreModuleType } from "@/utils/store";
import {WebSocket} from "@/services/websocket";

export interface WebSocketData {
  connStatus: string
  room: string
}

export interface ModuleType extends StoreModuleType<WebSocketData> {
  state: WebSocketData;
  mutations: {
    saveConnStatus: Mutation<WebSocketData>;
    saveRoom: Mutation<WebSocketData>;
  };
  actions: {
    connect: Action<WebSocketData, WebSocketData>;
    changeStatus: Action<WebSocketData, WebSocketData>;
  };
}

const initState: WebSocketData = {
  connStatus: '',
  room: 'room',
}

const StoreModel: ModuleType = {
  namespaced: true,
  name: 'WebSocket',
  state: {
    ...initState
  },
  mutations: {
    saveConnStatus(state, payload) {
      console.log('saveConnection', payload)
      state.connStatus = payload
    },
    saveRoom (state, payload) {
      console.log('saveRoom', payload)
      state.room = payload
    }
  },
  actions: {
    async connect(_ctx, room) {
      console.log("connect to websocket")
      await WebSocket.init(false)

      const msg = {act: 'init'}
      WebSocket.sentMsg(room, JSON.stringify(msg))

      return true;
    },
    async changeStatus({ commit }, status) {
      console.log("changeStatus")
      commit('saveConnStatus', status)

      return true;
    },
  }
}

export default StoreModel;
