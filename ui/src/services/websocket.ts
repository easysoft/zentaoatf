import * as neffos from 'neffos.js';
import { getCurrentInstance } from 'vue';
import {NSConn} from "neffos.js";
import {ComponentInternalInstance, ComponentPublicInstance} from "@vue/runtime-core";

const WebSocketPath = 'api/v1/ws';
export const WebSocketBaseDev = 'ws://127.0.0.1:8085/';
export type WsObject = {
  conn?: any;
};
export type WsEvent = {
  room: string;
  code: string;
  data: any;
};

export const WsEventName = 'ws_event'
export const WsDefaultNameSpace = 'default'
// export const WsDefaultRoom = 'default'

export class WebSocket {
  static conn: NSConn

  static async init(proxy: ComponentPublicInstance | any): Promise<any> {
    console.log(`init websocket`)
    if (!WebSocket.conn) {
      try {
        const conn = await neffos.dial(getWebSocketApi(), {
          default: {
            _OnNamespaceConnected: (nsConn, msg) => {
              if (nsConn.conn.wasReconnected()) {
                console.log('re-connected after ' + nsConn.conn.reconnectTries.toString() + ' trie(s)')
              }

              console.log('connected to namespace: ' + msg.Namespace)
              WebSocket.conn = nsConn
            },
            _OnNamespaceDisconnect: (_nsConn, msg) => {
              console.log('disconnected from namespace: ' + msg.Namespace)
            },
            OnVisit: (_nsConn, msg) => {
              console.log('OnVisit', msg)
            },
            // implement in webpage
            OnChat: (_nsConn, msg) => {
              console.log('OnChat in util cls', msg, msg.Room + ': response ' + msg.Body)
              proxy.$pub(WsEventName, {room: msg.Room, msg: msg.Body});
            }
          }
        })

        await conn.connect(WsDefaultNameSpace)

      } catch (err) {
        console.log(err)
      }
    }
    return WebSocket
  }

  static joinRoomAndSend(roomName: string, msg: string): void {
    if (!WebSocket.conn) return

    WebSocket.conn.joinRoom(roomName).then((room) => {
      console.log(`success to join room ${roomName}`)
      WebSocket.conn.room(roomName).emit('OnChat', msg)

    }).catch(err => {
      console.log(`fail to join room ${roomName}`, err)
    })
  }
  static sentMsg(roomName: string, msg: string): void {
    console.log(`send msg to room ${roomName}`)
    if (!WebSocket.conn) return

    WebSocket.conn.leaveAll().then(() =>
        this.joinRoomAndSend(roomName, msg)
    )
  }
}

export function getWebSocketApi (): string {
  const isProd = process.env.NODE_ENV === 'production'

  let wsUri = ''
  if (!isProd) {
    wsUri = WebSocketBaseDev
  } else {
    const loc = window.location

    if (loc.protocol === 'https:') {
      wsUri = 'wss:'
    } else {
      wsUri = 'ws:'
    }
    wsUri += '//' + loc.host
    wsUri += loc.pathname
  }

  return wsUri + WebSocketPath
}
