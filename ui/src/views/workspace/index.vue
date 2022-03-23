<template>
  <div>
    <a-card :bordered="false">
      <template #title>
        <div class="t-card-toolbar">
          <div class="left">
            {{t('workspace')}}
          </div>
          <div class="right">
            <a-select @change="onSearch" v-model:value="queryParams.enabled" :options="statusArr" class="status-select">
            </a-select>
            <a-input-search @change="onSearch" @search="onSearch" v-model:value="queryParams.keywords"
                            placeholder="输入关键字搜索" style="width:270px;margin-left: 16px;"/>
          </div>
        </div>
      </template>

      <template #extra>
        <a-button v-if="currProduct?.id" type="primary" @click="create()">
          <template #icon><PlusCircleOutlined /></template>
          {{t('create_workspace')}}
        </a-button>
      </template>

      <div>
        <a-table
            row-key="id"
            :columns="columns"
            :data-source="models"
            :loading="loading"
            :pagination="{
                ...pagination,
                onChange: (page) => {
                    getList(page);
                },
                onShowSizeChange: (page, size) => {
                    pagination.pageSize = size
                    getList(page);
                },
            }"
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
  </div>
</template>

<script lang="ts">
import {computed, ComputedRef, defineComponent, onMounted, ref, Ref, watch} from "vue";
import {useStore} from "vuex";

import {Empty, Form, message, Modal, notification} from "ant-design-vue";
import { PlusCircleOutlined } from '@ant-design/icons-vue';

import {StateType} from "./store";
import {useRouter} from "vue-router";
import {momentUtcDef} from "@/utils/datetime";
import {useI18n} from "vue-i18n";
import {PaginationConfig, QueryParams} from "@/types/data";
import debounce from "lodash.debounce";
import {ZentaoData} from "@/store/zentao";
import {disableStatusDef} from "@/utils/decorator";
import {disableStatusMap} from "@/utils/const";

const useForm = Form.useForm;

interface WorkspaceListSetupData {
  t: (key: string | number) => string;

  statusArr: Ref,
  queryParams: Ref,
  pagination: ComputedRef<PaginationConfig>;

  columns: any;
  currProduct: ComputedRef
  models: ComputedRef<any[]>;
  loading: Ref<boolean>;
  getList: (curr) => void

  create: () => void;
  edit: (id) => void;
  removeLoading: Ref<string[]>;
  remove: (item) => void;

  onSearch: () => void;
  disableStatus: (val) => string;
  momentUtc: (tm) => string;
  simpleImage: any
}

export default defineComponent({
  name: 'WorkspaceListPage',
  components: {
    PlusCircleOutlined,
  },
  setup(): WorkspaceListSetupData {
    const {t} = useI18n();
    const momentUtc = momentUtcDef
    const disableStatus = disableStatusDef

    const onSearch = debounce(() => {
      getList(1)
    }, 500);

    onMounted(() => {
      getList(1);
    })

    const columns = [
      {
        title: t('no'),
        dataIndex: 'index',
        width: 80,
        customRender: ({text, index}: { text: any; index: number }) =>
            (pagination.value.page - 1) * pagination.value.pageSize + index + 1,
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
        title: t('status'),
        dataIndex: 'status',
        slots: {customRender: 'status'},
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
    ];
    const statusArr = ref(disableStatusMap);

    const router = useRouter();

    const zentaoStore = useStore<{ zentao: ZentaoData }>();
    const currProduct = computed<any>(() => zentaoStore.state.zentao.currProduct);

    const store = useStore<{ Workspace: StateType }>();
    const models = computed<any[]>(() => store.state.Workspace.queryResult.result);
    const pagination = computed<PaginationConfig>(() => store.state.Workspace.queryResult.pagination);
    const queryParams = ref<QueryParams>({
      keywords: '', enabled: '1', page: pagination.value.page, pageSize: pagination.value.pageSize
    });

    const loading = ref<boolean>(true);
    const getList = (page: number) => {
      loading.value = true;
      store.dispatch('Workspace/list', {
        keywords: queryParams.value.keywords,
        enabled: queryParams.value.enabled,
        pageSize: pagination.value.pageSize,
        page: page});
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
      console.log(item)
      Modal.confirm({
        title: t('confirm_to_delete_workspace', {p: item.name}),
        okText: t('confirm'),
        cancelText: t('cancel'),
        onOk: () => {
          removeLoading.value = [item.id];
          store.dispatch('Workspace/delete', item.id).then((success) => {
            if (success) {
              message.success(t('delete_success'));
              getList(pagination.value.page)
            }
            removeLoading.value = [];
          })
        }
      });
    }

    return {
      t,

      statusArr,
      queryParams,
      pagination,

      currProduct,
      columns,
      models,
      loading,
      getList,
      create,
      edit,
      removeLoading,
      remove,

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