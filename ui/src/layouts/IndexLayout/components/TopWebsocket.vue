<template>
    <div></div>
</template>

<script lang="ts">
import {defineComponent, onMounted, onBeforeUnmount} from "vue";
import { useI18n } from "vue-i18n";

import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {notification} from "ant-design-vue";
import {ResultErr} from "@/utils/request";
import {WebSocket} from "@/services/websocket";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {WebSocketData} from "@/store/websoket";
import {WsMsg} from "@/types/data";

export default defineComponent({
    name: 'TopWebSocket',

    setup() {
      const { t } = useI18n();

      const websocketStore = useStore<{ WebSocket: WebSocketData }>();
      websocketStore.dispatch('WebSocket/connect')

      onMounted(() => {
        console.log('onMounted')
        bus.on(settings.eventWebSocketConnStatus, OnWebSocketEvent);
      })
      onBeforeUnmount( () => {
        bus.off(settings.eventWebSocketConnStatus, OnWebSocketEvent);
      })

      const OnWebSocketEvent = (data: any) => {
        console.log('OnWebSocketEvent in TopWebSocket', data.msg)

        const jsn = JSON.parse(data.msg) as WsMsg

        if (jsn.conn) { // update connection status
          websocketStore.dispatch('WebSocket/changeStatus', jsn.conn)
        }
      }

      return {
        t,
      }
    }
})
</script>