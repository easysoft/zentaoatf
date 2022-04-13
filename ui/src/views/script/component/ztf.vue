<template>
  <div class="ztf-script-main">

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

  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, ref, watch} from "vue";
import {useStore} from "vuex";
import {useI18n} from "vue-i18n";

import {ScriptData} from "../store";
import {resizeHeight, resizeWidth} from "@/utils/dom";
import {Empty, notification} from "ant-design-vue";

import {MonacoOptions} from "@/utils/const";
import bus from "@/utils/eventBus"
import MonacoEditor from "@/components/Editor/MonacoEditor.vue";
import {ZentaoData} from "@/store/zentao";

import ScriptExecLogPage from "./execLog.vue";
import settings from "@/config/settings";
import {ExecStatus} from "@/store/exec";

export default defineComponent({
  name: 'ZtfScriptPage',
  components: {
    ScriptExecLogPage,
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

    const scriptStore = useStore<{ Script: ScriptData }>();
    const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

    let script = computed<any>(() => scriptStore.state.Script.detail);
    let scriptCode = ref('')
    let lang = ref('')
    const editorOptions = ref(MonacoOptions)

    watch(script, () => {
      console.log('watch script', script)

      if (script.value) {
        scriptCode.value = script.value.code ? script.value.code : t('empty')
        lang.value = script.value.lang
      } else {
        scriptCode.value = ''
        lang.value = ''
      }
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

    onMounted(() => {
      console.log('onMounted')

      setTimeout(() => {
        resizeHeight('right-content', 'editor-panel', 'splitter-v', 'logs-panel',
            100, 100, 90)
      }, 600)
    })

    return {
      t,
      currSite,
      currProduct,

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
      extract,
      stop,

      editorRef,
    }
  }

})
</script>

<style lang="less" scoped>

.ztf-script-main {
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
    height: calc(100%);

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
        font-family: monospace;

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

</style>

<style lang="less">
.monaco-editor {
  padding: 10px 0;
}
</style>
