<template>
  <Panel title="测试结果">
    <template #toolbar-buttons>
      <Button class="rounded pure" hint="更多操作" icon="more-vert" />
    </template>
    <ResultList />
  </Panel>
</template>

<script setup lang="ts">
import Panel from './Panel.vue';
import Button from './Button.vue';
import ResultList from './ResultList.vue';
</script>

<style scoped>
</style>
