<template>
  <div class="script-main">
    <div id="main">
      <div id="left">
        <ScriptTreePage></ScriptTreePage>
      </div>

      <div id="splitter-h"></div>

      <div id="right">
        <ZtfScriptPage v-if="currWorkspace?.type === 'ztf'"></ZtfScriptPage>
        <UnitScriptPage v-if="currWorkspace?.type !== 'ztf'"></UnitScriptPage>
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
import ZtfScriptPage from "./component/ztf.vue"
import UnitScriptPage from "./component/unit.vue"

export default defineComponent({
  name: 'ScriptListPage',
  components: {
    ScriptTreePage,
    ZtfScriptPage,
    UnitScriptPage,
  },
  setup() {
    const { t } = useI18n();

    const zentaoStore = useStore<{ Zentao: ZentaoData }>();
    const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
    const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

    const scriptStore = useStore<{ Script: ScriptData }>();
    const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

    onMounted(() => {
      console.log('onMounted')
      setTimeout(() => {
        resizeWidth('main', 'left', 'splitter-h', 'right', 280, 800)
      }, 600)
    })

    return {
      t,
      currSite,
      currProduct,

      currWorkspace,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
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
    }
  }
}
</style>

<style lang="less">
.monaco-editor {
  padding: 10px 0;
}
</style>
