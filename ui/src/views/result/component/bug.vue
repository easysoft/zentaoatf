<template>
  <a-modal
      :title="t('submit_bug_to_zentao')"
      :destroy-on-close="true"
      :mask-closable="false"
      :visible="true"
      :onCancel="onCancel"
      :footer="null"
      width="800px"
  >
    <div>
      <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-item :label="t('title')" v-bind="validateInfos.title">
          <a-input v-model:value="modelRef.title" />
        </a-form-item>

        <a-form-item :label="t('module')">
          <a-select v-model:value="modelRef.module">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in modules" :key="item.code" :value="item.code+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item :label="t('category')">
          <a-select v-model:value="modelRef.type">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in types" :key="item.code" :value="item.code+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item :label="t('version')">
          <a-select v-model:value="modelRef.openedBuild" mode="multiple">
            <a-select-option v-for="item in builds" :key="item.code" :value="item.code+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item :label="t('severity')">
          <a-select v-model:value="modelRef.severity">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in severities" :key="item.code" :value="item.code+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item :label="t('priority')">
          <a-select v-model:value="modelRef.pri">
            <a-select-option key="" value="">&nbsp;</a-select-option>
            <a-select-option v-for="item in priorities" :key="item.code" :value="item.code+''">{{ item.name }}</a-select-option>
          </a-select>
        </a-form-item>

        <a-form-item :label="t('step')">
          <a-textarea v-model:value="modelRef.steps" :auto-size="{ minRows: 5, maxRows: 8 }" />
        </a-form-item>

        <a-form-item :wrapper-col="{ span: wrapperCol.span, offset: labelCol.span }"
                     :class="{'t-dir-right': !isWin}" class="t-right">
          <a-button type="primary" @click="onFinish" class="t-btn-gap">{{ t('submit') }}</a-button>
          <a-button @click="() => onCancel()" class="t-btn-gap">{{ t('cancel') }}</a-button>
        </a-form-item>

      </a-form>

    </div>

  </a-modal>
</template>

<script lang="ts">
import {defineComponent, onMounted, PropType, reactive, ref, Ref} from "vue";
import {Form} from 'ant-design-vue';
import { prepareBugData } from "@/services/bug";
import { queryBugFields } from "@/services/zentao";
import {useI18n} from "vue-i18n";
import {isWindows} from "@/utils/comm";
const useForm = Form.useForm;

export default defineComponent({
  name: 'BugForm',
  props: {
    model: {
      type: Object as PropType<any>,
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

  setup(props) {
    const { t } = useI18n();
    const isWin = isWindows()

    const rules = reactive({
      title: [
        { required: true, message: t('pls_title') },
      ],
      product: [
        { required: true, message: t('pls_product') },
      ],
    });

    const modelRef = ref<any>({})

    let modules = ref([])
    let types = ref([])
    let builds = ref([])
    let severities = ref([])
    let priorities = ref([])

    const { resetFields, validate, validateInfos } = useForm(modelRef, rules);

    const getBugData = () => {
      prepareBugData(props.model).then((jsn) => {
        modelRef.value = jsn.data
        modelRef.value.module = ''
        modelRef.value.severity = ''+modelRef.value.severity
        modelRef.value.pri = ''+modelRef.value.pri

        getBugFields()
      })
    }
    getBugData()

    const getBugFields = () => {
      queryBugFields().then((jsn) => {
        modules.value = jsn.data.modules
        types.value = jsn.data.type
        builds.value = jsn.data.build
        severities.value = jsn.data.severity
        priorities.value = jsn.data.pri
      })
    }

    const onFinish = async () => {
      console.log('onFinish', modelRef)

      validate().then(() => {
        props.onSubmit(modelRef.value);
      }).catch(err => { console.log('') })
    }

    onMounted(()=> {
      console.log('onMounted')
    })

    return {
      t,
      isWin,
      labelCol: { span: 6 },
      wrapperCol: { span: 16 },
      rules,
      validate,
      validateInfos,
      resetFields,

      modules,
      types,
      builds,
      severities,
      priorities,

      modelRef,
      onFinish,
    }
  }
})
</script>