<template>
  <div class="batch-run-button">
    <ButtonGroup>
      <Button @click="_handleDropDownMenuClick(execBy)" class="rounded border primary-pale" icon="run-all"
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
import {ref} from "vue";
import Button from './Button.vue';
import ButtonGroup from './ButtonGroup.vue';
import DropdownMenu from './DropdownMenu.vue';import {useI18n} from "vue-i18n";
import {getExecBy} from "@/utils/cache";
const { t } = useI18n();

const execBy = ref('opened')
getExecBy().then((execByVal) => {
  execBy.value = execByVal
})

function _handleDropDownMenuClick(event) {
    console.log('_handleDropDownMenuClick', event.key);
}
</script>

<style lang="less" scoped>
.batch-run-button {

}
</style>
