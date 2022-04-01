<template>
  <div class="site-main">
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
                            :placeholder="t('input_keyword_to_search')" style="width:270px;margin-left: 16px;"/>
          </div>
        </div>

      </template>
      <template #extra>
        <a-button type="primary" @click="create()">
          <template #icon><PlusCircleOutlined /></template>
          {{t('create_site')}}
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
                    list(page);
                },
                onShowSizeChange: (page, size) => {
                    pagination.pageSize = size
                    list(page);
                },
            }"
            :locale="localeConf"
        >
          <template #status="{ record }">
            {{ disableStatus(record.disabled) }}
          </template>
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

    <a-modal
        :title="t('confirm_to_delete_site')"
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
import {getInitStatus, setInitStatus} from "@/utils/cache";
import {isWindows} from "@/utils/comm";

const useForm = Form.useForm;

export default defineComponent({
  name: 'SiteListPage',
  components: {
    PlusCircleOutlined,
  },
  setup() {
    const {t, locale} = useI18n();
    const isWin = isWindows()
    const momentUtc = momentUtcDef
    const disableStatus = disableStatusDef

    const localeConf = {} as any
    getInitStatus().then((initStatus) => {
      console.log('initStatus', initStatus)
      if (!initStatus) {
        localeConf.emptyText = t('pls_add_zentao_site')

        setTimeout(() => {
          setInitStatus()
        }, 1000)
      }
    })

    onMounted(() => {
      console.log('onMounted')
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
      ]
    }
    setColumns()

    const statusArr = ref(disableStatusMap);

    const router = useRouter();
    const zentaoStore = useStore<{ zentao: ZentaoData }>();
    const store = useStore<{ Site: StateType }>();
    const models = computed<any[]>(() => store.state.Site.queryResult.result);
    const pagination = computed<PaginationConfig>(() => store.state.Site.queryResult.pagination);
    const queryParams = ref<QueryParams>({
      keywords: '', enabled: '1', page: pagination.value.page, pageSize: pagination.value.pageSize
    });

    const confirmVisible = ref(false)
    const model = ref({} as any)

    const loading = ref<boolean>(true);
    const list = (page: number) => {
      loading.value = true;
      store.dispatch('Site/list', {
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

    const create = () => {
      console.log('create')
      router.push(`/site/edit/0`)
    }
    const edit = (id) => {
      console.log('edit')
      router.push(`/site/edit/${id}`)
    }

    const removeLoading = ref<string[]>([]);
    const remove = (item) => {
      model.value = item
      confirmVisible.value = true
    }
    const removeConfirmed = async () => {
      removeLoading.value = [model.value.id];
      store.dispatch('Site/delete', model.value.id).then((success) => {
        zentaoStore.dispatch('zentao/fetchSitesAndProduct').then((success) => {
          message.success(t('delete_success'));
          list(pagination.value.page)

          removeLoading.value = [];
          confirmVisible.value = false
        })
      })
    }

    return {
      t,
      isWin,

      statusArr,
      queryParams,
      pagination,

      columns,
      models,
      loading,
      list,
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
      localeConf,
    }
  }

})
</script>

<style lang="less">
.site-main {
  .ant-table-placeholder {
    color: #000c17;
  }
}

</style>

<style lang="less" scoped>
</style>