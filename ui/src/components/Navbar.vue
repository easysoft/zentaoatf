<template>
  <nav id="navbar" class="surface row divider gap-lg align-center padding-h">
    <Logo />
    <SiteNav />

    <BatchRunButton />

    <div class="gap flex-auto"></div>
    <AppToolbar />
    <WindowToolbar />
  </nav>
</template>

<script setup lang="ts">
import Logo from './Logo.vue';
import WindowToolbar from './WindowToolbar.vue';
import AppToolbar from './AppToolbar.vue';
import SiteNav from '@/views/site/SiteNav.vue';
import BatchRunButton from '@/views/exec/BatchRunButton.vue';
</script>

<style>
#navbar {
  height: var(--navbar-height, 40px);
}

#navbar button {
  -webkit-app-region: no-drag;
}
</style>
