<template>
  <Toolbar>
<!--    <Button class="rounded pure" icon="search" :hint="t('search')" />-->
    <Button class="rounded pure" icon="settings" :hint="t('settings')" />
  </Toolbar>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Toolbar from './Toolbar.vue';
import {useI18n} from "vue-i18n";
const { t } = useI18n();


</script>

<style lang="less" scoped>
</style>
