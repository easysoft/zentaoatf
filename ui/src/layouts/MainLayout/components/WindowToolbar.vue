<template>
  <Toolbar>
    <Button class="rounded pure" icon="fullscreen-on" hint="全屏" />
    <Button class="rounded pure" icon="window-minimize" hint="最小化" />
    <Button class="rounded pure" icon="window-maximize" hint="最大化" />
    <Button class="rounded pure" icon="window-close" hint="关闭窗口" />
  </Toolbar>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Toolbar from './Toolbar.vue';
</script>
