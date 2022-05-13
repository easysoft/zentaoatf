<template>
  <div class="tab-page-sites">
    <div>TabPage: 测试结果页面</div>
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
