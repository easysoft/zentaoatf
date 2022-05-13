<template>
  <Panel>
    <template #header>
      <ButtonList :gap="2" class="workdir-panel-nav">
        <Button id="displayByMenuToggle"
                :label="t('by_' + displayBy)"
                labelClass="strong"
                class="rounded pure padding-sm-h"
                suffix-icon="caret-down"/>

        <Dropdown :items="displayTypes"
                  :checkedKey="displayBy"
                  @click="onDisplayByChanged"
                  toggle="#displayByMenuToggle">
        </Dropdown>

        <Dropdown>
          <Button class="rounded pure padding-sm-h" label="按套件" suffix-icon="caret-down" />
          <template #menu>
            <List>
              <ListItem :checked="true">{{ t('by_suite') }}</ListItem>
              <ListItem :checked="false">{{ t('by_task') }}</ListItem>
            </List>
          </template>
        </Dropdown>
      </ButtonList>
    </template>

    <template #toolbar-buttons>
      <Button @click="expandAllOrNot" class="rounded pure" hint="t('collapse')" icon="subtract-square-multiple" iconSize="1.4em" />
      <Button @click="expandAllOrNot" class="rounded pure" hint="t('expand_all')" icon="dismiss-square-multiple" iconSize="1.4em" />

      <Button class="rounded pure" :hint="t('batch_select')" icon="select-all-on" />
      <Button class="rounded pure" :hint="t('create_workspace')" icon="folder-add" />
      <Button class="rounded pure" :hint="t('more_actions')" icon="more-vert" />
    </template>

    <WorkDir />
  </Panel>
</template>

<script setup lang="ts">
import {computed, defineComponent, onMounted, onUnmounted, ref, watch} from "vue";

import Panel from './Panel.vue';
import Button from './Button.vue';
import ButtonList from './ButtonList.vue';
import WorkDir from './WorkDir.vue';
import Dropdown from './Dropdown.vue';
import List from './List.vue';
import ListItem from './ListItem.vue';
import {getScriptDisplayBy, setScriptDisplayBy} from "@/utils/cache";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import debounce from "lodash.debounce";
import {useI18n} from "vue-i18n";
import {ScriptData} from "@/views/script/store";
const { t } = useI18n();

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const store = useStore<{ Script: ScriptData }>();

let displayBy = ref('workspace')
let isExpand = ref(false);
const filerType = ref('')
const filerValue = ref('')

const displayTypes = ref([
  {key: 'workspace', title: t('by_workspace')},
  {key: 'module', title: t('by_module')},
])

const loadDisplayBy = async () => {
  displayBy.value = await getScriptDisplayBy(currSite.value.id, currProduct.value.id)
}

const initData = debounce(async () => {
  console.log('init')
  if (!currSite.value.id) return

  await loadDisplayBy()
}, 50)

watch(currProduct, () => {
  console.log('watch currProduct', currProduct.value.id)
  initData()
}, {deep: true})

// only do it when switch from another pages, otherwise will called by watching currProduct method.
if (filerValue.value.length === 0) initData()

const onDisplayByChanged = () => {
  console.log('onDisplayBy')
  if (displayBy.value === 'workspace') displayBy.value = 'module'
  else displayBy.value = 'workspace'

  setScriptDisplayBy(displayBy.value, currSite.value.id, currProduct.value.id)

  loadScripts()
}

const loadScripts = async () => {
  console.log(`filerType: ${filerType.value}, filerValue: ${filerValue.value}`)

  const params = {displayBy: displayBy.value, filerType: filerType.value, filerValue: filerValue.value} as any
  store.dispatch('Script/listScript', params)
}

const expandAllOrNot = () => {
  console.log('expandAllOrNot')
}

</script>

<style scoped>
.workdir-panel-nav {
  margin-left: -6px;
}
</style>
