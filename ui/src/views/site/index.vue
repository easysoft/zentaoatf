<template>
  <div>
    <a-card :bordered="false">
      <template #title>
        <div class="t-card-toolbar">
          <div class="left">
            禅道站点列表
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
        <a-button type="primary" @click="() => edit(0)">新建禅道站点</a-button>
      </template>

      <div>
        <a-table
            row-key="seq"
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
          <template #seq="{ text }">
            {{ text }}
          </template>
          <template #execBy="{ record }">
            {{ execBy(record) }}
          </template>
          <template #startTime="{ record }">
            <span v-if="record.startTime">{{ momentTime(record.startTime) }}</span>
          </template>
          <template #duration="{ record }">
            {{ record.duration }}秒
          </template>
          <template #result="{ record }">
            <span class="t-pass t-status">
              {{ record.pass }}&nbsp;
              <icon-svg type="pass"></icon-svg>&nbsp;
              ({{ percent(record.pass, record.total) }})
            </span>
            <span class="t-fail t-status">
              {{ record.fail }}&nbsp;
              <icon-svg type="fail"></icon-svg>&nbsp;
              ({{ percent(record.fail, record.total) }})
            </span>
            <span class="t-skip t-status">
              {{ record.skip }}&nbsp;
              <icon-svg type="skip"></icon-svg>&nbsp;
              ({{ percent(record.skip, record.total) }})
            </span>
          </template>
          <template #action="{ record }">
            <a-button @click="() => viewResult(record)" type="link" size="small">{{ t('view') }}</a-button>
            <a-button @click="() => deleteExec(record)" type="link" size="small"
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
import {StateType} from "./store";
import {useRouter} from "vue-router";
import {momentTimeDef} from "@/utils/datetime";
import {useI18n} from "vue-i18n";
import IconSvg from "@/components/IconSvg/index";
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

  removeLoading: Ref<string[]>;
  remove: (item) => void;

  onSearch: () => void;
  momentTime: (tm) => string;
  simpleImage: any
}

export default defineComponent({
  name: 'SiteListPage',
  components: {
    IconSvg,
  },
  setup(): SiteListSetupData {
    const {t} = useI18n();
    const momentTime = momentTimeDef

    const onSearch = debounce(() => {
      getList(1)
    }, 500);

    onMounted(() => {
      getList(1);
    })

    const columns = [
      {
        title: t('no'),
        dataIndex: 'seq',
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
        dataIndex: 'createTime',
        slots: {customRender: 'createTime'},
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
    const models = computed<any[]>(() => store.state.Site.queryResult.data);
    const pagination = computed<PaginationConfig>(() => store.state.Site.queryResult.pagination);
    const queryParams = ref<QueryParams>({
      keywords: '', enabled: '1',
      page: pagination.value.current, pageSize: pagination.value.pageSize
    });

    const loading = ref<boolean>(true);
    const getList = (current: number) => {
      loading.value = true;
      store.dispatch('Site/list', {
        keywords: queryParams.value.keywords,
        enabled: queryParams.value.enabled,
        pageSize: pagination.value.pageSize,
        page: current});
      loading.value = false;
    }

    onMounted(() => {
      console.log('onMounted')
      getList(1);
    })

    const removeLoading = ref<string[]>([]);
    const remove = (item) => {
      Modal.confirm({
        title: t('confirm_to_delete_result'),
        okText: t('confirm'),
        cancelText: t('cancel'),
        onOk: async () => {
          removeLoading.value = [item.seq];
          const res: boolean = await store.dispatch('History/delete', item.seq);
          if (res === true) {
            message.success(t('delete_success'));
            await getList(pagination.value.current);
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
      removeLoading,
      remove,

      onSearch,
      momentTime,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }
  }

})
</script>

<style lang="less" scoped>
.ant-card-extra {

}

</style>