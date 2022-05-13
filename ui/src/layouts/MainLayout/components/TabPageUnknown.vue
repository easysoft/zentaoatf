<template>
  <div class="tab-page-sites">
    <div>TabPage: 不支持的标签页类型</div>
    <div>{{JSON.stringify(tab)}}</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import { PageTab } from "@/store/tabs";

defineProps<{
    tab: PageTab
}>();
</script>
