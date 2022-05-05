<template>
  <a-modal
      :title="t('pls_name')"
      :destroy-on-close="true"
      :mask-closable="false"
      :visible="true"
      :onCancel="onCancel"
      :footer="null"
  >
    <div>
      <a-form :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-item :label="t('name')">
          <a-input v-model:value="modelRef.name" />
        </a-form-item>
        <a-form-item></a-form-item>

        <a-form-item :wrapper-col="{ span: wrapperCol.span, offset: labelCol.span }"
                     :class="{'t-dir-right': !isWin}" class="t-right">
          <a-button type="primary" @click="onFinish" class="t-btn-gap">{{t('submit')}}</a-button>
          <a-button @click="() => onCancel()" class="t-btn-gap">{{t('cancel')}}</a-button>
        </a-form-item>

      </a-form>

    </div>

  </a-modal>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, PropType, reactive, ref, Ref} from "vue";
import {message, Form} from 'ant-design-vue';
import {queryProduct, queryTask} from "@/services/zentao";
import {useI18n} from "vue-i18n";
import {isWindows} from "@/utils/comm";
import {useStore} from "vuex";
import {ZentaoData} from "@/store/zentao";
const useForm = Form.useForm;

export default defineComponent({
  name: 'NameForm',
  props: {
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

    const zentaoStore = useStore<{ Zentao: ZentaoData }>();
    const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

    const rules = reactive({
    });

    const modelRef = reactive<any>({taskId: ''})
    let tasks = ref([])

    const { resetFields, validate, validateInfos } = useForm(modelRef, rules);

    const listTask = () => {
      queryTask(currProduct.value.id).then((jsn) => {
        tasks.value = jsn.data
      })
    }
    listTask()

    const onFinish = async () => {
      console.log('onFinish', modelRef)

      validate().then(() => {
        props.onSubmit(modelRef);
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

      modelRef,
      onFinish,
    }
  }
})
</script>