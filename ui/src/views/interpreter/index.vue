<template>
  <a-card>
    <template #title>
      <div class="t-card-toolbar">
        <div class="left">
          {{t('interpreter')}}
        </div>
      </div>

    </template>
    <template #extra>
      <a-button v-if="interpreters" @click="create()" type="primary">
        <template #icon><PlusCircleOutlined /></template>
        {{t('create_interpreter')}}
      </a-button>
    </template>

    <div v-if="!interpreters" style="padding: 20px;">
      非Windows平台中，请参照<a-link to="https://ztf.im/book/ztf/ztf-about-26.html">此文</a-link>将可执行文件加入PATH变量中，
      即可在任意目录中执行测试，不需要为各种语言设置运行环境。
    </div>

    <a-table
        v-if="interpreters"
        row-key="id"
        :columns="columns"
        :data-source="interpreters"
        :loading="loading"
        :pagination="false"
    >
      <template #lang="{ record }">
        {{languageMap[record.lang].name}}
      </template>

      <template #createdAt="{ record }">
        <span v-if="record.createdAt">{{ momentUtc(record.createdAt) }}</span>
      </template>

      <template #action="{ record }">
        <a-button @click="() => edit(record)" type="link" size="small">{{ t('edit') }}</a-button>
        <a-button @click="() => remove(record)" type="link" size="small">{{ t('delete') }}
        </a-button>
      </template>

    </a-table>

    <a-modal
        :title="formTitle"
        v-if="formVisible"
        :visible="true"
        :onCancel="onCancel"
        width="800px"
        :destroy-on-close="true"
        :mask-closable="false"
        :footer="null"
    >
      <EditInterpreterForm
          :model="interpreter"
          :onClose="onSave"
      />
    </a-modal>

  </a-card>

</template>
<script lang="ts">
import {defineComponent, ref, reactive, computed, watch, ComputedRef, Ref, toRaw, toRef} from "vue";
import { useI18n } from "vue-i18n";
import { PlusCircleOutlined } from '@ant-design/icons-vue';
import {message, Empty} from 'ant-design-vue';

import EditInterpreterForm from './component/edit.vue';

import {getLangSettings} from "./service";
import {listInterpreter, removeInterpreter} from "@/views/interpreter/service";
import {momentUtcDef} from "@/utils/datetime";
import ALink from "@/components/ALink/index.vue";

export default defineComponent({
  name: 'InterpreterList',
  components: {
    ALink,
    EditInterpreterForm, PlusCircleOutlined,
  },
  setup(props) {
    const {t, locale} = useI18n();
    const momentUtc = momentUtcDef

    let languageMap = ref<any>({})
    const getInterpretersA = async () => {
      const data = await getLangSettings()
      languageMap.value = data.languageMap
    }
    getInterpretersA()

    let interpreters = ref<any>([])
    let interpreter = reactive<any>({})
    const formTitle = computed(() => {
      console.log('interpreter.id', interpreter.id)
      return interpreter.value.id ? t('edit_interpreter') : t('create_interpreter')
    })

    watch(locale, () => {
      console.log('watch locale', locale)
      setColumns()
    }, {deep: true})

    const columns = ref([] as any[])
    const setColumns = () => {
      columns.value = [
        {
          title: t('no'),
          dataIndex: 'index',
          width: 80,
          customRender: ({text, index}: { text: any; index: number }) => index + 1,
        },
        {
          title: t('lang'),
          dataIndex: 'lang',
          slots: {customRender: 'lang'},
        },
        {
          title: t('interpreter_path'),
          dataIndex: 'path',
        },
        {
          title: t('create_time'),
          dataIndex: 'createdAt',
          slots: {customRender: 'createdAt'},
        },
        {
          title: t('opt'),
          key: 'action',
          width: 260,
          slots: {customRender: 'action'},
        },
      ]
    }
    setColumns()

    const loading = ref<boolean>(true);
    const list = () => {
      loading.value = true;

      listInterpreter().then((json => {
        console.log('---', json)

        if (json.code === 0) {
          interpreters.value = json.data
        }
      }))

      loading.value = false;
    }
    list()

    const formVisible = ref(false)
    const setFormVisible = (val: boolean) => {
      formVisible.value = val;
    };
    
    const create = () => {
      interpreter.value = {}
      setFormVisible(true)
    }
    const edit = (item) => {
      interpreter.value = item
      setFormVisible(true)
    }

    const onSave = async () => {
      console.log('onSave')
      setFormVisible(false)
      list()
    }

    const onCancel = async () => {
      console.log('onCancel')
      setFormVisible(false)
    }

    const remove = async (item) => {
      await removeInterpreter(item.id)
      list()
    }

    return {
      t,
      momentUtc,
      languageMap,

      columns,
      loading,
      interpreters,
      interpreter,

      formTitle,
      formVisible,
      setFormVisible,
      create,
      edit,
      remove,
      onSave,
      onCancel,

      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }

  }
})
</script>

<style lang="less" scoped>

.interpreter-header {
  margin: 5px 30px;
  padding-bottom: 6px;
  border-bottom: 1px solid #f0f0f0;
}

.interpreter-item {
  margin: 5px 30px;

}

</style>
