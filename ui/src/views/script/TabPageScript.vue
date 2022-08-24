<template>
    <div class="ztf-script-main">
        <div id="editor-panel" class="tab-page-script editor-panel">
            <MonacoEditor v-if="scriptCode !== '' && scriptCode !== ScriptFileNotExist" v-model:value="scriptCode"
                :language="lang" :options="editorOptions" class="editor" ref="editorRef" @change="editorChange" />
        </div>
    </div>
</template>

<script setup lang="ts">
import notification from "@/utils/notification";
import { defineProps, defineExpose, computed, ref, watch, onMounted, onBeforeUnmount } from "vue";
import { PageTab, TabsData } from "@/store/tabs";
import { useStore } from "vuex";
import { ScriptData } from "@/views/script/store";
import { MonacoOptions, ScriptFileNotExist } from "@/utils/const";
import { resizeHeight } from "@/utils/dom";
import { useI18n } from "vue-i18n";
import MonacoEditor from "@/components/MonacoEditor.vue";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";
import { StateType as GlobalData } from "@/store/global";

const { t } = useI18n();
const props = defineProps<{
    tab: PageTab
}>();

const store = useStore<{ Script: ScriptData, global: GlobalData,tabs: TabsData }>();
const global = computed<any>(() => store.state.global.tabIdToWorkspaceIdMap);
const script = computed<any>(() => store.state.Script.detail);
const scriptCode = ref('')
const currentScript = ref({} as any)

const lang = ref('')
const editorOptions = ref(MonacoOptions)
const editorRef = ref<InstanceType<typeof MonacoEditor>>()
const activeID = computed((): string => {
  return store.state.tabs.activeID;
});

watch(script, () => {
    if(script.value == undefined || (currentScript.value.path!==undefined && currentScript.value.path !== script.value.path)){
        return
    }
    console.log('watch script', script)
    if (script.value) {
        if (script.value.code === ScriptFileNotExist) {
            scriptCode.value = ScriptFileNotExist
            lang.value = ''

            return
        }
        if(scriptCode.value == '' || currentScript.value.code == scriptCode.value){
          scriptCode.value = script.value.code ? script.value.code : t('empty');
		}
        lang.value = script.value.lang
        setTimeout(() => {
            resizeHeight('ztf-script-main', 'editor-panel', 'splitter-v', 'logs-panel',
              100, 100, 90)
        }, 600)
    } else {
        scriptCode.value = ''
        lang.value = ''
    }
    currentScript.value = script.value
    store.dispatch('tabs/update', {
        id: props.tab.id,
        title: props.tab.title,
        changed: script.value.code != scriptCode.value,
        type: 'script',
        data: props.tab.data
    });
}, { deep: false })


const editorChange = (newScriptCode) => {
    let oldScriptCode = script.value.code;
    scriptCode.value = newScriptCode;
    let changed = newScriptCode === oldScriptCode ? false : true;
    store.dispatch('tabs/update', {
        id: props.tab.id,
        title: props.tab.title,
        changed: changed,
        type: 'script',
        data: props.tab.data
    });
}

const save = (item) => {
  if(item.path.indexOf(currentScript.value.path) == -1 ){
    return;
  }
  const code = editorRef.value?.getValue()
  const tabId = currentScript.value.workspaceType === 'ztf' && currentScript.value.path.indexOf('.exp') !== currentScript.value.path.length - 4
        ? 'script-' + currentScript.value.path : 'code-' + currentScript.value.path
  store.dispatch('Script/updateCode',{
        workspaceId: global.value[tabId],
        path: currentScript.value.path,
        code: code
    }).then(() => {
        notification.success({
          message: t('save_success'),
        })
        store.dispatch('Script/getScript', {type: 'file', ...currentScript.value})
      })
}

onMounted(() => {
  console.log('onMounted')
  bus.on(settings.eventScriptSave, save);
  document.addEventListener('keydown', function(e){
    const ele = document.getElementsByClassName('vfm--fixed')
    const isFocus = Array.prototype.findIndex.call(ele, function(vftEle){
      return vftEle.style.display != 'none';
    });
    if(isFocus > -1){
        return 
    }
    if (e.keyCode == 83 && (navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey)){
        e.preventDefault();
        save({path: activeID.value})
     }
});
})
onBeforeUnmount( () => {
  console.log('onBeforeUnmount')
  bus.off(settings.eventScriptSave, save);
  document.removeEventListener('keydown', function(e){
    if (e.keyCode == 83 && (navigator.platform.match("Mac") ? e.metaKey : e.ctrlKey)){
        e.preventDefault();
        save(currentScript.value)
     }
});
})

defineExpose({
    save
});
</script>

<style lang="less" scoped>
.ztf-script-main {
    height: 100%;
    display: flex;
    flex-direction: column;

    #editor-panel {
        flex: 1;

        padding: 0;
    }

    .toolbar {
        padding: 4px 10px;
        height: 40px;
        text-align: right;

        .ant-btn {
            margin: 0 5px;
        }
    }
}
</style>

<style lang="less">
.editor-panel {
  height: 10%;
}
</style>
