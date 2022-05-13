<template>
  <div class="tab-page-script">
    <div>TabPage: 测试脚本文件</div>
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
