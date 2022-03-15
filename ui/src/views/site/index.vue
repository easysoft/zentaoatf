<template>
  <div>
    <a-card :bordered="false">
      <template #title>
        <div class="t-card-toolbar">
          <div class="left">
            {{t('zentao_site')}}
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
        <a-button type="primary" @click="() => edit(0)">
          <template #icon><PlusCircleOutlined /></template>
          新建站点
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

          <template #createdAt="{ record }">
            <span v-if="record.createdAt">{{ momentUtc(record.createdAt) }}</span>
          </template>

          <template #action="{ record }">
            <a-button @click="() => edit(record.id)" type="link" size="small">{{ t('edit') }}</a-button>
            <a-button @click="() => remove(record.id)" type="link" size="small"
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

import {Empty, Form, message, Modal} from "ant-design-vue";
import { PlusCircleOutlined } from '@ant-design/icons-vue';

import {StateType} from "./store";
import {useRouter} from "vue-router";
import {momentUtcDef} from "@/utils/datetime";
import {useI18n} from "vue-i18n";
import {PaginationConfig, QueryParams} from "@/types/data";
import debounce from "lodash.debounce";

const useForm = Form.useForm;

interface SiteListSetupData {
  t: (key: string | number) => string;

  statusArr: Ref,
  queryParams: Ref,
  pagination: ComputedRef<PaginationConfig>;

  columns: any;
  models: ComputedRef<any[]>;
  loading: Ref<boolean>;
  getList: (curr) => void

  edit: (id) => void;
  removeLoading: Ref<string[]>;
  remove: (id) => void;

  onSearch: () => void;
  momentUtc: (tm) => string;
  simpleImage: any
}

export default defineComponent({
  name: 'SiteListPage',
  components: {
    PlusCircleOutlined,
  },
  setup(): SiteListSetupData {
    const {t} = useI18n();
    const momentUtc = momentUtcDef

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
        title: t('zentao_url'),
        dataIndex: 'url',
      },
      {
        title: t('username'),
        dataIndex: 'username',
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
    const statusArr = ref([
      {
        label: '所有状态',
        value: '',
      },
      {
        label: '启用',
        value: '1',
      },
      {
        label: '禁用',
        value: '0',
      },
    ]);

    const router = useRouter();
    const store = useStore<{ Site: StateType }>();
    const models = computed<any[]>(() => store.state.Site.queryResult.result);
    const pagination = computed<PaginationConfig>(() => store.state.Site.queryResult.pagination);
    const queryParams = ref<QueryParams>({
      keywords: '', enabled: '1', page: pagination.value.page, pageSize: pagination.value.pageSize
    });

    const loading = ref<boolean>(true);
    const getList = (page: number) => {
      loading.value = true;
      store.dispatch('Site/list', {
        keywords: queryParams.value.keywords,
        enabled: queryParams.value.enabled,
        pageSize: pagination.value.pageSize,
        page: page});
      loading.value = false;
    }
    getList(1);

    onMounted(() => {
      console.log('onMounted')
    })

    const edit = (id) => {
      console.log('edit')
      router.push(`/site/edit/${id}`)
    }

    const removeLoading = ref<string[]>([]);
    const remove = (id) => {
      Modal.confirm({
        title: t('confirm_to_delete_site'),
        okText: t('confirm'),
        cancelText: t('cancel'),
        onOk: async () => {
          removeLoading.value = [id];
          const res: boolean = await store.dispatch('Site/delete', id);
          if (res === true) {
            message.success(t('delete_success'));
            await getList(pagination.value.page);
          }
          removeLoading.value = [];
        }
      });
    }

    return {
      t,

      statusArr,
      queryParams,
      pagination,

      columns,
      models,
      loading,
      getList,
      edit,
      removeLoading,
      remove,

      onSearch,
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