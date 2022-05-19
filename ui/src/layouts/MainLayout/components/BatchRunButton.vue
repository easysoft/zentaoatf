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
          {key: 'selected', title: t('exec_selected')},
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
import Button from './Button.vue';
import ButtonGroup from './ButtonGroup.vue';
import DropdownMenu from './DropdownMenu.vue';import {useI18n} from "vue-i18n";
import {getExecBy} from "@/utils/cache";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import {useStore} from "vuex";
import {ScriptData} from "@/views/script/store";
import {TabsData} from "@/store/tabs";
import {TabNavItem} from "@/layouts/MainLayout/components/TabsNav.vue";
const { t } = useI18n();

const store = useStore<{ tabs: TabsData }>();
const tabs = computed<TabNavItem[]>(() => {
  return store.getters['tabs/list'];
});

const scriptStore = useStore<{ Script: ScriptData }>();
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);
const treeDataMap = computed<any>(() => scriptStore.state.Script.treeDataMap);
const selectedNodes = computed<any>(() => scriptStore.state.Script.checkedNodes);

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
    console.log(arr);
    bus.emit(settings.eventExec, { execType: 'ztf', scripts: arr });

  } else if (execBy.value === 'opened') {
    const openedScripts = getOpenedScripts()
    bus.emit(settings.eventExec, { execType: 'ztf', scripts: openedScripts });

  }

  return
}

const getOpenedScripts = () => {
  const openedScripts = tabs.value.filter((item, index) => {
    return item.type === 'script'
  }).map((script: any) => {
    return { path: script.data, workspaceId: currWorkspace.value.id }
  });

  return openedScripts
}

</script>

<style lang="less" scoped>
.batch-run-button {

}
</style>
