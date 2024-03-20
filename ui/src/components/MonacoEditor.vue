<template>
  <div class="monaco-editor-vue3" :style="style" ref="elRef"></div>
</template>

<script setup lang="ts">
import { defineProps, onBeforeUnmount, computed, CSSProperties, watch, shallowRef, defineEmits, onMounted, defineExpose } from 'vue';
import * as monaco from 'monaco-editor';
import { useElementSize } from '@vueuse/core'

const props = defineProps({
    width: {type: [String, Number], default: '100%'},
    height: {type: [String, Number], default: '100%'},
    value: String,
    language: {type: String, default: 'javascript'},
    theme: {type: String, default: 'vs'},
    options: {type: Object, default() {return {};}},
});

const emit = defineEmits([
    'editorWillMount',
    'editorDidMount',
    'change'
]);

const elRef = shallowRef<HTMLElement>();
const editorRef = shallowRef<monaco.editor.IStandaloneCodeEditor>();

const style = computed(()=>{
    const {width, height} = props;
    const fixedWidth = (typeof width === 'string' && width.includes('%')) ? width : `${width}px`
    const fixedHeight = (typeof height === 'string' && height.includes('%')) ? height : `${height}px`
    return {
        width: fixedWidth,
        height: fixedHeight,
        textAlign: 'left'
    } as CSSProperties;
});

const {height: containerHeight} = useElementSize(elRef);

watch(containerHeight, (_newVal, _oldValue) => {
    if (editorRef.value) {
       editorRef.value.layout();
    }
});

watch(() => props.options, () => {
    if (editorRef.value) {
        editorRef.value.updateOptions(props.options);
    }
});

watch(() => props.value, () => {
    if (editorRef.value && props.value && editorRef.value.getValue() !== props.value) {
        editorRef.value.setValue(props.value);
    }
});

watch(() => props.language, () => {
    if (editorRef.value) {
        const model = editorRef.value.getModel();
        if (model) {
            monaco.editor.setModelLanguage(model, props.language);
        }
    }
});

watch(() => props.theme, () => {
    monaco.editor.setTheme(props.theme);
});

onMounted(() => {
    const editorOptions = {
        value: props.value,
        theme: props.theme,
        automaticLayout: true,
        ...props.options,
        scrollbar: {
            useShadows: false,
            verticalScrollbarSize: 6,
            horizontalScrollbarSize: 6
        },
        language: props.language
    };
    emit('editorWillMount', editorOptions);

    if (elRef.value) {
        const editor = monaco.editor.create(elRef.value, editorOptions);
        editorRef.value = editor;

        editor.onDidChangeModelContent(event => {
            const editorValue = editor.getValue()
            emit('change', editorValue, event);
        });

        emit('editorDidMount', editor);
        setTimeout(() => {
            editor.getAction('editor.action.formatDocument')?.run();
        }, 100);
    }
});

onBeforeUnmount(() => {
    if (editorRef.value) {
       editorRef.value.dispose();
    }
});

function getValue() {
    if (editorRef.value) {
        return editorRef.value.getValue();
    }
}

defineExpose({
    getValue
});
</script>
