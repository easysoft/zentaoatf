<template>
  <div class="batch-run-button">
    <ButtonGroup>
      <Button @click="_handleDropDownMenuClick" class="rounded border primary-pale" icon="run-all"
              :label="t('exec_'+execBy)" />
      <Button id="batchRunMenuToggle"
              class="rounded border primary-pale padding-0"
              iconClass="muted"
              icon="caret-down"
              iconSize="1em"
              style="width: 20px" />
    </ButtonGroup>

    <DropdownMenu
        :items="[
          {key: 'all', title: t('exec_all')},
          {key: 'previous', title: t('exec_previous')},
          {key: 'opened', title: t('exec_opened')}
        ]"
        :checkedKey="execBy"
        :activeKey="execBy"
        toggle="#batchRunMenuToggle"
        position="bottom-right"
        @click="_handleDropDownMenuClick"
    />
  </div>
</template>

<script setup lang="ts">
import {computed, ref} from "vue";
import Button from '@/components/Button.vue';
import ButtonGroup from '@/components/ButtonGroup.vue';
import DropdownMenu from '@/components/DropdownMenu.vue';
import {useI18n} from "vue-i18n";
import {getExecBy} from "@/utils/cache";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {useStore} from "vuex";
import {ScriptData} from "@/views/script/store";
import {TabsData} from "@/store/tabs";
import {TabNavItem} from "@/components/TabsNav.vue";
const { t } = useI18n();

const store = useStore<{ tabs: TabsData, Script: ScriptData }>();
const tabs = computed<TabNavItem[]>(() => {
  return store.getters['tabs/list'];
});

const currWorkspace = computed<any>(() => store.state.Script.currWorkspace);
const treeDataMap = computed<any>(() => store.state.Script.treeDataMap);
const selectedNodes = computed<any>(() => store.state.Script.checkedNodes);

const execBy = ref('opened');
(async function () {
  const execByVal = await getExecBy()
  execBy.value = execByVal
})();

function _handleDropDownMenuClick(event) {
  if (event.key) execBy.value = event.key
  console.log('_handleDropDownMenuClick', execBy.value);

  if (execBy.value === 'selected') {
    console.log(selectedNodes.value);

    let arr = [] as string[]
    Object.keys(selectedNodes.value).forEach(k => {
      if (selectedNodes.value[k] === true && treeDataMap.value[k]?.type === 'file') {
        arr.push(treeDataMap.value[k])
      }
    })

    bus.emit(settings.eventExec, { execType: 'ztf', scripts: arr });

  } else if (execBy.value === 'opened') {
    const openedScripts = getOpenedScripts()
    bus.emit(settings.eventExec, { execType: 'ztf', scripts: openedScripts });

  } else if (execBy.value === 'previous') {
    bus.emit(settings.eventExec, { execType: 'previous' });

  } else if (execBy.value === 'all') {
    let arr = [] as string[]
    Object.keys(treeDataMap.value).forEach(k => {
      if (treeDataMap.value[k].workspaceType === 'ztf' && treeDataMap.value[k]?.type === 'file') {
        arr.push(treeDataMap.value[k])
      }
    })

    bus.emit(settings.eventExec, { execType: 'ztf', scripts: arr });
  }
}

const getOpenedScripts = () => {
  return tabs.value.filter((item) => {
    return item.type === 'script'
  }).map((script: any) => {
    return { path: script.data, workspaceId: currWorkspace.value.id }
  });
}

</script>

<style lang="less" scoped>
.batch-run-button {

}
</style>
