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
import { defineProps, defineExpose } from "vue";
import { PageTab } from "@/store/tabs";
import { useStore } from "vuex";
import { ScriptData } from "@/views/script/store";
import { computed, ref, watch, onMounted, onBeforeUnmount } from "vue";
import { MonacoOptions, ScriptFileNotExist } from "@/utils/const";
import { resizeHeight, resizeWidth } from "@/utils/dom";
import { useI18n } from "vue-i18n";
import MonacoEditor from "@/components/MonacoEditor.vue";
import bus from "@/utils/eventBus";
import settings from "@/config/settings";

const { t } = useI18n();
const props = defineProps<{
    tab: PageTab
}>();

const store = useStore<{ Script: ScriptData }>();
const script = computed<any>(() => store.state.Script.detail);
const currWorkspace = computed<any>(() => store.state.Script.currWorkspace);
const scriptCode = ref('')
const currentScript = ref({} as any)

const lang = ref('')
const editorOptions = ref(MonacoOptions)
const editorRef = ref<InstanceType<typeof MonacoEditor>>()

const init = ref(false)
watch(script, () => {
    if(script.value == undefined || (currentScript.value.path!==undefined && currentScript.value.path !== script.value.path)){
        return
    }
    console.log('watch script', script)
    if(!isFirstLoad.value){
        scriptCode.value = String(editorRef.value?.getValue())
        isFirstLoad.value = true
    }
    currentScript.value = script.value
    if (script.value) {
        if (script.value.code === ScriptFileNotExist) {
            scriptCode.value = ScriptFileNotExist
            lang.value = ''

            return
        }
        
        scriptCode.value = script.value.code ? script.value.code : t('empty')
        lang.value = script.value.lang
        setTimeout(() => {
            resizeHeight('ztf-script-main', 'editor-panel', 'splitter-v', 'logs-panel',
              100, 100, 90)
        }, 600)
    } else {
        scriptCode.value = ''
        lang.value = ''
    }
    store.dispatch('tabs/update', {
        id: props.tab.id,
        title: props.tab.title,
        changed: script.value.code != scriptCode.value,
        type: 'script',
        data: props.tab.data
    });
}, { deep: false })

const isFirstLoad = ref(false) // update code from MonacoEditor when first load scriptCode from store.

const editorChange = (newScriptCode) => {
    let oldScriptCode = scriptCode.value;
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
console.log(item.path)
if(item.path.indexOf(currentScript.value.path) == -1 ){
    return;
}
const code = editorRef.value?.getValue()
  store.dispatch('Script/updateCode',{
        workspaceId: currWorkspace.value.id,
        path: currentScript.value.path,
        code: code
    }).then(() => {
        console.info("success")
        store.dispatch('Script/getScript', {type: 'file', ...currentScript.value})
        notification.success({
          message: t('save_success'),
        })
      })
}

onMounted(() => {
  console.log('onMounted')
  bus.on(settings.eventScriptSave, save);
})
onBeforeUnmount( () => {
  console.log('onBeforeUnmount')
  bus.off(settings.eventScriptSave, save);
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
