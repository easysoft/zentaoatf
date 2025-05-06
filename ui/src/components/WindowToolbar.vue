<template>
  <div class="window-toolbar">
    <Toolbar>
      <template v-if="isElectron">
        <div class="version-container">
          <span class="version-text">V{{version}}</span>
        </div>
        <div class="window-controls">
          <Button v-if="!fullScreenDef" @click="fullScreen" class="rounded pure" icon="fullscreen-on" iconSize="1.5em"
                  :hint="t('fullscreen')" />
          <Button v-if="fullScreenDef" @click="fullScreen" class="rounded pure" icon="fullscreen-off" iconSize="1.5em"
                  :hint="t('exitFullscreen')" />

          <Button @click="minimize" class="rounded pure" icon="window-minimize" iconSize="1.5em"
                  :hint="t('minimize')" />

          <Button v-if="!maximizeDef" @click="maximize" class="rounded pure" icon="window-maximize" iconSize="1.5em"
                  :hint="t('maximize')" />
          <Button v-if="maximizeDef" @click="maximize" class="rounded pure" icon="window-restore" iconSize="1.5em"
                  :hint="t('restore')" />

          <Button @click="exit" class="rounded pure" icon="window-close" iconSize="1.5em"
                  :hint="t('close')" />
        </div>
      </template>
    </Toolbar>
  </div>
</template>

<style lang="less" scoped>
  button {
    -webkit-app-region: no-drag;
  }
  .window-toolbar {
    position: fixed;
    right: 0;
    display: flex;
    align-items: center;
  }
  .version-container {
    position: fixed;
    display: flex;
    align-items: left;
    right: 110px;
    z-index: 1000;
  }
  .window-controls {
    display: flex;
    align-items: right;
  }
  .version-text {
    font-size: 12px;
    color: green;
    -webkit-app-region: no-drag;
    min-width: 80px;
    font-weight: bold;
  }
</style>

<script setup lang="ts">
declare global {
  interface Window {
    require: (module: string) => any;
  }
}
import Button from './Button.vue';
import Toolbar from './Toolbar.vue';
import {ref, onMounted} from "vue";
import {useI18n} from "vue-i18n";
import {useRouter} from "vue-router";
import {getElectron} from "@/utils/comm";
import settings from "@/config/settings";

const version = ref('');
const {t} = useI18n();
const router = useRouter();
const isElectron = ref(getElectron())
let ipcRenderer: any = null;

onMounted(() => {
  if (isElectron.value) {
    ipcRenderer = window.require('electron').ipcRenderer;

    ipcRenderer.on('version', (event: any, ver: string) => {
      version.value = ver;
    });

    try {
      setTimeout(() => {
          ipcRenderer.send(settings.electronMsg, 'getVersion');
      }, 100);
    } catch (error) {
      console.error('Error in version IPC:', error);
    }
  }
})

const fullScreenDef = ref(false)
const fullScreen = (): void => {
  console.log('fullScreen')
  fullScreenDef.value = !fullScreenDef.value
  ipcRenderer?.send(settings.electronMsg, 'fullScreen')
}

const maximizeDef = ref(true)
const minimize = (): void => {
  console.log('minimize')
  ipcRenderer?.send(settings.electronMsg, 'minimize')
}

const maximize = (): void => {
  console.log('maximize')
  ipcRenderer?.send(settings.electronMsg, maximizeDef.value ? 'unmaximize' : 'maximize')
  maximizeDef.value = !maximizeDef.value
}

const exit = (): void => {
  console.log('exit')
  ipcRenderer?.send(settings.electronMsg, 'exit')
}
</script>
