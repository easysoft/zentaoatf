<template>
  <Toolbar>
    <Button class="rounded pure" icon="panel-left" hint="切换显示左侧边栏" />
    <Button class="rounded pure" icon="panel-bottom" hint="切换显示底部边栏" />
    <Button class="rounded pure" icon="panel-right" hint="切换显示右侧边栏" />
  </Toolbar>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Toolbar from './Toolbar.vue';
</script>
