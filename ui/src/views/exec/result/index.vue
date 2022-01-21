<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
        <div>
          <a-button type="primary" @click="() => record()" id="com-deeptest-record" class="act-btn">录制</a-button>
          <a-button @click="() => playback()" class="act-btn">播放</a-button>
        </div>
      </template>
      <template #extra>
        <a-button type="link" @click="() => back()">返回</a-button>
      </template>

      <br /> <!--WebSocket Test-->
      <div>
        <div><a-input id="input" type="text" v-model:value="wsMsg.in" /></div>
        <div><a-button id="sendBtn" @click="sendWs">Send</a-button></div>
        <div>
          <pre>{{ wsMsg.out }}</pre>
        </div>
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

      const id = +router.currentRoute.value.params.id
      console.log('id', id)

      const back = ():void =>  {
        router.push(`/execution/list`)
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
  .act-btn {
    margin-right: 20px;
  }

  .execution {
    .title {
      font-weight: bolder;
    }
    .desc {
      .step {
        display: flex;
        .cmd {
          flex: 1;
        }
        .capture {
          width: 600px;

          img {
            height: 50px;
          }
        }
      }
    }
  }
</style>
