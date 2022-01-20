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
              <a-form :label-col="labelCol" :wrapper-col="wrapperCol">

                <a-form-item label="产品" v-bind="validateInfos.productId">
                  <a-select v-model:value="model.productId" @change="selectProduct">
                    <a-select-option key="" value="">&nbsp;</a-select-option>
                    <a-select-option v-for="item in products" :key="item.id" :value="item.id">{{item.name}}</a-select-option>
                  </a-select>
                </a-form-item>

                <a-form-item label="模块" v-bind="validateInfos.moduleId">
                  <a-select v-model:value="model.moduleId">
                    <a-select-option key="" value="">&nbsp;</a-select-option>
                    <a-select-option v-for="item in modules" :key="item.id" :value="item.id"><span v-html="item.name"></span></a-select-option>
                  </a-select>
                </a-form-item>

              </a-form>
            </div>

            <div id="resize"></div>

            <div id="content">
              <div id="logs">
                <span v-html="wsMsg.out"></span>
              </div>
            </div>
          </div>

        </a-card>
    </div>
</template>

<script lang="ts">
import {ComputedRef, defineComponent, ref, Ref, reactive, computed, onMounted, getCurrentInstance, watch} from "vue";
import { validateInfos } from 'ant-design-vue/lib/form/useForm';
import {message, Form} from 'ant-design-vue';
const useForm = Form.useForm;

import { ExecutionBy } from '../data.d';
import {useStore} from "vuex";
import {ProjectData} from "@/store/project";
import {ZentaoData} from "@/store/zentao";

import {useRouter} from "vue-router";
import {getCache} from "@/utils/localCache";
import settings from "@/config/settings";
import {WebSocket, WsEventName} from "@/services/websocket";
import {resizeWidth, scroll, SetWidth} from "@/utils/dom";
import {genExecInfo} from "@/views/exec/service";

interface ExecCasePageSetupData {
  model: Ref;

  wsMsg: any,
  exec: (keys) => void;
  stop: (keys) => void;
  isRunning: Ref<string>;
  back: () => void;

  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos,
  products: ComputedRef<any[]>;
  modules: ComputedRef<any[]>;
  selectProduct:  (item) => void;
}

export default defineComponent({
    name: 'ExecutionSuitePage',
    components: {
    },
    setup(): ExecCasePageSetupData {
      const storeProject = useStore<{ project: ProjectData }>();
      const currConfig = computed<any>(() => storeProject.state.project.currConfig);

      const store = useStore<{zentao: ZentaoData}>();
      const products = computed<any[]>(() => store.state.zentao.products);
      const modules = computed<any[]>(() => store.state.zentao.modules);

      store.dispatch('zentao/fetchProducts')
      watch(currConfig, (currConfig)=> {
        store.dispatch('zentao/fetchProducts')
      })

      const formRef = ref();

      const model = reactive<ExecutionBy>({
        productId: '',
        moduleId: '',
      } as ExecutionBy);

      const rules = reactive({
        productId: [
          { required: true, message: '请选择产品' },
        ],
        moduleId: [
          { required: true, message: '请选择模块' },
        ],
      });

      const { resetFields, validate, validateInfos } = useForm(model, rules);

      const selectProduct = (item) => {
        console.log('selectProduct', item)
        if (!item) return

        store.dispatch('zentao/fetchModules', item)
      };

      const router = useRouter();

      let init = true;
      let isRunning = ref('false');
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

          wsMsg.out += genExecInfo(jsn, i)
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
        validate().then(() => {
          getCache(settings.currProject).then(
              (projectPath) => {
                const msg = Object.assign({act: 'execModule', projectPath: projectPath}, model)
                console.log('msg', msg)

                wsMsg.out += '\n'
                WebSocket.sentMsg(room, JSON.stringify(msg))
              }
          )
        }).catch(err => {console.log('validate fail', err)});
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

        formRef,
        labelCol: { span: 6 },
        wrapperCol: { span: 16 },
        rules,
        validate,
        validateInfos,
        resetFields,

        products,
        modules,
        selectProduct,

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
    padding: 20px 3px;
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