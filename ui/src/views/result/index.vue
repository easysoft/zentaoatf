<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
        {{ t('exec_result') }}
      </template>
      <template #extra>
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
            <span v-if="record.startTime">{{ momentUtc(record.startTime) }}</span>
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
            <a-button @click="() => view(record)" type="link" size="small">{{ t('view') }}</a-button>
            <a-button @click="() => remove(record)" type="link" size="small"
                      :loading="removeLoading.includes(record.seq)">{{ t('remove') }}
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
import {momentUnixDef, momentUtcDef, percentDef} from "@/utils/datetime";
import {execByDef} from "@/utils/testing";
import {WorkspaceData} from "@/store/workspace";
import {hideMenu} from "@/utils/dom";
import {useI18n} from "vue-i18n";
import IconSvg from "@/components/IconSvg/index";
import {PaginationConfig, QueryParams} from "@/types/data";
import {disableStatusMap} from "@/utils/const";
import debounce from "lodash.debounce";
import {disableStatusDef} from "@/utils/decorator";

const useForm = Form.useForm;

interface ListResultSetupData {
  t: (key: string | number) => string;

  statusArr: Ref,
  queryParams: Ref,
  pagination: ComputedRef<PaginationConfig>;

  columns: any;
  models: ComputedRef<any[]>;
  loading: Ref<boolean>;
  list: (v) => void
  view: (item) => void;
  removeLoading: Ref<string[]>;
  remove: (id) => void;

  onSearch: () => void;
  disableStatus: (val) => string;
  momentUtc: (tm) => string;

  execBy: (item) => string;
  percent: (numb, total) => string;
  simpleImage: any
}

export default defineComponent({
  name: 'ResultListPage',
  components: {
    IconSvg,
  },
  setup(): ListResultSetupData {
    const {t} = useI18n();
    const momentUtc = momentUtcDef
    const disableStatus = disableStatusDef

    const execBy = execByDef
    const percent = percentDef

    const columns = [
      {
        title: t('no'),
        dataIndex: 'seq',
      },
      {
        title: t('exec_type'),
        dataIndex: 'execBy',
        slots: {customRender: 'execBy'},
      },
      {
        title: t('start_time'),
        dataIndex: 'startTime',
        slots: {customRender: 'startTime'},
      },
      {
        title: t('duration'),
        dataIndex: 'duration',
        slots: {customRender: 'duration'},
      },
      {
        title: t('result'),
        dataIndex: 'result',
        slots: {customRender: 'result'},
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
    const store = useStore<{ result: StateType }>();

    const models = computed<any[]>(() => store.state.result.queryResult.result)
    const pagination = computed<PaginationConfig>(() => store.state.result.queryResult.pagination);
    const queryParams = ref<QueryParams>({
      keywords: '', enabled: '1', page: pagination.value.page, pageSize: pagination.value.pageSize
    });

    const loading = ref<boolean>(true);
    const list = (page: number) => {
      loading.value = true;
      store.dispatch('result/list', {
        keywords: queryParams.value.keywords,
        enabled: queryParams.value.enabled,
        pageSize: pagination.value.pageSize,
        page: page});
      loading.value = false;
    }
    list(1);

    const onSearch = debounce(() => {
      list(1)
    }, 500);

    onMounted(() => {
      console.log('onMounted')
    })

    // 查看
    const view = (item) => {
      router.push(`/exec/history/${item.testType}/${item.seq}`)
    }

    // 删除
    const removeLoading = ref<string[]>([]);
    const remove = (item) => {
      Modal.confirm({
        title: t('confirm_to_remove_result'),
        okText: t('confirm'),
        cancelText: t('cancel'),
        onOk: async () => {
          removeLoading.value = [item.seq];
          const res: boolean = await store.dispatch('History/remove', item.seq);
          if (res === true) {
            message.success(t('remove_success'));
            await list(1);
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
      list,

      view,
      removeLoading,
      remove,

      onSearch,
      disableStatus,
      momentUtc,

      execBy,
      percent,
      simpleImage: Empty.PRESENTED_IMAGE_SIMPLE,
    }
  }

})
</script>

<style lang="less" scoped>
.exec-button {
  padding-left: 23px;
  .exec-icon {
    display: inline-block;
    margin-right: 5px;
  }
  .button-text {
    display: inline-block;
    margin-right: 6px;
  }
}

</style>