<template>
  <main id="main" class="column single surface-light relative no-overflow">
    <Navbar class="flex-none" />

    <div id="mainContent" class="flex-auto">
      <Splitpanes id="mainRow">
        <Pane :size="20" id="leftPane">
          <WorkDirPanel />
        </Pane>

        <Pane id="centerPane">

          <Splitpanes id="centerColumn" horizontal>
            <Pane id="tabsPane">
              <TabsContainer class="height-full" />
            </Pane>
            <Pane v-show="showLogPanel" :size="30" id="bottomPane">
              <LogPanel />
            </Pane>
          </Splitpanes>

        </Pane>

        <Pane :size="20" id="rightPane">
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

import { Splitpanes, Pane } from 'splitpanes';
import Navbar from './components/Navbar.vue';
import WorkDirPanel from './components/WorkDirPanel.vue';
import LogPanel from './components/LogPanel.vue';
import TabsContainer from './components/TabsContainer.vue';
import ResultListPanel from './components/ResultListPanel.vue';
import Websocket from './components/Websocket.vue';
import settings from "@/config/settings";
import {onBeforeUnmount, onMounted, ref} from "vue";
import bus from "@/utils/eventBus";
import {notification} from "ant-design-vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const showLogPanel = ref(false)

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

onMounted(() => {
  console.log('onMounted ztf')
  bus.on(settings.eventExec, onExecStartEvent)
  bus.on(settings.eventNotify, notify);
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

  #mainContent {
    height: calc(100% - 40px);
    -webkit-app-region: no-drag;
  }
}
#leftPane {
  min-width: var(--pane-left-min-width, 200px);
}

#centerPane {
  min-width: var(--pane-center-min-width, 100px);

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
