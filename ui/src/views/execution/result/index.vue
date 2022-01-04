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

      <div class="execution">
        <div class="title">测试步骤</div>
        <div class="desc">
          <div v-for="(step, index) in execution.steps" :key="index" class="step">
            <div class="cmd">{{step.action}} {{step.selector}}  {{step.value}}</div>
            <div class="capture" style="border: 3px cornflowerblue;"><img :src="step.image"></div>
          </div>
        </div>
      </div>

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
import {StateType as ListStateType} from "@/views/execution/store";
import {useRouter} from "vue-router";
import {ActionRecordStart, EventName, EventNodeId, ScopeDeeptest, ActionRecordedMsg} from "@/utils/const";
import {WebSocket, WsEventName} from "@/services/websocket";
import {getToken} from "@/utils/localToken";
import {ExecutionItem, StepItem} from "@/views/execution/data";

interface DesignExecutionPageSetupData {
  execution: ExecutionItem;

  loading: Ref<boolean>;
  getExecution:  (current: number) => Promise<void>;
  record: () => void;
  playback: () => void;
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
      const getExecution = async (id: number): Promise<void> => {
        loading.value = true;
        // await store.dispatch('ListExecution/getExecution', {
        // });
        loading.value = false;
      }

      const record = ():void =>  {
        console.log('record')

        window.postMessage({
          scope: ScopeDeeptest,
          content: {
            act: ActionRecordStart,
          }
        }, "*");
      }
      const playback = ():void =>  {
        console.log('playback')
      }

      const back = ():void =>  {
        router.push(`/~/execution/list`)
      }

      let init = true;
      let wsMsg = reactive({in: '', out: ''});

      let room: string | null = ''
      getToken().then((token) => {
        room = token
      })

      const sendWs = () => {
        console.log('sendWs');
        WebSocket.sentMsg(room, wsMsg.in);
        wsMsg.out = wsMsg.out + 'client: ' + wsMsg.in + '\n';
      };

      const { proxy } = getCurrentInstance() as any;
      WebSocket.init(proxy)

      onMounted(() => {
        const eventNode = document.getElementById(EventNodeId)
        if (eventNode) {
          eventNode.addEventListener(EventName, function () {
            const msg = JSON.parse(eventNode.innerText);
            console.log('====', msg);

            if (msg.scope !== ScopeDeeptest && msg.content.act !== ActionRecordedMsg) {
              return
            }

            const data :StepItem = msg.content.data
            execution.steps.push({
              action: data.action,
              selector: data.selector,
              value: data.value? data.value:'',
              image: data.image,
            })
          });
        }

        getExecution(1);
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
        getExecution,
        record,
        playback,
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
