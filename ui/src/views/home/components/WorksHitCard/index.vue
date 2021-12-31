<template>
    <a-card
      class="homeBoxCard"
      :title="t('page.home.workshitcard.card-title')"
    >
      <a-table
        size="small"
        rowKey="id"
        :columns="columns"
        :dataSource="list"
        :loading="loading"
        :pagination="pagination"
        @change="(p) => {
          getList(p.current || 1);
        }"
      />
    </a-card>
</template>
<script lang="ts">
import { computed, ComputedRef, defineComponent, onMounted, Ref, ref } from "vue";
import { useStore } from "vuex";
import { useI18n } from "vue-i18n";
import { StateType as HomeStateType } from "../../store";
import { PaginationConfig } from "../../data";
import { TableListItem } from "./data";

interface WorksHitCardSetupData {
    t: (key: string | number) => string;
    columns: any;
    list: ComputedRef<TableListItem[]>;
    pagination: ComputedRef<PaginationConfig>;
    loading: Ref<boolean>;
    getList: (current: number) => Promise<void>;
}

export default defineComponent({
    name: 'WorksHitCard',
    setup(): WorksHitCardSetupData {
        const store = useStore<{ Home: HomeStateType}>();
        const { t } = useI18n();

        // 分页
        const pagination = computed<PaginationConfig>(() => store.state.Home.worksHitData.pagination);

        // 数据
        const list = computed<TableListItem[]>(()=> store.state.Home.worksHitData.list);

        // 列
        const columns = [
            {
                dataIndex: 'index',
                title: t('page.home.workshitcard.card.table-column-number'),
                width: 80,
                customRender: ({ index }: {index: number}) => {
                    return (pagination.value.current - 1) * pagination.value.pageSize + index + 1;
                },
            },
            {
                dataIndex: 'title',
                title: t('page.home.workshitcard.card.table-column-title'),
            },
            {
                dataIndex: 'hit',
                title: t('page.home.workshitcard.card.table-column-hit'),
            },
        ];

        // 读取数据 func
        const loading = ref<boolean>(true);
        const getList = async (current: number): Promise<void> => {
            loading.value = true;
            await store.dispatch('Home/queryWorksHitData', {
                per: pagination.value.pageSize,
                page: current,
            });
            loading.value = false;
        }

        onMounted(()=> {
           getList(1);
        })


        return {
            t,
            columns,
            list,
            pagination,
            loading,
            getList
        }
    }
})
</script>
<style lang="less" scoped>
.homeBoxCard {
  margin-bottom: 24px;
  ::v-deep(.ant-card-head) {
    padding-left: 12px;
    padding-right: 12px;
  }
  ::v-deep(.ant-card-body) {
    padding: 12px;
  }
  ::v-deep(.ant-divider-horizontal) {
    margin: 8px 0;
  }
}
</style>