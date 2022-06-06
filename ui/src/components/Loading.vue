<template>
  <div class="center">
    <div class="padding">
      <div class="spinner"></div>
      <slot>
        <div class="muted space-sm-top">{{title ?? t('loading')}}</div>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import {defineProps} from 'vue';
import {useI18n} from "vue-i18n";

const {t} = useI18n();

defineProps<{title?: string}>();
</script>
