<template>
  <a-card>
    <template #title>
      {{ modelRef.id > 0 ? t('edit_workspace') : t('create_workspace')}}
    </template>

    <template #extra>
      <div class="opt">
        <a-button @click="back">{{ t('back') }}</a-button>
      </div>
    </template>

    <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-item :label="t('name')" v-bind="validateInfos.name">
        <a-input v-model:value="modelRef.name"
                 @blur="validate('name', { trigger: 'blur' }).catch(() => {})" />
      </a-form-item>

      <a-form-item :label="t('path')" v-bind="validateInfos.path">
        <a-input-search v-if="isElectron" v-model:value="modelRef.path" @search="selectDir" spellcheck="false"
                        @blur="validate('path', { trigger: 'blur' }).catch(() => {})">
          <template #enterButton>
            <a-button>选择</a-button>
          </template>
        </a-input-search>

        <a-input v-if="!isElectron" v-model:value="modelRef.path" spellcheck="false"
                 @blur="validate('path', { trigger: 'blur' }).catch(() => {})"/>
      </a-form-item>

      <a-form-item :label="t('type')" v-bind="validateInfos.type">
        <a-select @change="selectType" v-model:value="modelRef.type">
          <a-select-option v-for="item in testTypes" :key="item.value" :value="item.value">
            {{item.label}}
          </a-select-option>
        </a-select>
      </a-form-item>

      <a-form-item v-if="showCmd" :label="t('test_cmd')" v-bind="validateInfos.cmd">
        <a-textarea v-model:value="modelRef.cmd"
                    :auto-size="{ minRows: 3, maxRows: 6 }" />
        <span>使用客户端执行测试时，-p参数的值会被替换成当前产品的ID。命令行执行时，请提供产品ID的数字。</span>
      </a-form-item>

      <a-form-item :wrapper-col="{ span: 14, offset: 4 }">
        <a-button type="primary" @click.prevent="submitForm">{{t('save')}}</a-button> &nbsp;
        <a-button style="margin-left: 10px" @click="resetFields">{{t('reset')}}</a-button>
      </a-form-item>
    </a-form>

  </a-card>

</template>
<script lang="ts">
import {defineComponent, ref, reactive, computed, Ref,watch, ComputedRef} from "vue";
import {useRouter} from "vue-router";
import { useI18n } from "vue-i18n";

import { validateInfos } from 'ant-design-vue/lib/form/useForm';
import {Form, notification} from 'ant-design-vue';
const useForm = Form.useForm;

import {useStore} from "vuex";
import {WorkspaceData} from "@/views/workspace/store";
import {ZentaoData} from "@/store/zentao";
import {arrToMap} from "@/utils/array";
import {ztfTestTypesDef, unitTestTypesDef} from "@/utils/const";

interface WorkspaceFormSetupData {
  t: (key: string | number) => string;
  isElectron: Ref<boolean>;
  testTypes: Ref<any[]>
  cmdMap: Ref
  showCmd: Ref<boolean>;
  modelRef: ComputedRef;
  labelCol: any
  wrapperCol: any
  validate: any
  validateInfos: validateInfos
  selectType:  (item) => void;
  selectDir: () => void
  submitForm:  () => void;
  resetFields:  () => void;
  back:  () => void;
}

export default defineComponent({
  name: 'WorkspaceForm',
  components: {
  },
  setup(props): WorkspaceFormSetupData {
    const { t } = useI18n();
    const router = useRouter();
    const isElectron = ref(!!window.require)
    const testTypes = ref([...ztfTestTypesDef, ...unitTestTypesDef])
    const cmdMap = ref(arrToMap(testTypes.value))

    const store = useStore<{ Workspace: WorkspaceData }>();
    const modelRef = computed(() => store.state.Workspace.detailResult);
    const showCmd = computed(() => { return modelRef.value.type !== 'ztf' })

    const get = async (id: number): Promise<void> => {
      await store.dispatch('Workspace/get', id);
    }
    const id = +router.currentRoute.value.params.id
    get(id)

    const rules = reactive({
      name: [{ required: true, message: t('pls_name'), trigger: 'blur' }],
      path: [{ required: true, message: t('pls_workspace_path'), trigger: 'blur' }],
      type: [{ required: true, message: t('pls_workspace_type') }],
      cmd: [
        {
          trigger: 'blur',
          validator: async (rule: any, value: string) => {
            if (modelRef.value.type !== 'ztf' && (value === '' || !value)) {
              throw new Error(t('pls_cmd'));
            }
          }
        },
      ],
    });

    const selectType = (item) => {
      console.log('selectType', item)
      modelRef.value.cmd = cmdMap.value[modelRef.value.type].cmd
    }

    const selectDir = () => {
      console.log('selectDir')

      if (isElectron.value) {
        const {dialog} = window.require('@electron/remote');
        dialog.showOpenDialog({
          properties: ['openDirectory']
        }).then(result => {
          if (result.filePaths && result.filePaths.length > 0) {
            modelRef.value.path = result.filePaths[0]
          }
        }).catch(err => {
          console.log(err)
        })
      }
    }

    const { resetFields, validate, validateInfos } = useForm(modelRef, rules);
    const submitForm = () => {
      validate()
        .then(() => {
          console.log(modelRef.value);
          store.dispatch('Workspace/save', modelRef.value).then((success) => {
            if (success) {
              notification.success({message: t('save_success')});
              router.push(`/workspace/list`)
            }
          })
        })
      .catch((e) => {console.log(e)})
    };

    const back = () => {
      console.log('back')
      router.push(`/workspace/list`)
    }

    return {
      t,
      isElectron,
      testTypes,
      cmdMap,
      showCmd,
      labelCol: { span: 4 },
      wrapperCol: { span: 14 },
      modelRef,

      resetFields,
      validate,
      validateInfos,
      selectType,
      selectDir,
      submitForm,
      back,
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
