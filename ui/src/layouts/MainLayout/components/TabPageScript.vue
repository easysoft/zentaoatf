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
import { computed, ref, watch } from "vue";
import { MonacoOptions, ScriptFileNotExist } from "@/utils/const";
import { resizeHeight, resizeWidth } from "@/utils/dom";
import { useI18n } from "vue-i18n";
import MonacoEditor from "@/components/Editor/MonacoEditor.vue";

const { t } = useI18n();
const props = defineProps<{
    tab: PageTab
}>();

const scriptStore = useStore<{ Script: ScriptData }>();
const script = computed<any>(() => scriptStore.state.Script.detail);
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);
const scriptCode = ref('')
const path = ref('')

const lang = ref('')
const editorOptions = ref(MonacoOptions)
const editorRef = ref<InstanceType<typeof MonacoEditor>>()

const init = ref(false)
watch(script, () => {
    if(script.value == undefined || (path.value !== '' && path.value !== script.value.path)){
        return
    }
    console.log('watch script', script)
    path.value = path.value === '' ? script.value.path : path.value
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
}, { deep: false })

const editorChange = (newScriptCode) => {
    newScriptCode = newScriptCode.replace(/\n$/, '');
    let oldScriptCode = scriptCode.value;
    oldScriptCode = oldScriptCode.replace(/\n$/, '');
    let changed = newScriptCode == oldScriptCode ? false : true;
    scriptStore.dispatch('tabs/update', {
        id: props.tab.id,
        title: props.tab.title,
        changed: changed,
        type: 'script',
        data: props.tab.data
    });
}

const save = () => {
    const code = editorRef.value?.getValue()
    scriptStore.dispatch('Script/updateCode',{
        workspaceId: currWorkspace.value.id,
        path: script.value.path,
        code: code
    }).then(() => {
        console.info("success")
        notification.success({
          message: t('save_success'),
        })
      })
}

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
