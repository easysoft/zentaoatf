<template>
  <div class="script-main">
    <div id="main">
      <div id="left">
        <ScriptTreePage></ScriptTreePage>
      </div>

      <div id="splitter-h"></div>

      <div id="right">
        <!-- ZTF Test -->
        <template v-if="currWorkspace?.type === 'ztf'">
          <div class="toolbar" v-if="scriptCode !== ''">
            <a-button :disabled="isRunning === 'true'" @click="execSingle" type="primary" class="t-btn-gap">
              {{ t('exec') }}
            </a-button>
            <a-button v-if="isRunning === 'true'" @click="execStop" class="t-btn-gap">{{ t('stop') }}</a-button>

            <a-button @click="extract" type="primary">{{ t('extract_step') }}</a-button>

            <a-button @click="save" type="primary">{{ t('save') }}</a-button>
          </div>

          <div id="right-content" class="right-content">
            <!-- Exec Single Script -->
            <template v-if="script">
              <div id="editor-panel" class="editor-panel">
                <MonacoEditor
                    v-if="scriptCode !== ''"
                    v-model:value="scriptCode"
                    :language="lang"
                    :options="editorOptions"
                    class="editor"
                    ref="editorRef"
                />
              </div>

              <div id="splitter-v" class="splitter-v"></div>

              <div id="logs-panel" class="logs-panel">
                <ScriptExecLogPage></ScriptExecLogPage>
              </div>
            </template>

            <!-- Exec Selected Script -->
            <template v-if="!script">
              <div class="logs-panel">
                <ScriptExecLogPage></ScriptExecLogPage>
              </div>
            </template>
          </div>
        </template>

        <!-- Unit Test -->
        <template v-if="currWorkspace?.type !== 'ztf'">
          <div class="toolbar" v-if="scriptCode !== ''">
            <a-button @click="save" type="primary">{{ t('save') }}</a-button>
          </div>

          <div class="right-content">
            <template v-if="script">
              <MonacoEditor
                  v-if="scriptCode !== ''"
                  v-model:value="scriptCode"
                  :language="lang"
                  :options="editorOptions"
                  class="editor"
                  ref="editorRef"
              />
            </template>

            <template v-if="!script">
              <div class="unit-panel">
                <a-form :model="modelUnit" layout="inline">
                  <a-form-item :label="t('test_cmd')">
                    <a-input
                        v-model:value="modelUnit.cmd"
                        @keydown.down="down"
                        @keydown.up="up"
                        style="width:380px;"/>
                  </a-form-item>

                  <a-form-item>
                    <a-checkbox v-model:checked="modelUnit.submitResult">{{t('submit_result_to_zentao')}}</a-checkbox>
                  </a-form-item>


                  <a-form-item>
                    <a-button :disabled="isRunning === 'true' || !modelUnit.cmd" @click="execUnit" type="primary" class="t-btn-gap">
                      {{ t('exec') }}
                    </a-button>
                    <a-button v-if="isRunning === 'true'" @click="execStop" class="t-btn-gap">
                      {{ t('stop') }}
                    </a-button>
                  </a-form-item>

                  <a-form-item>
                    <span class="t-tips">{{t('cmd_nav')}}</span>
                  </a-form-item>
                </a-form>
              </div>

              <div class="logs-panel">
                <ScriptExecLogPage></ScriptExecLogPage>
              </div>
            </template>
          </div>
        </template>

      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  onMounted,
  ref,
  watch
} from "vue";
import {useStore} from "vuex";
import {useI18n} from "vue-i18n";

import {ScriptData} from "./store";
import {resizeWidth, resizeHeight} from "@/utils/dom";
import {Empty, notification} from "ant-design-vue";

import {MonacoOptions} from "@/utils/const";
import bus from "@/utils/eventBus"
import MonacoEditor from "@/components/Editor/MonacoEditor.vue";
import {ZentaoData} from "@/store/zentao";

import ScriptTreePage from "./component/tree.vue";
import ScriptExecLogPage from "./component/execLog.vue";
import settings from "@/config/settings";
import {get} from "@/views/workspace/service";
import {getCmdHistories, setCmdHistories} from "@/utils/cache";
import {ExecStatus} from "@/store/exec";

