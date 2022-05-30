<template>
  <div>
    <a-card :bordered="false">
      <template #extra>
        <a-button type="primary" @click="create()">
          <template #icon>
            <PlusCircleOutlined/>
          </template>
          {{ t('create_workspace') }}
        </a-button>
      </template>

      <div>
        <a-table
            row-key="id"
            :columns="columns"
            :data-source="models"
            :loading="loading"
        >
          <template #status="{ record }">
            {{ disableStatus(record.disabled) }}
          </template>
          <template #createdAt="{ record }">
            <span v-if="record.createdAt">{{ momentUtc(record.createdAt) }}</span>
          </template>

          <template #action="{ record }">
            <a-button @click="() => edit(record.id)" type="link" size="small">{{ t('edit') }}</a-button>
            <a-button @click="() => remove(record)" type="link" size="small"
                      :loading="removeLoading.includes(record.seq)">{{ t('delete') }}
            </a-button>
          </template>

        </a-table>
      </div>
    </a-card>

    <a-modal
        :title="t('confirm_to_delete_workspace')"
        v-if="confirmVisible"
        :visible="true"
        :destroy-on-close="true"
    >
      <template #footer>
        <div :class="{'t-dir-right': !isWin}" class="t-right">
          <a-button @click="removeConfirmed()" type="primary" class="t-btn-gap">{{ t('confirm') }}</a-button>
          <a-button @click="confirmVisible = false" class="t-btn-gap">{{ t('cancel') }}</a-button>
        </div>
      </template>
    </a-modal>
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, ref, Ref, watch} from "vue";
import {useStore} from "vuex";

import {Empty, Form, message, Modal} from "ant-design-vue";
import {PlusCircleOutlined} from '@ant-design/icons-vue';

import {WorkspaceData} from "./store";
import {useRouter} from "vue-router";
import {momentUtcDef} from "@/utils/datetime";
import {useI18n} from "vue-i18n";
import {PaginationConfig, QueryParams} from "@/types/data";
import debounce from "lodash.debounce";
import {ZentaoData} from "@/store/zentao";
import {disableStatusDef} from "@/utils/decorator";
import {disableStatusMap} from "@/utils/const";
import {isWindows} from "@/utils/comm";

const useForm = Form.useForm;

export default defineComponent({
  name: 'WorkspaceListPage',
  components: {
    PlusCircleOutlined,
  },
  setup() {
    const {t, locale} = useI18n();
    const isWin = isWindows()
    const momentUtc = momentUtcDef
    const disableStatus = disableStatusDef

    const confirmVisible = ref(false)
    const model = ref({} as any)

    const onSearch = debounce(() => {
      getList(1)
    }, 500);

    onMounted(() => {
      getList(1);
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
        customRender: ({text, index}: { text: any; index: number }) =>
            index + 1,
        },
        {
          title: t('name'),
          dataIndex: 'name',
        },
        {
          title: t('path'),
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

    const statusArr = ref(disableStatusMap);

    const router = useRouter();

    const zentaoStore = useStore<{ Zentao: ZentaoData }>();
    const currProduct = computed<any>(() => zentaoStore.state.Zentao.currProduct);

    const store = useStore<{ WorkspaceOld: WorkspaceData }>();
    const models = computed<any[]>(() => store.state.WorkspaceOld.queryResult?.result);

    const loading = ref<boolean>(true);
    const getList = (page: number) => {
      loading.value = true;
      store.dispatch('WorkspaceOld/query', {});
      loading.value = false;
    }
    getList(1);

    watch(currProduct, () => {
      console.log('watch currProduct', currProduct.value.id)
      getList(1);
    }, {deep: true})

    onMounted(() => {
      console.log('onMounted')
    })

    const create = () => {
      console.log('create')
      router.push(`/workspace/edit/0`)
    }
    const edit = (id) => {
      console.log('edit')
      router.push(`/workspace/edit/${id}`)
    }

    const removeLoading = ref<string[]>([]);
    const remove = (item) => {
      model.value = item
      confirmVisible.value = true
    }
    const removeConfirmed = async () => {
      removeLoading.value = [model.value.id];
      store.dispatch('Workspace/delete', model.value.id).then((success) => {
        if (success) {
          message.success(t('delete_success'));
          getList(1)

          removeLoading.value = [];
          confirmVisible.value = false
        }
      })
    }

    return {
      t,
      isWin,

      statusArr,

      currProduct,
      columns,
      models,
      loading,
      getList,
      create,
      edit,
      removeLoading,
      confirmVisible,
      remove,
      removeConfirmed,

      onSearch,
      disableStatus,
      momentUtc,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }
  }

})
</script>

<style lang="less" scoped>
.ant-card-extra {

}

</style>