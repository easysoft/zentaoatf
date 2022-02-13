<template>
    <div class="indexlayout-main-conent">
        <a-card :bordered="false">
          <template #title>
            {{t('exec')}}{{t('unit_test')}}
          </template>
          <template #extra>
            <div class="opt">
              <a-button v-if="isRunning == 'false'" @click="exec" type="primary">{{t('exec')}}</a-button>
              <a-button v-if="isRunning == 'true'" @click="stop" type="primary">{{t('stop')}}</a-button>

              <a-button @click="back" type="link">{{t('back')}}</a-button>
            </div>
          </template>

          <div id="main">
            <div id="left">
              <a-form :label-col="labelCol" :wrapper-col="wrapperCol">

                <a-form-item :label="t('product')" v-bind="validateInfos.productId">
                  <a-select v-model:value="model.productId">
                    <a-select-option key="" value="">&nbsp;</a-select-option>
                    <a-select-option v-for="item in products" :key="item.id" :value="item.id+''">{{item.name}}</a-select-option>
                  </a-select>
                </a-form-item>

                <a-form-item :label="t('framework_tool')" v-bind="validateInfos.testTool">
                  <a-select v-model:value="model.testTool">
                    <a-select-option key="" value="">&nbsp;</a-select-option>
                    <a-select-option v-for="item in unitTestFrameworks.list" :key="item" :value="item">
                      {{unitTestFrameworks.map[item]}}
                    </a-select-option>
                  </a-select>
                </a-form-item>

                <a-form-item :label="t('build_tool')" v-if="model.testTool=='junit' || model.testTool=='testng'"
                             v-bind="validateInfos.buildTool">
                  <a-select v-model:value="model.buildTool">
                    <a-select-option key="" value="">&nbsp;</a-select-option>
                    <a-select-option v-for="item in unitTestTools.data[model.testTool]" :key="item" :value="item">
                      {{unitTestTools.map[item]}}
                    </a-select-option>
                  </a-select>
                </a-form-item>

                <a-form-item :label="t('test_cmd')" v-bind="validateInfos.cmd">
                  <a-textarea v-model:value="model.cmd" placeholder="mvn clean package test"
                              :auto-size="{ minRows: 3, maxRows: 6 }" />
                </a-form-item>

                <a-form-item :label="t('submit_to_zentao')" v-bind="validateInfos.submitResult">
                  <a-switch v-model:checked="model.submitResult" />
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

import {useStore} from "vuex";
import {ProjectData} from "@/store/project";
import {ZentaoData} from "@/store/zentao";

import {useRouter} from "vue-router";
import {getCache} from "@/utils/localCache";
import settings from "@/config/settings";
import {WebSocket, WsEventName} from "@/services/websocket";
import {resizeWidth, scroll} from "@/utils/dom";
import {genExecInfo} from "@/views/exec/service";
import {getUnitTestFrameworks, getUnitTestTools} from "@/utils/testing";
import throttle from "lodash.debounce";
import {useI18n} from "vue-i18n";

interface ExecCasePageSetupData {
  t: (key: string | number) => string;
  labelCol: any
  wrapperCol: any

  model: Ref;
  products: ComputedRef<any[]>;
  unitTestFrameworks: Ref
  unitTestTools: Ref

  wsMsg: any,
  exec: (keys) => void;
  stop: (keys) => void;
  back: () => void;
  isRunning: Ref<string>;

  rules: any
  validate: any
  validateInfos: validateInfos,
  resetFields:  () => void;
}

export default defineComponent({
    name: 'ExecutionSuitePage',
    components: {
    },
    setup(): ExecCasePageSetupData {
      const { t } = useI18n();

      const router = useRouter();

      const unitTestFrameworks = getUnitTestFrameworks()
      const unitTestTools = getUnitTestTools()

      const storeProject = useStore<{ project: ProjectData }>();
      const currConfig = computed<any>(() => storeProject.state.project.currConfig);

      const store = useStore<{zentao: ZentaoData}>();
      const products = computed<any[]>(() => store.state.zentao.products);

      const fetchProducts = throttle((): void => {
        store.dispatch('zentao/fetchProducts').catch((error) => {
          if (error.response.data.code === 2000) router.push(`/config`)
        })
      }, 600)
      fetchProducts()
      watch(currConfig, ()=> {
        fetchProducts()
      })

      let init = true;
      let isRunning = ref('false');
      let wsMsg = reactive({in: '', out: ''});

      let room = ''
      getCache(settings.currProject).then((token) => {
        room = token || ''
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
          i++
          scroll('logs')
        });
        init = false;
      }

      onMounted(() => {
        console.log('onMounted')
        resizeWidth('main', 'left', 'resize', 'content', 280, 800)
        initWsConn()
      })

      const model = reactive<any>({productId: '', testTool: '', buildTool: '', cmd: ''});

      const rules = reactive({
        productId: [
          { required: true, message: t('pls_product') },
        ],
        testTool: [
          { required: true, message: t('pls_test_tool') },
        ],
        buildTool: [
          { required: true, message: t('pls_build_tool') },
        ],
        cmd: [
          { required: true, message: t('pls_cmd') },
        ],
      });

      const { resetFields, validate, validateInfos } = useForm(model, rules);

      const exec = (): void => {
        console.log("exec")

        validate().then(() => {
          getCache(settings.currProject).then(
              (projectPath) => {
                const msg = Object.assign({act: 'execUnit', projectPath: projectPath}, model)
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
        t,
        labelCol: { span: 6 },
        wrapperCol: { span: 16 },

        products,
        model,
        unitTestFrameworks,
        unitTestTools,
        wsMsg,

        rules,
        validate,
        validateInfos,
        resetFields,

        isRunning,
        exec,
        stop,
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
    padding: 16px;
    overflow: auto;

    #logs {
      margin: 0;
      padding: 0;
      height: calc(100% - 10px);
      width: 100%;
      overflow-y: auto;
      white-space: pre-wrap;
      word-wrap: break-word;
      font-family:monospace;
    }
  }
}
</style>