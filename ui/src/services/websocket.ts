import * as neffos from 'neffos.js';
import {NSConn} from "neffos.js";

import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import { getCache } from '@/utils/localCache';

export type WsEvent = {
  room: string;
  code: string;
  data: any;
};

export const WsDefaultNameSpace = 'default'

export class WebSocket {
  static conns: Record<string, NSConn>

  static async init(reConn, appApiHost): Promise<any> {
    console.log(`init websocket`, WebSocket.conns, appApiHost)
    if(WebSocket.conns == undefined) WebSocket.conns = {};
    if (reConn || WebSocket.conns[appApiHost] == undefined) {
      try {
        const conn = await neffos.dial(await getWebSocketApi(appApiHost), {
          default: {
            _OnNamespaceConnected: (nsConn, msg) => {
              if (nsConn.conn.wasReconnected()) {
                console.log('re-connected after ' + nsConn.conn.reconnectTries.toString() + ' trie(s)')
              }

              console.log('connected to namespace: ' + msg.Namespace)
              WebSocket.conns[appApiHost] = nsConn
              bus.emit(settings.eventWebSocketConnStatus, {msg: '{"conn": "success"}'});
            },

            _OnNamespaceDisconnect: (_nsConn, msg) => {
              delete WebSocket.conns[appApiHost]
              console.log('disconnected from namespace: ' + msg.Namespace)
            },

            OnVisit: (_nsConn, msg) => {
              console.log('OnVisit', msg)
            },

            // implement in webpage
            OnChat: (_nsConn, msg) => {
              console.log('OnChat in util cls', msg, msg.Room + ': response ' + msg.Body)
              bus.emit(settings.eventWebSocketMsg, {room: msg.Room, msg: msg.Body});
            }
          }
        })

        await conn.connect(WsDefaultNameSpace)

      } catch (err) {
        console.log('failed connect to websocket', err)
        bus.emit(settings.eventWebSocketConnStatus, {msg: '{"conn": "fail"}'});
      }
    }
    return WebSocket
  }

  static joinRoomAndSend(roomName: string, msg: string, appApiHost:string): void {
    // if (WebSocket.conns[appApiHost] && WebSocket.conns[appApiHost].room(roomName)) {
    //   WebSocket.conns[appApiHost].room(roomName).emit('OnChat', msg)
    // } else {
      WebSocket.init(true, appApiHost).then(() => {
        WebSocket.conns[appApiHost].joinRoom(roomName).then((_room) => {
          console.log(`success to join room "${roomName}"`)
          _room.emit('OnChat', msg)
        }).catch(err => {
          console.log(`fail to join room ${roomName}`, err)
          bus.emit(settings.eventWebSocketConnStatus, {msg: '{"conn": "fail"}'});
          WebSocket.conns[appApiHost].disconnect()
          this.joinRoomAndSend(roomName, msg, appApiHost)
        })
      })
    // }
  }

  static sentMsg(roomName: string, msg: string, appApiHost: string): void {
    console.log(`send msg to room "${roomName}"`)
    if (WebSocket.conns[appApiHost]){
        WebSocket.conns[appApiHost].disconnect().then(() =>
        this.joinRoomAndSend(roomName, msg, appApiHost)
        )
    }else{
        this.joinRoomAndSend(roomName, msg, appApiHost)
    }
  }
}

export async function getWebSocketApi (appApiHost): Promise<string> {
  const isProd = process.env.NODE_ENV === 'production'
  const loc = window.location
  console.log(`${isProd}, ${loc.toString()}`)
  const serverUrl = await getServerUrl();
  const apiHost = appApiHost && appApiHost != 'local' ? appApiHost + process.env.VUE_APP_APISUFFIX : serverUrl

  const url = apiHost.replace('http', 'ws') + '/ws'
  console.log(`websocket url = ${url}`, appApiHost)

  return url
}

export async function getServerUrl(): Promise<string>{
  let serverURL = await getCache(settings.currServerURL);
  if (!serverURL || serverURL == 'local') {
    serverURL = process.env.VUE_APP_APIHOST;
  } else {
    serverURL = String(serverURL) + 'api/v1';
  }
  return serverURL
}