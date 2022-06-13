<template>
  <main id="main" class="column single surface-light relative no-overflow">
    <modals-container></modals-container>
    <Navbar class="flex-none" />

    <div id="mainContent" class="flex-auto">
      <Splitpanes id="mainRow" ref="mainRow">
        <Pane :size="20" id="leftPane" :min-size="minLeftPane" :max-size="50">
          <WorkDirPanel />
        </Pane>

        <Pane id="centerPane" :min-size="paneMinSize">
          <Splitpanes id="centerColumn" ref="centerColumn" horizontal v-on:resized="onSplitpanesResized($event)">
            <Pane id="tabsPane" :size='globalStore.getters["global/editorPaneSize"]'>
              <TabsContainer class="height-full" />
            </Pane>
            <Pane v-show="showLogPanel" :size='globalStore.getters["global/logPaneSize"]' id="bottomPane" :min-size="minBottomPane" :max-size="50">
              <LogPanel />
            </Pane>
          </Splitpanes>
        </Pane>

        <Pane :size="20" id="rightPane" :min-size="minRightPane" :max-size="80">
          <ResultListPanel />
        </Pane>
      </Splitpanes>
    </div>

    <Websocket></Websocket>
  </main>
</template>

<script setup lang="ts">
import './style/index.less';
import 'splitpanes/dist/splitpanes.css'

import { Splitpanes, Pane, PaneProps } from 'splitpanes';
import Navbar from '@/components/Navbar.vue';
import WorkDirPanel from '@/views/script/WorkDirPanel.vue';
import LogPanel from '@/views/exec/LogPanel.vue';
import TabsContainer from '@/components/TabsContainer.vue';
import ResultListPanel from '@/views/result/ResultListPanel.vue';
import Websocket from '@/components/Websocket.vue';
import settings from "@/config/settings";
import {onBeforeUnmount, onMounted, ref} from "vue";
import bus from "@/utils/eventBus";
import notification from "@/utils/notification";
import { useI18n } from "vue-i18n";
import {StateType} from "@/store/global"
import { useStore } from 'vuex';
import { useHeightToPercent, useWidthToPercent } from '@/components/hooks/use-pixels-to-percent';

const { t } = useI18n();
const showLogPanel = ref(false);
const mainRow = ref();
const centerColumn = ref();
const [minLeftPane, minRightPane] = useWidthToPercent(mainRow, 220, 200);
const [minBottomPane] = useHeightToPercent(centerColumn, 70);

const paneMinSize = ref(0)

const globalStore = useStore<{global: StateType}>()

const onExecStartEvent = () => {
  console.log('onExecStartEvent')
  showLogPanel.value = true
}

const notify = (result: any) => {
  if (!result.httpCode) result.httpCode = 100
  const msg = result.httpCode === 200 ? t('biz_'+result.resultCode) : t('http_'+result.httpCode)
  const desc = result.resultMsg ? result.resultMsg : ''

  notification.error({
    message: msg,
    description: desc,
  });
}

const onSplitpanesResized = (evt : PaneProps[]) => {
  let logpaneProps = evt[1];
  globalStore.commit("global/setLogPaneSize", logpaneProps.size)
}

onMounted(() => {
  console.log('onMounted ztf')
  bus.on(settings.eventExec, onExecStartEvent)
  bus.on(settings.eventNotify, notify);

  const resize_ob = new ResizeObserver((entries) => {
    console.log(1, entries)
    const width = entries[0].contentRect.width
    console.log(2, width)

    paneMinSize.value = (660 * 100) / width
    console.log(3, paneMinSize.value)
  })
  resize_ob.observe(mainRow.value.container)
})
onBeforeUnmount( () => {
  bus.off(settings.eventExec, onExecStartEvent)
  bus.off(settings.eventNotify, notify);
})

</script>

<style lang="less" scoped>
#main {
  height: 100vh;
  width: 100vw;
  min-width: 1100px;

  #mainContent {
    height: calc(100% - 40px);
    -webkit-app-region: no-drag;
  }
}
#leftPane {
  min-width: var(--pane-left-min-width, 230px);
}

#centerPane {
  min-width: 500px;

  #tabsPane {
    min-width: var(--pane-tabs-min-width, 100px);
  }
  #bottomPane {
    min-width: var(--pane-bottom-min-width, 100px);
  }
}

#rightPane {
  min-width: var(--pane-right-min-width, 200px);
}

</style>
