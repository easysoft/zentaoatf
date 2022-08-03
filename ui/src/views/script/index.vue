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
} from "vue";
import {useStore} from "vuex";
import {useI18n} from "vue-i18n";

import {ScriptData} from "./store";
import {resizeWidth} from "@/utils/dom";
import {Empty} from "ant-design-vue";

import {ZentaoData} from "@/store/zentao";

import ScriptTreePage from "./component/tree.vue";
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
        resizeWidth('main', 'left', 'splitter-h', 'right', 380, 800)
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
    height: 100%;
    width: 100%;

    #left {
      float: left;
      width: 30%;

      height: 100%;
      padding: 0;
    }

    #splitter-h {
      float: left;
      width: 0px;
      border: solid 1px #D0D7DE;
      height: 100%;
      cursor: ew-resize;

      &.active {
        border-color: #a9aeb4;
      }
    }

    #right {
      float: left;
      width: calc(70% - 2px);
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
