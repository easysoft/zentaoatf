<template>
    <div class="ztf-script-main">
        <div id="right-content">
            <div id="editor-panel" class="tab-page-script">
                <MonacoEditor v-if="scriptCode !== '' && scriptCode !== ScriptFileNotExist" v-model:value="scriptCode"
                    :language="lang" :options="editorOptions" class="editor" ref="editorRef" @change="editorChange" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { defineProps } from "vue";
import { PageTab } from "@/store/tabs";
import { useStore } from "vuex";
import { ScriptData } from "@/views/script/store";
import { computed, defineComponent, onMounted, onUnmounted, ref, watch } from "vue";
import { MonacoOptions, ScriptFileNotExist } from "@/utils/const";
import { resizeHeight, resizeWidth } from "@/utils/dom";
import { useI18n } from "vue-i18n";
import MonacoEditor from "@/components/Editor/MonacoEditor.vue";

const { t } = useI18n();
const props = defineProps<{
    tab: PageTab
}>();

let tabMap = ref({})
const scriptStore = useStore<{ Script: ScriptData }>();
let script = computed<any>(() => scriptStore.state.Script.detail);
let scriptCode = ref('')

let lang = ref('')
const editorOptions = ref(MonacoOptions)

watch(script, () => {
    console.log('watch script', script)

    if (script.value) {
        if (script.value.code === ScriptFileNotExist) {
            scriptCode.value = ScriptFileNotExist
            lang.value = ''

            return
        }

        scriptCode.value = script.value.code ? script.value.code : t('empty')
        lang.value = script.value.lang
        setTimeout(() => {
            resizeHeight('right-content', 'editor-panel', 'splitter-v', 'logs-panel',
                100, 100, 90)
        }, 600)
    } else {
        scriptCode.value = ''
        lang.value = ''
    }
}, { deep: true })

const editorChange = (newScriptCode) => {
    newScriptCode = newScriptCode.replace(/\n$/, '');
    let changed = newScriptCode == scriptCode.value ? false : true;
    scriptStore.dispatch('tabs/update', {
        id: props.tab.id,
        title: props.tab.title,
        changed: changed,
        type: 'script',
        data: props.tab.data
    });
}
</script>

<style lang="less">
#editor-panel {
    .script_file_not_exist {
        padding: 10px;
    }
}
</style>

<style lang="less" scoped>
.ztf-script-main {
    flex: 1;
    height: 100%;

    .toolbar {
        padding: 4px 10px;
        height: 40px;
        text-align: right;

        .ant-btn {
            margin: 0 5px;
        }
    }

    #right-content {
        height: calc(100% - 40px);

        display: flex;
        flex-direction: column;

        #editor-panel {
            flex: 1;

            padding: 0 6px 0 8px;
            overflow: auto;
        }

        #splitter-v {
            width: 100%;
            height: 2px;
            background-color: #D0D7DE;
            cursor: ns-resize;

            &.active {
                background-color: #a9aeb4;
            }
        }

        #logs-panel {
            height: 220px;
        }

        .logs-panel {
            height: 100%;
        }
    }
}
</style>

<style lang="less">
.monaco-editor {
    padding: 10px 0;
}
</style>