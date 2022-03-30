<template>
  <div class="script-main">
    <div v-if="currProduct.id" id="main">
      <div id="left">
        <ScriptTreePage></ScriptTreePage>
      </div>

      <div id="splitter-h"></div>

      <div id="right">
        <div class="toolbar">
          <template v-if="scriptCode !== ''">
            <a-button @click="execSingle" type="primary">{{ t('exec') }}</a-button>

            <a-button @click="extract">{{ t('extract_step') }}</a-button>
          </template>
        </div>

        <div id="right-content">
          <template v-if="script">
            <div id="editor-panel">
              <MonacoEditor
                  v-if="scriptCode !== ''"
                  class="editor"
                  :value="scriptCode"
                  :language="lang"
                  :options="editorOptions"
              />
            </div>

            <div id="splitter-v"></div>

            <div id="logs-panel">
              <ScriptExecLogPage></ScriptExecLogPage>
            </div>
          </template>

          <template v-if="!script">
            <div class="logs-panel">
              <ScriptExecLogPage></ScriptExecLogPage>
            </div>
          </template>
        </div>

      </div>
    </div>
    <div v-if="!currProduct.id">
      <a-empty :image="simpleImage"/>
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
import SyncFromZentao from "./component/syncFromZentao.vue";

export default defineComponent({
  name: 'ScriptListPage',
  components: {
    ScriptTreePage, ScriptExecLogPage,
    MonacoEditor,
  },
  setup() {
    const { t } = useI18n();

    const zentaoStore = useStore<{ zentao: ZentaoData }>();
    const currSite = computed<any>(() => zentaoStore.state.zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.zentao.currProduct);

    let tree = ref(null)

    const scriptStore = useStore<{ script: ScriptData }>();
    let script = computed<any>(() => scriptStore.state.script.detail);
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

    const execSingle = () => {
      console.log('exec', script.value)

      const data = [script.value]
      bus.emit(settings.eventExec, data);
    }

    const extract = () => {
      console.log('extract', script.value)

      scriptStore.dispatch('script/extractScript', script.value).then(() => {
        notification.success({
          message: t('extract_success'),
        })
      }).catch(() => {
        notification.error({
          message: t('extract_fail'),
        });
      })
    }

    onMounted(() => {
      console.log('onMounted', tree)

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

      tree,
      script,
      scriptCode,
      lang,
      editorOptions,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,

      execSingle,
      extract,
      stop,
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

      #right-content {
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
  }
}
</style>

<style lang="less">
.monaco-editor {
  padding: 10px 0;
}
</style>
