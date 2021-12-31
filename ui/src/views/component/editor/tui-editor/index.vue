<template>
    <div className="indexlayout-main-conent">
        <a-card :bordered="false" title="Editor:">
           <TuiEditor v-model="val" />
        </a-card>
        <a-card :bordered="false" title="Content:" style="margin-top:10px">
           <TuiEditorViewer :value="val" />
        </a-card>
    </div>
</template>
<script lang="ts">
import { defineComponent, ref } from "vue";
import TuiEditor from "@/components/TuiEditor/index.vue";
import TuiEditorViewer from '@/components/TuiEditor/viewer.vue';
export default defineComponent({
    name: 'ComponentEditorCKE',
    components: {
        TuiEditor,
        TuiEditorViewer
    },
    setup() {
        const val = ref<string>('### This is Editor');
        return {
            val
        }
    }
})
</script>
