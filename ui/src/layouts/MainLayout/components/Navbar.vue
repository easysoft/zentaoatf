<template>
  <nav id="navbar" class="surface row divider gap-lg align-center padding-h">
    <Logo />
    <SiteNav class="space-left" />
    <BatchRunButton class="space-left" />
    <LastTextInfo />
    <div class="gap flex-auto"></div>
    <AppToolbar />
    <LayoutToolbar />
    <WindowToolbar />
  </nav>
</template>

<script setup lang="ts">
import Logo from './Logo.vue';
import WindowToolbar from './WindowToolbar.vue';
import AppToolbar from './AppToolbar.vue';
import LayoutToolbar from './LayoutToolbar.vue';
import SiteNav from './SiteNav.vue';
import BatchRunButton from './BatchRunButton.vue';
import LastTextInfo from './LastTextInfo.vue';
</script>

<style>
#navbar {
  height: var(--navbar-height, 40px);
}
</style>
