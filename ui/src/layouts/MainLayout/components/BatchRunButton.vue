<template>
  <ButtonGroup>
    <Button class="rounded border primary-pale" icon="run-all"
            :label="t('exec_selected')" />
    <Button class="rounded border primary-pale padding-0" iconClass="muted" icon="caret-down" iconSize="1em" style="width: 20px" id="batchRunMenuToggle" />
  </ButtonGroup>
  <DropdownMenu
    :items="[
      {id: 'all', name: t('exec_all')},
      {id: 'last', name: t('exec_previous')},
      {id: 'select', name: t('exec_selected')},
      {id: 'open', name: t('exec_opened')}
    ]"
    keyName="id"
    checkedKey="select"
    activeKey="select"
    :replaceFields="{title: 'name'}"
    toggle="#batchRunMenuToggle"
    position="bottom-right"
    @click="_handleDropDownMenuClick"
  />
</template>

<script setup lang="ts">
import Button from './Button.vue';
import ButtonGroup from './ButtonGroup.vue';
import DropdownMenu from './DropdownMenu.vue';import {useI18n} from "vue-i18n";
const { t } = useI18n();

function _handleDropDownMenuClick(event) {
    console.log('_handleDropDownMenuClick', event);
}
</script>
