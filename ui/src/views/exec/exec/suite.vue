<template>
    <div class="indexlayout-main-conent">
        <a-card :bordered="false">
          <template #title>
            执行套件
          </template>
          <template #extra>
            <div class="opt">
              <a-button v-if="isRunning == 'false'" @click="exec" type="primary">执行</a-button>
              <a-button v-if="isRunning == 'true'" @click="stop" type="primary">停止</a-button>

              <a-button @click="back" type="link">返回</a-button>
            </div>
          </template>

          <div id="main">
            <div id="left">
            </div>

            <div id="resize"></div>

            <div id="content">
            </div>
          </div>

        </a-card>
    </div>
</template>

<script lang="ts">
import {ComputedRef, defineComponent, ref, Ref, reactive, computed, onMounted, getCurrentInstance} from "vue";
import { SelectTypes } from 'ant-design-vue/es/select';
import {Execution} from '../data.d';
import {useStore} from "vuex";

import { Props } from 'ant-design-vue/lib/form/useForm';
import { message, Modal, Form } from "ant-design-vue";
const useForm = Form.useForm;

import CreateForm from './components/CreateForm.vue';
import UpdateForm from './components/UpdateForm.vue';

import {StateType as ListStateType} from "../store";
import debounce from "lodash.debounce";
import {useRoute, useRouter} from "vue-router";
import {getCache} from "@/utils/localCache";
import settings from "@/config/settings";
import {WebSocket, WsEventName} from "@/services/websocket";
import {resizeWidth, scroll, SetWidth} from "@/utils/dom";

interface ExecCasePageSetupData {
  model: any

  wsMsg: any,
  exec: (keys) => void;
  stop: (keys) => void;
  isRunning: Ref<string>;
  back: () => void;
}

export default defineComponent({
    name: 'ExecutionSuitePage',
    components: {
    },
    setup(): ExecCasePageSetupData {
      const router = useRouter();
      const model = {}

      let init = true;
      let isRunning = ref('');
      let wsMsg = reactive({in: '', out: ''});

      let room: string | null = ''
      getCache(settings.currProject).then((token) => {
        room = token
      })

      const {proxy} = getCurrentInstance() as any;
      WebSocket.init(proxy)

      let i = 1
      if (init) {
        proxy.$sub(WsEventName, (data) => {
          console.log(data[0].msg);
          const jsn = JSON.parse(data[0].msg)

          if ('isRunning' in jsn) {
            isRunning.value = jsn.isRunning
          }

          let msg = jsn.msg
          msg = msg.replace(/^"+/,'').replace(/"+$/,'')
          msg = SetWidth(i++ + '. ', 40) + `<span>${msg}</span>`;

          let sty = ''
          if (jsn.category === 'exec') {
            sty = 'color: #009688;'
          } else if (jsn.category === 'output') {
            // sty = 'font-style: italic;'
          }

          msg = `<div style="${sty}"> ${msg} </div>`
          wsMsg.out += msg

          scroll('logs')
        });
        init = false;
      }

      onMounted(() => {
        console.log('onMounted')
        resizeWidth('main', 'left', 'resize', 'content', 280, 800)
        initWsConn()
      })

      const exec = (): void => {
        console.log("exec")

        getCache(settings.currProject).then (
            (projectPath) => {
              const msg = {act: 'execSuite', projectPath: projectPath}
              console.log('msg', msg)

              wsMsg.out += '\n'
              WebSocket.sentMsg(room, JSON.stringify(msg))
            }
        )
      }
      const stop = (): void => {
        console.log("stop")
        getCache(settings.currProject).then (
            (projectPath) => {
              const msg = {act: 'execStop', projectPath: projectPath}
              console.log('msg', msg)
              WebSocket.sentMsg(room, JSON.stringify(msg))
            }
        )
      }
      const initWsConn = (): void => {
        console.log("initWsConn")
        getCache(settings.currProject).then (
            (projectPath) => {
              const msg = {act: 'init', projectPath: projectPath}
              console.log('msg', msg)
              WebSocket.sentMsg(room, JSON.stringify(msg))
            }
        )
      }

      const back = (): void => {
        router.push(`/exec/history`)
      }

      return {
        model,
        wsMsg,
        exec,
        stop,
        isRunning,
        back,
      }
    }

})
</script>

<style lang="less" scoped>
.indexlayout-main-conent {
  height: calc(100% - 50px);
}

#main {
  display: flex;
  height: 100%;

  #left {
    width: 380px;
    height: 100%;
    padding: 3px;
  }

  #resize {
    width: 2px;
    height: 100%;
    background-color: #D0D7DE;
    cursor: ew-resize;

    &.active {
      background-color: #a9aeb4;
    }
  }

  #content {
    flex: 1;
    height: 100%;
  }
}
</style>