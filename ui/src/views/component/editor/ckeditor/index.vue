<template>
    <div className="indexlayout-main-conent">
        <a-card :bordered="false" title="Editor:">
           <CKEditor v-model="val" />
        </a-card>
        <a-card :bordered="false" title="Content:" style="margin-top:10px">
           <div v-html="val"></div>
        </a-card>
    </div>
</template>
<script lang="ts">
import { defineComponent, ref } from "vue";
import CKEditor from "@/components/CKEditor/index.vue";
export default defineComponent({
    name: 'ComponentEditorCKE',
    components: {
        CKEditor
    },
    setup() {
        const val = ref<string>('<h1>This is Editor</h1>');
        return {
            val
        }
    }
})
</script>
