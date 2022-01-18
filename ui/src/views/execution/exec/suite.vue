<template>
    <div class="indexlayout-main-conent">
        <a-card :bordered="false">

        </a-card>
    </div>
</template>

<script lang="ts">
import {ComputedRef, defineComponent, ref, Ref, reactive, computed, onMounted} from "vue";
import { SelectTypes } from 'ant-design-vue/es/select';
import {Execution} from '../data.d';
import {useStore} from "vuex";

import { Props } from 'ant-design-vue/lib/form/useForm';
import { message, Modal, Form } from "ant-design-vue";
const useForm = Form.useForm;

import CreateForm from './components/CreateForm.vue';
import UpdateForm from './components/UpdateForm.vue';

import {StateType as ListStateType} from "../store";
import debounce from "lodash.debounce";
import {useRoute, useRouter} from "vue-router";

interface ExecCasePageSetupData {
  model: any
}

export default defineComponent({
    name: 'ExecutionListPage',
    components: {
    },
    setup(): ExecCasePageSetupData {
      const model = {}

      return {
        model
      }
    }

})
</script>

<style lang="less" scoped>

</style>