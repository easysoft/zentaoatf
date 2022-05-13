<template>
  <div class="workdir">
    <Tree :data="treeData" :checkable="checkable" ref="treeRef" />

    <Form labelCol="50px" wrapperCol="60">
      <FormItem name="name" :label="t('title')" :info="validateInfos.name">
        <input v-model="modelRef.name" />
      </FormItem>
      <FormItem name="email" :label="t('email')" :info="validateInfos.email">
        <input v-model="modelRef.email" />
      </FormItem>
      <FormItem name="num" :label="t('number')" :info="validateInfos.num">
        <input v-model="modelRef.num" />
      </FormItem>

      <FormItem size="small">
        <button @click="submit" type="button">{{ t('submit') }}</button>
        <button @click="reset" type="button">{{ t('reset') }}</button>
      </FormItem>
    </Form>

    <ScriptTreePage></ScriptTreePage>
  </div>
</template>

<script setup lang="ts">
import ScriptTreePage from "../../../views/script/component/tree.vue";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
import {computed, onMounted, defineExpose, ref, reactive} from "vue";
import {ScriptData} from "@/views/script/store";
import {resizeWidth} from "@/utils/dom";

import Row from "./Row.vue";
import Col from "./Col.vue";
import Form from "./Form.vue";
import FormItem from "./FormItem.vue";
import {useForm} from "@/utils/form";
import Tree from "./Tree.vue";

const { t } = useI18n();

const zentaoStore = useStore<{ Zentao: ZentaoData }>();
const currSite = computed<any>(() => zentaoStore.state.Zentao.currSite);
const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

const scriptStore = useStore<{ Script: ScriptData }>();
const currWorkspace = computed<any>(() => scriptStore.state.Script.currWorkspace);

onMounted(() => {
  console.log('onMounted')
  setTimeout(() => {
    resizeWidth('main', 'left', 'splitter-h', 'right', 380, 800)
  }, 600)
})

const modelRef = ref({})
const rulesRef = ref({
  name: [
    {required: true, msg: 'Please input name.'},
  ],
  email: [
    {required: true, msg: 'Please input email.'},
    {email: true, msg: 'Please check email format.'},
  ],
  num: [
    {required: true, msg: 'Please input num.'},
    {regex: '^[0-9]*$', msg: 'Please input a number.'},
  ],
})

const { validate, reset, validateInfos } = useForm(modelRef, rulesRef);
const submit = () => {
  if (validate()) {
    console.log('TODO')
  }
}

const treeRef = ref<{isAllCollapsed: () => boolean, toggleAllCollapsed: () => void}>();

const treeData = reactive([{
    title: 'ztf_lang_test',
    collapsed: false,
    checkable: false,
    id: 'root',
    toolbarItems: [
        {hint: 'Add sub folder', icon: 'folder-add'},
        {hint: 'Add file', icon: 'file-add'},
    ],
    children: [
        {
            id: '1',
            title: 'bat',
            children: [
                {
                    id: 'test',
                    title: 'test.bat'
                },
                {
                    id: 'test1',
                    title: 'test_fast.bat'
                }
            ]
        }, {
            id: '5',
            title: 'javascript',
            children: []
        }, {
            id: '6',
            title: 'file',
        }
    ]
}, {
    title: 'demo',
    collapsed: true,
    checkable: false,
    id: 'demo',
    children: [
        {
            id: 'demo1',
            title: 'demo1.txt'
        }
    ]
}, {
    title: 'selenium',
    collapsed: true,
    checkable: false,
    id: 'selenium',
    children: [
        {
            id: 'selenium1',
            title: 'vendor',
            children: [
                {
                    id: 'seleniumtest',
                    title: 'chrome.php'
                },
                {
                    id: 'seleniumtest1',
                    title: 'firefox.js'
                }
            ]
        }, {
            id: 'browser',
            title: 'javascript',
            children: []
        }
    ]
}]);
const checkable = ref(false);

function toggleCheckable(toggle?: boolean) {
    if (toggle === undefined) {
        toggle = !checkable.value;
    }
    checkable.value = toggle;
}

defineExpose({
    get isCheckable() {
        return checkable.value;
    },
    get isAllCollapsed() {
        return treeRef.value?.isAllCollapsed();
    },
    toggleAllCollapsed() {
        return treeRef.value?.toggleAllCollapsed();
    },
    toggleCheckable,
});
</script>

<style lang="less" scoped>
.workdir {

}
</style>
