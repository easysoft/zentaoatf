<template>
  <Toolbar :id="isElectron ? 'setting-btn':''">
    <Button class="rounded pure" icon="settings" iconSize="1.5em" :hint="t('settings')" @click="openSettings" />
  </Toolbar>
  <SettingsModal
    v-if="showSettingsModal"
    :show="showSettingsModal"
    @cancel="settingsModalClose"
    :showOkBtn="false"
    :showCancelBtn="false"
  />
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Toolbar from './Toolbar.vue';
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import SettingsModal from '@/views/settings/SettingsModal.vue';
import { ref } from "vue";
import { getElectron } from "@/utils/comm";

const { t } = useI18n();
const store = useStore<{ Zentao: ZentaoData }>();
const isElectron = ref(getElectron());

const showSettingsModal = ref(false);

const openSettings = () => {
    showSettingsModal.value = true;
}

const settingsModalClose = () => {
    showSettingsModal.value = false;
}

</script>
<style>
#setting-btn{
    position: fixed;
    right: 100px;
}
</style>