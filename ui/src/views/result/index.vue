<template>
  <div class="indexlayout-main-conent">
    <a-card :bordered="false">
      <template #title>
        {{ t('test_result') }}
      </template>
      <template #extra>
      </template>

      <div>
        <a-table
            row-key="no"
            :columns="columns"
            :data-source="models"
            :loading="loading"
            :pagination="{
                ...pagination,
                onChange: (page) => {
                    list(page);
                },
                onShowSizeChange: (page, size) => {
                    pagination.pageSize = size
                    list(page);
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
import {momentUnixDef, percentDef} from "@/utils/datetime";
import {execByDef} from "@/utils/testing";
import {useI18n} from "vue-i18n";
import IconSvg from "@/components/IconSvg/index";
import {PaginationConfig, QueryParams} from "@/types/data";
import {disableStatusMap} from "@/utils/const";
import debounce from "lodash.debounce";
import {disableStatusDef} from "@/utils/decorator";

const useForm = Form.useForm;

export default defineComponent({
  name: 'ResultListPage',
  components: {
    IconSvg,
  },
  setup() {
    const {t} = useI18n();
    const momentTime = momentUnixDef
    const disableStatus = disableStatusDef

    const execBy = execByDef
    const percent = percentDef

    const columns = [
      {
        title: '工作目录',
        dataIndex: 'workspaceName',
      },
      {
        title: '序号',
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
      router.push(`/result/${item.testType}/${item.workspaceId}/${item.seq}`)
    }

    // 删除
    const removeLoading = ref<string[]>([]);
    const remove = (item) => {
      Modal.confirm({
        title: t('confirm_to_delete_result'),
        okText: t('confirm'),
        cancelText: t('cancel'),
        onOk: async () => {
          removeLoading.value = [item.seq];
          const res: boolean = await store.dispatch('result/delete',
              {workspaceId: item.workspaceId, seq: item.seq});
          if (res === true) {
            message.success(t('delete_success'));
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
      momentTime,

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