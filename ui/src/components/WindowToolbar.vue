<template>
  <div class="window-toolbar">
    <Toolbar>
      <template v-if="isElectron">
        <Button v-if="!fullScreenDef" @click="fullScreen" class="rounded pure" icon="fullscreen-on" iconSize="1.5em"
                :hint="t('fullscreen')" />
        <Button v-if="fullScreenDef" @click="fullScreen" class="rounded pure" icon="fullscreen-off" iconSize="1.5em"
                :hint="t('exit_fullscreen')" />

        <Button @click="minimize" class="rounded pure" icon="window-minimize" iconSize="1.5em"
                :hint="t('minimize')" />

        <Button v-if="!maximizeDef" @click="maximize" class="rounded pure" icon="window-maximize" iconSize="1.5em"
                :hint="t('maximize')" />
        <Button v-if="maximizeDef" @click="maximize" class="rounded pure" icon="window-restore" iconSize="1.5em"
                :hint="t('restore')" />

        <Button @click="exit" class="rounded pure" icon="window-close" iconSize="1.5em"
                :hint="t('close')" />
      </template>
    </Toolbar>
  </div>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Toolbar from './Toolbar.vue';
import {ref} from "vue";
import {useI18n} from "vue-i18n";
import {useRouter} from "vue-router";
import {getElectron} from "@/utils/comm";
import settings from "@/config/settings";

const {t} = useI18n();
const router = useRouter();

const isElectron = ref(getElectron())

const fullScreenDef = ref(false)
const fullScreen = (): void => {
  console.log('fullScreen')
  fullScreenDef.value = !fullScreenDef.value

  const { ipcRenderer } = window.require('electron')
  ipcRenderer.send(settings.electronMsg, 'fullScreen')
}

const maximizeDef = ref(true)
const minimize = (): void => {
  console.log('minimize')

  const { ipcRenderer } = window.require('electron')
  ipcRenderer.send(settings.electronMsg, 'minimize')
}
const maximize = (): void => {
  console.log('maximize')

  const { ipcRenderer } = window.require('electron')
  ipcRenderer.send(settings.electronMsg, maximizeDef.value ? 'unmaximize' : 'maximize')
  maximizeDef.value = !maximizeDef.value
}

const exit = (): void => {
  console.log('exit')
  const { ipcRenderer } = window.require('electron')
  ipcRenderer.send(settings.electronMsg, 'exit')
}
</script>

<style lang="less" scoped>
  button {
    -webkit-app-region: no-drag;
  }
  .window-toolbar{
    position: fixed;
    right: 0;
  }
</style>
