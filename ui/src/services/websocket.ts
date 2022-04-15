import * as neffos from 'neffos.js';
import {NSConn} from "neffos.js";
import {ComponentPublicInstance} from "@vue/runtime-core";

const WebSocketPath = 'api/v1/ws';

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
              proxy.$pub(WsEventName, {msg: '{"conn": "success"}'});
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
        console.log('failed connect to websocket', err)
        proxy.$pub(WsEventName, {msg: '{"conn": "fail"}'});
      }
    }
    return WebSocket
  }

  static joinRoomAndSend(roomName: string, msg: string): void {
    if (!WebSocket.conn) return

    WebSocket.conn.joinRoom(roomName).then((room) => {
      console.log(`success to join room "${roomName}"`)
      WebSocket.conn.room(roomName).emit('OnChat', msg)

    }).catch(err => {
      console.log(`fail to join room ${roomName}`, err)
    })
  }

  static sentMsg(roomName: string, msg: string): void {
    console.log(`send msg to room "${roomName}"`)
    if (!WebSocket.conn) return

    WebSocket.conn.leaveAll().then(() =>
        this.joinRoomAndSend(roomName, msg)
    )
  }
}

export function getWebSocketApi (): string {
  const isProd = process.env.NODE_ENV === 'production'

  const loc = window.location
  console.log(`${isProd}, ${loc.toString()}`)

  const wsUri = process.env.VUE_APP_APIHOST.replace('http', 'ws')
  const url = wsUri + WebSocketPath
  console.log(`websocket url = ${url}`)

  return url
}
