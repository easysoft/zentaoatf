<template>
  <div class="unit-script-main">
    <!-- Codes -->
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

    <!-- Exec -->
    <template v-if="!script">
      <div class="unit-panel">
        <a-form :model="modelUnit" layout="inline">
          <a-form-item :label="t('test_cmd')">
            <a-input
                v-model:value="modelUnit.cmd"
                @keydown="keydown"
                style="width:380px;"/>
          </a-form-item>

          <a-form-item v-if="currProduct.id">
            <a-checkbox v-model:checked="modelUnit.submitResult">{{ t('submit_result_to_zentao') }}</a-checkbox>
          </a-form-item>

          <a-form-item>
            <a-button :disabled="isRunning === 'true' || !modelUnit.cmd" @click="execUnit" type="primary"
                      class="t-btn-gap">
              {{ t('exec') }}
            </a-button>
            <a-button v-if="isRunning === 'true'" @click="execStop" class="t-btn-gap">
              {{ t('stop') }}
            </a-button>
          </a-form-item>

          <a-form-item>
            <span class="t-tips">{{ t('cmd_nav') }}</span>
          </a-form-item>
        </a-form>
      </div>

      <div class="logs-panel">
        <ScriptExecLogPage></ScriptExecLogPage>
      </div>
    </template>
  </div>
</template>

<script lang="ts">
import {computed, defineComponent, onMounted, ref, watch} from "vue";
import {useStore} from "vuex";
import {useI18n} from "vue-i18n";

import {ScriptData} from "../store";
import {Empty, notification} from "ant-design-vue";

import {MonacoOptions} from "@/utils/const";
import bus from "@/utils/eventBus"
import MonacoEditor from "@/components/Editor/MonacoEditor.vue";
import {ZentaoData} from "@/store/zentao";

import ScriptExecLogPage from "./execLog.vue";
import settings from "@/config/settings";
import {get} from "@/views/workspace/service";
import {getCmdHistories, setCmdHistories} from "@/utils/cache";
import {ExecStatus} from "@/store/exec";

export default defineComponent({
  name: 'UnitScriptPage',
  components: {
    ScriptExecLogPage,
    MonacoEditor,
  },
  setup() {
    const {t} = useI18n();

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
      console.log('watch script', script.value)

      if (script.value) {
        scriptCode.value = script.value.code ? script.value.code : t('empty')
        lang.value = script.value.lang
      } else {
        scriptCode.value = ''
        lang.value = ''
      }
    }, {deep: true})

    watch(currWorkspace, () => {
      console.log('watch currWorkspace', currWorkspace)

      loadCmdHistories()
      loadWorkspaceCmd()
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

    const histories = ref([] as any[])
    const historyIndex = ref(0)

    const loadCmdHistories = async () => {
      histories.value = await getCmdHistories(currWorkspace.value.id)
      historyIndex.value = histories.value? histories.value.length : 0
    }
    const loadWorkspaceCmd = async () => {
      get(currWorkspace.value.id).then((json) => {
        modelUnit.value = Object.assign({cmd: json.data.cmd}, currWorkspace.value)
      })
    }

    const keydown = (e) => {
      console.log('keydown', e.code)
      if (e.code === 'ArrowUp') {
        if (historyIndex.value > 0) historyIndex.value--
        modelUnit.value.cmd = histories.value[historyIndex.value]
      } else if (e.code === 'ArrowDown') {
        if (historyIndex.value < histories.value.length - 1) historyIndex.value++
        modelUnit.value.cmd = histories.value[historyIndex.value]
      }
    }

    onMounted(() => {
      console.log('onMounted')

      loadCmdHistories()
      loadWorkspaceCmd()
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
      execStop,
      execUnit,
      stop,

      histories,
      historyIndex,
      keydown,
      editorRef,
    }
  }

})
</script>

<style lang="less" scoped>

.unit-script-main {
  height: 100%;

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

</style>

<style lang="less">
.monaco-editor {
  padding: 10px 0;
}
</style>
