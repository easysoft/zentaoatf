<template>
  <a-modal
      :destroy-on-close="true"
      :mask-closable="false"
      :title="t('create_project')"
      :visible="visible"
      :onCancel="onCancel"
  >
    <template #footer>
      <a-button key="submit" type="primary" @click="onFinish">{{t('save')}}</a-button>
      <a-button key="back" @click="() => onCancel()">{{t('cancel')}}</a-button>
    </template>

    <div>
      <a-form :labelCol="{ span: 4 }" :wrapper-col="{span:20}">
        <a-form-item :label="t('path')" v-bind="validateInfos.path" :placeholder="t('project_path')">
          <a-input-search
              v-if="isElectron"
              v-model:value="modelRef.path"
              @search="selectDir"
              spellcheck="false"
          >
            <template #enterButton>
              <a-button>选择</a-button>
            </template>
          </a-input-search>

          <a-input
              v-if="!isElectron"
              v-model:value="modelRef.path"
              spellcheck="false" />
        </a-form-item>

        <a-form-item :label="t('name')">
          <a-input v-model:value="modelRef.name" :placeholder="t('use_dir_name')"/>
        </a-form-item>

        <a-form-item :label="t('type')" v-bind="validateInfos.type">
          <a-select v-model:value="modelRef.type">
            <a-select-option key="func" value="func">{{t('test_type_ztf')}}</a-select-option>
            <a-select-option key="unit" value="unit">{{t('test_type_other')}}</a-select-option>
          </a-select>
        </a-form-item>

      </a-form>

    </div>

  </a-modal>
</template>

<script lang="ts">
import {defineComponent, onMounted, PropType, reactive, ref, Ref} from "vue";
import {Interpreter} from "@/views/config/data";
import {useI18n} from "vue-i18n";
import {Form} from "ant-design-vue";
import { validateInfos } from 'ant-design-vue/lib/form/useForm';

interface ProjectCreateFormSetupData {
  t: (key: string | number) => string;
  isElectron: Ref<boolean>;
  selectDir: () => void;
  modelRef: Ref<Interpreter>
  validateInfos: validateInfos
  onFinish: () => Promise<void>;
}

export default defineComponent({
  name: 'ProjectCreateForm',
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    onCancel: {
      type: Function,
      required: true
    },
    onSubmit: {
      type: Function as PropType<(values: any) => void>,
      required: true
    }
  },

  components: {},

  setup(props): ProjectCreateFormSetupData {
    const { t } = useI18n();

    const modelRef = ref<any>({path: '', type: 'func'})
    const rulesRef = reactive({
      path: [ { required: true, message: t('pls_project_path') } ],
      type: [ { required: true, message: t('pls_project_type') } ],
    });

    const { validate, validateInfos } = Form.useForm(modelRef, rulesRef);

    const isElectron = ref(!!window.require)
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

    const onFinish = async () => {
      console.log('onFinish')

      validate().then(() => {
        props.onSubmit(modelRef);
      }).catch(err => { console.log('') })
    }

    onMounted(()=> {
      console.log('onMounted')
    })

    return {
      t,
      isElectron,
      selectDir,
      modelRef,
      validateInfos,
      onFinish
    }
  }
})
</script>