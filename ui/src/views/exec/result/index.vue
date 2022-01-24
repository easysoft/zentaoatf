<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
       执行结果详情
      </template>
      <template #extra>
        <a-button type="link" @click="() => back()">返回</a-button>
      </template>

      <div class="main">
        细节
      </div>

    </a-card>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, onBeforeUnmount, getCurrentInstance, ComputedRef, Ref, ref, reactive} from "vue";
import { useStore } from 'vuex';
import {StateType as ListStateType} from "@/views/exec/store";
import {useRouter} from "vue-router";
import {WebSocket, WsEventName} from "@/services/websocket";
import {getToken} from "@/utils/localToken";
import {ExecutionItem} from "@/views/exec/data";

interface DesignExecutionPageSetupData {
  execution: ExecutionItem;

  loading: Ref<boolean>;
  back: () => void;

  wsMsg: any,
  sendWs: () => void;
}

export default defineComponent({
    name: 'ExecutionResultPage',
    setup(): DesignExecutionPageSetupData {
      const router = useRouter();
      const store = useStore<{ ListExecution: ListStateType}>();

      const execution = reactive<ExecutionItem>({steps: []})
      const loading = ref<boolean>(true);

      const seq = +router.currentRoute.value.params.seq
      console.log('seq', seq)

      const back = ():void =>  {
        router.push(`/exec/history`)
      }

      let init = true;
      let wsMsg = reactive({in: '', out: ''});

      let room = ''
      getToken().then((token) => {
        room = token || ''
      })

      const sendWs = () => {
        console.log('sendWs');
        WebSocket.sentMsg(room, wsMsg.in);
        wsMsg.out = wsMsg.out + 'client: ' + wsMsg.in + '\n';
      };

      const { proxy } = getCurrentInstance() as any;
      WebSocket.init(proxy)

      onMounted(() => {
        if (init) {
          proxy.$sub(WsEventName, (data) => {
            console.log(data[0].msg);
            wsMsg.out = wsMsg.out + 'server: ' + data[0].msg + '\n';
            console.log(wsMsg.out);
          });
          init = false;
        }
      });
      onBeforeUnmount(() => {
        proxy.$unsub(WsEventName, () => {
          console.log('unsub event ' + WsEventName);
        });
      });

      return {
        execution,
        loading,
        back,
        wsMsg,
        sendWs,
      }
    }
})
</script>

<style lang="less" scoped>
  .main {
    padding: 20px;
  }
</style>