export default defineComponent({
  name: 'ScriptListPage',
  components: {
    ScriptTreePage, ScriptExecLogPage,
    MonacoEditor,
  },
  setup() {
    const { t } = useI18n();

    let editorRef = ref(null as any)

    const zentaoStore = useStore<{ Zentao: ZentaoData }>();
    const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

    const execStore = useStore<{ Exec: ExecStatus }>();
    const isRunning = computed<any>(() => execStore.state.Exec.isRunning);

    let modelUnit = ref({} as any)

    const scriptStore = useStore<{ Script: ScriptData }>();
    const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

    let script = computed<any>(() => scriptStore.state.Script.detail);
    let scriptCode = ref('')
    let lang = ref('')
    const editorOptions = ref(MonacoOptions)

    watch(script, () => {
      console.log('watch script', script)

      if (script.value) {
        scriptCode.value = script.value.code
        lang.value = script.value.lang
      } else {
        scriptCode.value = ''
        lang.value = ''
      }
    }, {deep: true})

    watch(currWorkspace, () => {
      console.log('watch currWorkspace', currWorkspace)

      loadCmdHistories()

      get(currWorkspace.value.id).then((json) => {
        modelUnit.value = Object.assign({cmd: json.data.cmd}, currWorkspace.value)
      })
    }, {deep: true})

    const save = () => {
      console.log('save')

      const code = editorRef.value._getValue()
      scriptStore.dispatch('Script/updateCode',
          {workspaceId: currWorkspace.value.id, path: script.value.path, code: code}).then(() => {
        notification.success({
          message: t('save_success'),
        })
      })
    }

    const execSingle = () => {
      console.log('exec', script.value)

      bus.emit(settings.eventExec, {execType: 'ztf', scripts: [script.value]});
    }
    const execUnit = () => {
      console.log('execUnit', modelUnit.value)

      if (modelUnit.value.cmd !== histories.value[histories.value.length - 1]) histories.value.push(modelUnit.value.cmd)
      if (histories.value.length > 10) histories.value = histories.value.slice(histories.value.length - 10)
      setCmdHistories(currWorkspace.value.id, histories.value)
      historyIndex.value = histories.value.length

      const data = Object.assign({execType: 'unit'}, modelUnit.value)
      bus.emit(settings.eventExec, data);
    }
    const execStop = () => {
      console.log('execStop')
      const data = Object.assign({execType: 'stop'})
      bus.emit(settings.eventExec, data);
    }

    const extract = () => {
      console.log('extract', script.value)

      scriptStore.dispatch('Script/extractScript', script.value).then(() => {
        notification.success({
          message: t('extract_success'),
        })
      }).catch(() => {
        notification.error({
          message: t('extract_fail'),
        });
      })
    }

    const histories = ref([] as any[])
    const historyIndex = ref(0)

    const loadCmdHistories = async () => {
      if (currWorkspace.value.type === 'ztf') return

      histories.value = await getCmdHistories(currWorkspace.value.id)
      historyIndex.value = histories.value.length
    }

    const up = () => {
      console.log('up')
      if (historyIndex.value > 0) historyIndex.value--
      modelUnit.value.cmd = histories.value[historyIndex.value]
    }
    const down = () => {
      console.log('down')
      if (historyIndex.value < histories.value.length - 1) historyIndex.value++
      modelUnit.value.cmd = histories.value[historyIndex.value]
    }

    onMounted(() => {
      console.log('onMounted')

      setTimeout(() => {
        resizeWidth('main', 'left', 'splitter-h', 'right', 280, 800)
        resizeHeight('right-content', 'editor-panel', 'splitter-v', 'logs-panel',
            100, 100, 90)
      }, 600)
    })

    return {
      t,
      currSite,
      currProduct,

      modelUnit,
      currWorkspace,
      script,
      scriptCode,
      lang,
      editorOptions,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,

      save,
      isRunning,
      execSingle,
      execStop,
      execUnit,
      extract,
      stop,

      histories,
      historyIndex,
      up,
      down,
      editorRef,

      // labelCol: { span: 3, offset: 12 },
      // wrapperCol: { span: 3, offset: 12 },
    }
  }

})
</script>

<style lang="less" scoped>
.script-main {
  margin: 0px;
  height: 100%;

  #main {
    display: flex;
    height: 100%;
    width: 100%;

    #left {
      width: 380px;

      height: 100%;
      padding: 0;
    }

    #splitter-h {
      width: 0px;
      border: solid 1px #D0D7DE;
      height: 100%;
      cursor: ew-resize;

      &.active {
        border-color: #a9aeb4;
      }
    }

    #right {
      flex: 1;
      height: 100%;

      .toolbar {
        padding: 4px 10px;
        height: 40px;
        text-align: right;

        .ant-btn {
          margin: 0 5px;
        }
      }

      .right-content {
        height: calc(100% - 50px);

        display: flex;
        flex-direction: column;

        #editor-panel {
          flex: 1;

          padding: 0 6px 0 8px;
          overflow: auto;
        }

        #splitter-v {
          width: 100%;
          height: 2px;
          background-color: #D0D7DE;
          cursor: ns-resize;

          &.active {
            background-color: #a9aeb4;
          }
        }

        #logs-panel {
          height: 160px;

          .ws-status {
            padding-left: 8px;
            height: 44px;
            line-height: 44px;
            color: #333333;

            &.success {
              background-color: #DAF7E9;
              svg {
                color: #DAF7E9;
              }
            }
            &.error {
              background-color: #FFD6D0;
              svg {
                color: #FFD6D0;
              }
            }

            .text {
              display: inline-block;
              margin-left: 5px;
            }
            .icon-close {
              position: absolute;
              padding: 5px;
              line-height: 34px;
              right: 15px;
              cursor: pointer;
              svg {
                font-size: 8px;
                color: #333333;
              }
            }
          }

          #logs {
            margin: 0;
            padding: 10px;
            width: 100%;
            overflow-y: auto;
            white-space: pre-wrap;
            word-wrap: break-word;
            font-family:monospace;

            height: 100%;
            &.with-status {
              height: calc(100% - 45px);
            }
          }
        }

        .logs-panel {
          height: 100%;
        }
      }
    }

    // unit test
    #right {
      .right-content {
        .unit-panel {
          padding: 3px 5px;
          height: 40px;
          .ant-row {
            margin: 0 4px;
          }
        }

        .logs-panel {
          height: calc(100% - 40px);
        }
      }
    }
  }
}
</style>

<style lang="less">
.monaco-editor {
  padding: 10px 0;
}
</style>
