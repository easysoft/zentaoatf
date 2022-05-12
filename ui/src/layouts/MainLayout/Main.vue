<template>
  <main id="main" class="column single surface-light relative no-overflow">
    <Navbar class="flex-none" />
    <div id="mainContent" class="flex-auto">
      <Splitpanes id="mainRow">
        <Pane :size="20" id="leftPane"><WorkDirPanel /></Pane>
        <Pane id="centerPane">
          <Splitpanes id="centerColumn" horizontal>
            <Pane id="tabsPane">
              <TabsContainer />
            </Pane>
            <Pane :size="20" id="bottomPane">
              <LogPanel />
            </Pane>
          </Splitpanes>
        </Pane>
        <Pane :size="20" id="rightPane">
          <ResultListPanel />
        </Pane>
      </Splitpanes>
    </div>
  </main>
</template>

<script setup lang="ts">
import './style/var.css';
import './style/helpers.css';
import './style/main.css';
import './style/grid.less'
import 'splitpanes/dist/splitpanes.css'

import { Splitpanes, Pane } from 'splitpanes';
import Navbar from './components/Navbar.vue';
import WorkDirPanel from './components/WorkDirPanel.vue';
import LogPanel from './components/LogPanel.vue';
import TabsContainer from './components/TabsContainer.vue';
import ResultListPanel from './components/ResultListPanel.vue';
</script>

<style lang="less" scoped>
#main {
  height: 100vh;
  width: 100vw;

  #mainContent {
    -webkit-app-region: no-drag;
  }
}
#leftPane {
  min-width: var(--pane-left-min-width, 200px);
}
#centerPane {
  min-width: var(--pane-center-min-width, 100px);
}
#rightPane {
  min-width: var(--pane-right-min-width, 200px);
}
#tabsPane {
  min-width: var(--pane-tabs-min-width, 100px);
}
#bottomPane {
  min-width: var(--pane-bottom-min-width, 200px);
}
</style>
