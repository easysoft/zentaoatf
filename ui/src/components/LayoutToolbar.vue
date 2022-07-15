<template>
  <Toolbar>
    <Button class="rounded pure" icon="panel-left" :hint="t('switch_display_left_side')" />
    <Button class="rounded pure" icon="panel-bottom" :hint="t('switch_display_bottom_side')" />
    <Button class="rounded pure" icon="panel-right" :hint="t('switch_display_right_side')" />
  </Toolbar>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Toolbar from './Toolbar.vue';
import {useI18n} from "vue-i18n";
const { t } = useI18n();
</script>
