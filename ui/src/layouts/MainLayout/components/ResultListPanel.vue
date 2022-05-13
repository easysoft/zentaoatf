<template>
  <Panel :title="t('test_result')">
    <template #toolbar-buttons>
      <Button class="rounded pure" :hint="t('more_actions')" icon="more-vert" />
    </template>
    <ResultList />
  </Panel>
</template>

<script setup lang="ts">
import Panel from './Panel.vue';
import Button from './Button.vue';
import ResultList from './ResultList.vue';
import {useI18n} from "vue-i18n";
const { t } = useI18n();

</script>

<style scoped>
</style>
