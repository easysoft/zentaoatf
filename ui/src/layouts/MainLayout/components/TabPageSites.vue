<template>
  <div class="tab-page-sites">
    <div>{{ t('site_management') }}</div>
    <br/>
    <Table
        :is-loading="table.isLoading"
        :columns="table.columns"
        :rows="table.rows"
        :total="table.totalRecordCount"
        :sortable="table.sortable"
        :messages="table.messages"
        @do-search="doSearch"
        @is-finished="table.isLoading = false"
        :pageSize="3"
    ></Table>

  </div>
</template>

<script setup lang="ts">
import {defineProps, reactive} from "vue";
import {PageTab} from "@/store/tabs";
import Table from "./Table.vue";
import {useI18n} from "vue-i18n";
const {t} = useI18n();

const table = reactive({
  isLoading: false,
  columns: [
    {
      label: "ID",
      field: "id",
      width: "3%",
      sortable: true,
      isKey: true,
    },
    {
      label: "Name",
      field: "name",
      width: "10%",
      sortable: true,
    },
    {
      label: "Email",
      field: "email",
      width: "15%",
      sortable: true,
    },
  ],
  rows: [],
  totalRecordCount: 0,
  sortable: {
    order: "id",
    sort: "asc",
  } ,
} as any);

const doSearch = (offset, limit, order, sort) => {
  console.log('===', limit)

  table.isLoading = true;
  setTimeout(() => {
    table.isReSearch = offset == undefined ? true : false;

    if (sort == "asc") {
      table.rows = sampleData1(offset, limit);
    } else {
      table.rows = sampleData2(offset, limit);
    }
    table.totalRecordCount = 20;
    table.sortable.order = order;
    table.sortable.sort = sort;
  }, 50);
};

doSearch(0, 3, 'id', 'asc');

// Fake Data for 'asc' sortable
const sampleData1 = (offst, limit) => {
  offst = offst + 1;
  let data = [] as any[];
  for (let i = offst; i < offst + limit; i++) {
    data.push({
      id: i,
      name: "TEST" + i,
      email: "test" + i + "@example.com",
    });
  }
  return data;
};
// Fake Data for 'desc' sortable
const sampleData2 = (offst, limit) => {
  let data = [] as any[];
  for (let i = limit; i > offst + limit; i--) {
    data.push({
      id: i,
      name: "TEST" + i,
      email: "test" + i + "@example.com",
    });
  }
  return data;
};

defineProps<{
  tab: PageTab
}>();
</script>

<style lang="less" scoped>
.tab-page-sites {
  padding: 10px;
  height: 100%;
  overflow-y: auto;
}
</style>
