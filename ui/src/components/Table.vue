<template>
  <div class="z-table-main">
    <div class="z-row">
      <div class="z-full">
        <div v-if="isLoading" class="z-loading-mask">
          <div class="z-loading-content">
            <span style="color: white">Loading...</span>
          </div>
        </div>
        <table
            class="z-table z-table-hover z-table-bordered"
            ref="localTable">
          <thead class="z-thead">
            <tr class="z-thead-tr">
            <th v-if="hasCheckbox" class="z-thead-th z-checkbox-th">
              <div>
                <input
                    type="checkbox"
                    class="z-thead-checkbox"
                    v-model="setting.isCheckAll"
                    @click="checkAll"
                />
              </div>
            </th>
            <th v-for="(col, index) in columns"
                class="z-thead-th"
                :class="col.headerClasses"
                :key="index"
                :style="
                  Object.assign(
                    { width: col.width ? col.width : 'auto' },
                    col.headerStyles
                  )
                ">
              <div class="z-thead-column"
                   :class="{
                    'z-sortable': col.sortable,
                    'z-both': col.sortable,
                    'z-asc': setting.order === col.field && setting.sort === 'asc',
                    'z-desc': setting.order === col.field && setting.sort === 'desc',
                  }"
                   @click="col.sortable ? doSort(col.field) : false">
                {{ col.label }}
              </div>
            </th>
          </tr>
          </thead>
          <tbody v-if="rows.length > 0" class="z-tbody">
            <template v-if="isStaticMode">
            <tr v-for="(row, i) in localRows"
                :key="i"
                class="z-tbody-tr"
                :class="typeof rowClasses === 'function' ? rowClasses(row) : rowClasses"
                @click="$emit('row-clicked', row)">
              <td v-if="hasCheckbox" class="z-tbody-td">
                <div>
                  <input type="checkbox"
                         class="z-tbody-checkbox"
                         :ref="
                        (el) => {
                          rowCheckbox[i] = el;
                        }
                      "
                         :value="row[setting.keyColumn]"
                         :checked="row.checked"
                         @click="checked"/>
                </div>
              </td>
              <td v-for="(col, j) in columns"
                  :key="j"
                  class="z-tbody-td"
                  :class="col.columnClasses"
                  :style="col.columnStyles">
                <div v-if="col.display" v-html="col.display(row)"></div>
                <template v-else>
                  <div v-if="setting.isSlotMode && slots[col.field]">
                    <slot :name="col.field" :value="row"></slot>
                  </div>
                  <span v-else>{{ row[col.field] }}</span>
                </template>
              </td>
            </tr>
          </template>
            <template v-else>
            <tr v-for="(row, i) in rows"
                :key="i"
                class="z-tbody-tr"
                :class="typeof rowClasses === 'function' ? rowClasses(row) : rowClasses"
                @click="$emit('row-clicked', row)">
              <td v-if="hasCheckbox" class="z-tbody-td">
                <div>
                  <input type="checkbox"
                         class="z-tbody-checkbox"
                         :ref="
                        (el) => {
                          rowCheckbox[i] = el;
                        }
                      "
                         :value="row[setting.keyColumn]"
                         :checked="row.checked"
                         @click="checked"/>
                </div>
              </td>
              <td v-for="(col, j) in columns"
                  :key="j"
                  class="z-tbody-td"
                  :class="col.columnClasses"
                  :style="col.columnStyles">
                <div v-if="col.display" v-html="col.display(row)"></div>
                <div v-else>
                  <div v-if="setting.isSlotMode && slots[col.field]">
                    <slot :name="col.field" :value="row"></slot>
                  </div>
                  <span v-else>{{ row[col.field] }}</span>
                </div>
              </td>
            </tr>
          </template>
          </tbody>
        </table>
      </div>
    </div>

    <!-- pagination -->
    <div class="z-paging z-row" v-if="rows.length > 0 && !setting.isHidePaging">
        <div class="z-paging-info">
          <div role="status" aria-live="polite">
            {{ info.pagingInfo }}
          </div>
        </div>

        <div class="z-paging-change-div">
          <span class="z-paging-count-label">{{ info.pageSizeChangeLabel }}&nbsp;</span>
          <select class="z-paging-count-dropdown" v-model="setting.pageSize">
            <option v-for="pageOption in pageOptions"
                    :value="pageOption.value"
                    :key="pageOption.value">
              {{ pageOption.text }}
            </option>
          </select>&nbsp;

          <span class="z-paging-page-label">{{ info.gotoPageLabel }}&nbsp;</span>
          <select class="z-paging-page-dropdown" v-model="setting.page">
            <option v-for="n in setting.maxPage" :key="n" :value="parseInt(n)">
              {{ n }}
            </option>
          </select>
        </div>

        <div class="z-paging-pagination-div col-full col-md-4">
          <div class="dataTables_paginate">
            <ul class="z-paging-pagination-ul z-pagination">
              <li class="z-paging-pagination-page-li z-paging-pagination-page-li-first page-item"
                  :class="{ disabled: setting.page <= 1 }">
                <a class="z-paging-pagination-page-link z-paging-pagination-page-link-first page-link"
                   href="javascript:void(0)"
                   aria-label="Previous"
                   @click="setting.page = 1">
                  <span aria-hidden="true">&laquo;</span>
                  <span class="sr-only">First</span>
                </a>
              </li>
              <li class="z-paging-pagination-page-li z-paging-pagination-page-li-prev page-item"
                  :class="{ disabled: setting.page <= 1 }">
                <a class="z-paging-pagination-page-link z-paging-pagination-page-link-prev page-link"
                   href="javascript:void(0)"
                   aria-label="Previous"
                   @click="prevPage">
                  <span aria-hidden="true">&lt;</span>
                  <span class="sr-only">Prev</span>
                </a>
              </li>
              <li class="z-paging-pagination-page-li z-paging-pagination-page-li-number page-item"
                  v-for="n in setting.paging"
                  :key="n"
                  :class="{ disabled: setting.page === n }">
                <a class="z-paging-pagination-page-link z-paging-pagination-page-link-number page-link"
                   href="javascript:void(0)"
                   @click="movePage(n)"
                >{{ n }}</a>
              </li>
              <li
                  class="z-paging-pagination-page-li z-paging-pagination-page-li-next page-item"
                  :class="{ disabled: setting.page >= setting.maxPage }">
                <a class="z-paging-pagination-page-link z-paging-pagination-page-link-next page-link"
                   href="javascript:void(0)"
                   aria-label="Next"
                   @click="nextPage">
                  <span aria-hidden="true">&gt;</span>
                  <span class="sr-only">Next</span>
                </a>
              </li>
              <li class="z-paging-pagination-page-li z-paging-pagination-page-li-last page-item"
                  :class="{ disabled: setting.page >= setting.maxPage }">
                <a class="z-paging-pagination-page-link z-paging-pagination-page-link-last page-link"
                   href="javascript:void(0)"
                   aria-label="Next"
                   @click="setting.page = setting.maxPage">
                  <span aria-hidden="true">&raquo;</span>
                  <span class="sr-only">Last</span>
                </a>
              </li>
            </ul>
          </div>
        </div>
    </div>

    <div class="z-row" v-else>
      <div class="z-empty-msg col-full z-center">
        {{ messages.noDataAvailable }}
      </div>
    </div>
  </div>
</template>

<script lang="ts">
// https://github.com/linmasahiro/vue3-table-lite

import {computed, defineComponent, nextTick, onBeforeUpdate, onMounted, reactive, ref, watch,} from "vue";
import {useI18n} from "vue-i18n";

interface PageOption {
  value: number;
  text: number | string;
}

interface TableSetting {
  isSlotMode: boolean;
  isCheckAll: boolean;
  isHidePaging: boolean;
  keyColumn: string;
  page: number;
  pageSize: number;
  maxPage: number;
  offset: number;
  limit: number;
  paging: Array<number>;
  order: string;
  sort: string;
  pageOptions: Array<PageOption>;
}

interface Column {
  isKey: string;
  field: string;
}

interface PageMessage {
  pagingInfo: string
  pageSizeChangeLabel: string
  gotoPageLabel: string
  noDataAvailable: string
}

export default defineComponent({
  name: "Table",
  emits: ["return-checked-rows", "do-search", "is-finished", "get-now-page", "row-clicked"],
  props: {
    // is data loading
    isLoading: {
      type: Boolean,
      require: true,
    },
    // Whether to perform a re-query
    isReSearch: {
      type: Boolean,
      require: true,
    },
    // Presence of Checkbox
    hasCheckbox: {
      type: Boolean,
      default: false,
    },
    isCheckAll: {
      type: Boolean,
      default: false,
    },
    //Returns data type for checked of Checkbox
    checkedReturnType: {
      type: String,
      default: "key",
    },
    // title
    title: {
      type: String,
      default: "",
    },
    // Field
    columns: {
      type: Array,
      default: () => {
        return [];
      },
    },
    // data
    rows: {
      type: Array,
      default: () => {
        return [];
      },
    },
    // data row classes
    rowClasses: {
      type: [Array, Function],
      default: () => {
        return [];
      },
    },
    // Display the number of items on one page
    pageSize: {
      type: Number,
      default: 10,
    },
    // Total number of transactions
    total: {
      type: Number,
      default: 100,
    },
    // Current page number
    page: {
      type: Number,
      default: 1,
    },
    // Sort condition
    sortable: {
      type: Object,
      default: () => {
        return {
          order: "id",
          sort: "asc",
        };
      },
    },
    // Display text
    messages: {
      type: Object,
      default: {} as PageMessage,
    },
    // Static mode(no refresh server data)
    isStaticMode: {
      type: Boolean,
      default: false,
    },
    // V-slot mode
    isSlotMode: {
      type: Boolean,
      default: false,
    },
    // Hide paging
    isHidePaging: {
      type: Boolean,
      default: false,
    },
    // Modify page dropdown
    pageOptions: {
      type: Array,
      default: () => [
        {
          value: 10,
          text: 10,
        },
        {
          value: 20,
          text: 20,
        },
        {
          value: 50,
          text: 50,
        },
      ],
    },
  },
  setup(props, {emit, slots}) {
    const {t} = useI18n();

    const info = computed<any>(() => {
      return {
        pagingInfo: props.messages?.pagingInfo ? props.messages?.pagingInfo :
            t('page_info', {offset: setting.offset, limit: setting.limit, total: props.total}),
        pageSizeChangeLabel: props.messages?.pageSizeChangeLabel ? props.messages?.pageSizeChangeLabel :
            t('page_count'),
        gotoPageLabel: props.messages?.gotoPageLabel ? props.messages?.gotoPageLabel :
            t('page_goto'),
        noDataAvailable: props.messages?.noDataAvailable ? props.messages?.noDataAvailable :
            t('page_no_data'),
      }
    })

    let localTable = ref<HTMLElement | null>(null);

    // Validate dropdown values have page-size value or not
    let tmpPageOptions = props.pageOptions as Array<PageOption>;
    let defaultPageSize =
        props.pageOptions.length > 0 ? ref(tmpPageOptions[0].value) : ref(props.pageSize);
    if (tmpPageOptions.length > 0) {
      tmpPageOptions.forEach((v: PageOption) => {
        if (
            Object.prototype.hasOwnProperty.call(v, "value") &&
            Object.prototype.hasOwnProperty.call(v, "text") &&
            props.pageSize == v.value
        ) {
          defaultPageSize.value = v.value;
        }
      });
    }

    // Internal set value for components
    const setting: TableSetting = reactive({
      // Enable slot mode
      isSlotMode: props.isSlotMode,
      // Whether to select all
      isCheckAll: props.isCheckAll,
      // Hide paging
      isHidePaging: props.isHidePaging,
      // KEY field name
      keyColumn: computed(() => {
        let key = "";
        Object.assign(props.columns).forEach((col: Column) => {
          if (col.isKey) {
            key = col.field;
          }
        });
        return key;
      }),
      // current page number
      page: props.page,
      // Display count per page
      pageSize: !props.isHidePaging ? props.pageSize ? props.pageSize : defaultPageSize.value : 10000,
      // Maximum number of pages
      maxPage: computed(() => {
        if (props.total <= 0) {
          return 0;
        }
        let maxPage = Math.floor(props.total / setting.pageSize);
        let mod = props.total % setting.pageSize;
        if (mod > 0) {
          maxPage++;
        }
        return maxPage;
      }),
      // The starting value of the page number
      offset: computed(() => {
        return (setting.page - 1) * setting.pageSize + 1;
      }),
      // Maximum number of pages0
      limit: computed(() => {
        let limit = setting.page * setting.pageSize;
        return props.total >= limit ? limit : props.total;
      }),
      // Paging array
      paging: computed(() => {
        let startPage = setting.page - 2 <= 0 ? 1 : setting.page - 2;
        if (setting.maxPage - setting.page <= 2) {
          startPage = setting.maxPage - 4;
        }
        startPage = startPage <= 0 ? 1 : startPage;
        let pages = [] as number[];
        for (let i = startPage; i <= setting.maxPage; i++) {
          if (pages.length < 5) {
            pages.push(i);
          }
        }
        return pages;
      }),
      // Sortable for local
      order: props.sortable.order,
      sort: props.sortable.sort,
      pageOptions: computed(() => {
        const ops: PageOption[] = [];
        props.pageOptions?.forEach((o) => {
          ops.push({
            value: (o as PageOption).value,
            text: (o as PageOption).text,
          });
        });
        return ops;
      }),
    });

    // Data rows for local
    const localRows = computed(() => {
      // sort rows
      let property = setting.order;
      let sort_order = 1;
      if (setting.sort === "desc") {
        sort_order = -1;
      }
      let rows = props.rows
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      rows.sort((a: any, b: any): number => {
        if (a[property] < b[property]) {
          return -1 * sort_order;
        } else if (a[property] > b[property]) {
          return sort_order;
        } else {
          return 0;
        }
      });

      // return sorted and offset rows
      let result = [] as any[];
      for (let index = setting.offset - 1; index < setting.limit; index++) {
        if (rows[index]) {
          result.push(rows[index]);
        }
      }

      nextTick(function () {
        // 資料完成渲染後回傳私有元件
        callIsFinished();
      });

      return result;
    });

    ////////////////////////////
    //
    //  (Checkbox related operations)
    //

    // 定義Checkbox參照 (Define Checkbox reference)
    const rowCheckbox = ref([]);
    if (props.hasCheckbox) {
      /**
       * Execute before re-rendering
       */
      onBeforeUpdate(() => {
        // 每次更新前都把值全部清空 (Clear all values before each update)
        rowCheckbox.value = [];
      });

      /**
       * Check all checkboxes for monitoring
       */
    //   watch(
    //       () => setting.isCheckAll,
    //       (state: boolean) => {
    //         let isChecked: Array<string | unknown> = [];
    //         rowCheckbox.value.forEach((val: HTMLInputElement, i: number) => {
    //           if (val) {
    //             val.checked = state;
    //             if (val.checked) {
    //               if (props.checkedReturnType == "row") {
    //                 isChecked.push(localRows.value[i]);
    //               } else {
    //                 isChecked.push(val.value);
    //               }
    //             }
    //           }
    //         });
    //         // 回傳畫面上選上的資料 (Return the selected data on the screen)
    //         emit("return-checked-rows", isChecked);
    //       }
    //   );
    }

    const checkAll = (event: Event) => {
        event.stopPropagation();
        setting.isCheckAll = !setting.isCheckAll;
        let isChecked: Array<string | unknown> = [];
        rowCheckbox.value.forEach((val: HTMLInputElement, i: number) => {
            val.checked = setting.isCheckAll;
            if (val.checked) {
                if (props.checkedReturnType == "row") {
                isChecked.push(localRows.value[i]);
                } else {
                isChecked.push(val.value);
                }
            }
        });
        // 回傳畫面上選上的資料 (Return the selected data on the screen)
        emit("return-checked-rows", isChecked);
    }

    /**
     * Checkbox click event
     */
    const checked = (event: Event) => {
      event.stopPropagation();
      let isChecked: Array<string | unknown> = [];
      rowCheckbox.value.forEach((val: HTMLInputElement, i: number) => {
        if (val.checked) {
          if (props.checkedReturnType == "row") {
            isChecked.push(localRows.value[i]);
          } else {
            isChecked.push(val.value);
          }
        }
      });
      setting.isCheckAll = isChecked.length >= rowCheckbox.value.length;
      // 回傳畫面上選上的資料 (Return the selected data on the screen)
      emit("return-checked-rows", isChecked);
    };

    /**
     * Clear all selected data on the screen
     */
    const clearChecked = () => {
      rowCheckbox.value.forEach((val: HTMLInputElement) => {
        if (val && val.checked) {
          val.checked = false;
        }
      });
      // 回傳畫面上選上的資料 (Return the selected data on the screen)
      emit("return-checked-rows", []);
    };

    ////////////////////////////
    //
    //  (Sorting, page change, etc. related operations)
    //

    /**
     * Call execution sequencing
     */
    const doSort = (order: string) => {
      let sort = "asc";
      if (order == setting.order) {
        // 排序中的項目時 (When sorting items)
        if (setting.sort == "asc") {
          sort = "desc";
        }
      }
      let offset = (setting.page - 1) * setting.pageSize;
      let limit = setting.pageSize;
      setting.order = order;
      setting.sort = sort;

      emit("do-search", offset, limit, order, sort);

      // Clear the selected data on the screen
      if (setting.isCheckAll) {
        // It will be cleared when you cancel all selections
        setting.isCheckAll = false;
      } else {
        if (props.hasCheckbox) {
          clearChecked();
        }
      }
    };

    /**
     * Switch page number
     *
     * @param pageNum      number  New page number
     * @param prevPageNum  number  Current page number
     */
    const changePage = (pageNum: number, prevPageNum: number) => {
      setting.isCheckAll = false;
      let order = setting.order;
      let sort = setting.sort;
      let offset = (pageNum - 1) * setting.pageSize;
      let limit = setting.pageSize;
      if (!props.isReSearch || pageNum > 1 || pageNum == prevPageNum) {
        // Call query will only be executed if the page number is changed without re-query

        console.log("do-search", offset, limit)
        emit("do-search", offset, limit, order, sort);
      }
    };
    // Monitor page switching
    watch(() => setting.page, changePage);
    // Monitor manual page switching
    watch(
        () => props.page,
        (val) => {
          if (val <= 1) {
            setting.page = 1;
            emit("get-now-page", setting.page);
          } else if (val >= setting.maxPage) {
            setting.page = setting.maxPage;
            emit("get-now-page", setting.page);
          } else {
            setting.page = val;
          }
        }
    );

    /**
     * Switch display number
     */
    const changePageSize = () => {
      if (setting.page === 1) {
        changePage(setting.page, setting.page);
      } else {
        setting.page = 1;
        setting.isCheckAll = false;
      }
    };
    // Monitor display number switch
    watch(() => setting.pageSize, changePageSize);

    /**
     * Previous page
     */
    const prevPage = () => {
      if (setting.page == 1) {
        // If it is the first page, it will not be executed
        return false;
      }
      setting.page--;
    };

    /**
     * Move to the specified number of pages
     */
    const movePage = (page: number) => {
      setting.page = page;
    };

    /**
     * Next page
     */
    const nextPage = () => {
      if (setting.page >= setting.maxPage) {
        // 如果等於大於最大頁數，不與執行 (If it is equal to or greater than the maximum number of pages, no execution)
        return false;
      }
      setting.page++;
    };

    // Monitoring data changes
    watch(
        () => props.rows,
        () => {
          if (props.isReSearch || props.isStaticMode) {
            setting.page = 1;
          }
          nextTick(function () {
            // 資料完成渲染後回傳私有元件 (Return the private components after the data is rendered)
            if (!props.isStaticMode) {
              callIsFinished();
            }
          });
        }
    );

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const stringFormat = (template: string, ...args: any[]) => {
      return template.replace(/{(\d+)}/g, function (match, number) {
        return typeof args[number] != "undefined" ? args[number] : match;
      });
    };

    // Call 「is-finished」 Method
    const callIsFinished = () => {
      if (localTable.value) {
        let localElement = localTable.value.getElementsByClassName("is-rows-el");
        emit("is-finished", localElement);
      }
      emit("get-now-page", setting.page);
    };

    /**
     * Mounted Event
     */
    onMounted(() => {
      nextTick(() => {
        if (props.rows.length > 0) {
          callIsFinished();
        }
      });
    });

    if (props.hasCheckbox) {
      // 需要 Checkbox 時 (When Checkbox is needed)
      return {
        slots,
        localTable,
        localRows,
        setting,
        rowCheckbox,
        checked,
        checkAll,
        doSort,
        prevPage,
        movePage,
        nextPage,
        stringFormat,
      };
    } else {
      return {
        slots,
        localTable,
        localRows,
        setting,
        doSort,
        prevPage,
        movePage,
        nextPage,
        stringFormat,
        info,
      };
    }
  },
});
</script>

<style lang="less" scoped>
.z-table-main {

}
.z-checkbox-th {
  width: 1%;
}

.z-sortable {
  cursor: pointer;
  background-position: right;
  background-repeat: no-repeat;
  padding-right: 30px !important;

  &.z-asc {
    background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAYAAAByUDbMAAAAZ0lEQVQ4y2NgGLKgquEuFxBPAGI2ahhWCsS/gDibUoO0gPgxEP8H4ttArEyuQYxAPBdqEAxPBImTY5gjEL9DM+wTENuQahAvEO9DMwiGdwAxOymGJQLxTyD+jgWDxCMZRsEoGAVoAADeemwtPcZI2wAAAABJRU5ErkJggg==) no-repeat right 5px;
  }

  &.z-desc {
    background: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAYAAAByUDbMAAAAZUlEQVQ4y2NgGAWjYBSggaqGu5FA/BOIv2PBIPFEUgxjB+IdQPwfC94HxLykus4GiD+hGfQOiB3J8SojEE9EM2wuSJzcsFMG4ttQgx4DsRalkZENxL+AuJQaMcsGxBOAmGvopk8AVz1sLZgg0bsAAAAASUVORK5CYII=) no-repeat right -2px;
  }
}

.z-both {
  background-image: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABMAAAATCAQAAADYWf5HAAAAkElEQVQoz7X QMQ5AQBCF4dWQSJxC5wwax1Cq1e7BAdxD5SL+Tq/QCM1oNiJidwox0355mXnG/DrEtIQ6azioNZQxI0ykPhTQIwhCR+BmBYtlK7kLJYwWCcJA9M4qdrZrd8pPjZWPtOqdRQy320YSV17OatFC4euts6z39GYMKRPCTKY9UnPQ6P+GtMRfGtPnBCiqhAeJPmkqAAAAAElFTkSuQmCC");
}



.z-loading-mask {
  position: absolute;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  flex-flow: column;
  transition: opacity 0.3s ease;
}

.z-loading-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

select {
  width: auto;
  border: 1px solid #cccccc;
  background-color: #ffffff;
  height: auto;
  padding: 0;
  margin-bottom: 0;
}

.z-table {
  width: 100%;
  margin-bottom: 1rem;
  color: #212529;
  border-collapse: collapse;
}

th {
  text-align: inherit;
}

tr {
  display: table-row;
  vertical-align: inherit;
  border-color: inherit;
}

.z-table-bordered thead td,
.z-table-bordered thead th {
  border-bottom-width: 2px;
}

.z-table thead th {
  vertical-align: bottom;
  background-color: var(--color-darken-1);
  border-color: #dee2e6;
  border-bottom: 2px solid #dee2e6;
}

.z-table-bordered td,
.z-table-bordered th {
  border: 1px solid #dee2e6;
}

.z-table td,
.z-table th {
  padding: 0.75rem;
  vertical-align: top;
  border-top: 1px solid #dee2e6;
  vertical-align: middle;
  position: sticky;
  top: 0;
}

.z-table-hover tbody tr:hover {
  color: #212529;
  background-color: rgba(0, 0, 0, 0.075);
}

.z-row {
  display: -ms-flexbox;
  display: flex;
  -ms-flex-wrap: wrap;
  flex-wrap: wrap;
}

.z-pagination {
  margin: 2px 0;
  white-space: nowrap;
  justify-content: flex-end;
  display: -ms-flexbox;
  display: flex;
  padding-left: 0;
  list-style: none;
  border-radius: 0.25rem;
}

.z-paging {
  .z-paging-info {
    width: 160px;
    line-height: 38px;

    .z-paging-page-dropdown {
      width: 50px !important;
    }
  }

  .z-paging-change-div {
    width: 230px;
    line-height: 38px;

    .z-paging-page-dropdown {
      width: 50px !important;
    }
  }

  .z-paging-pagination-div {
    flex: 1;

    .page-link {
      position: relative;
      display: block;
      padding: 0.5rem 0.75rem;
      margin-left: -1px;
      line-height: 1.25;
      color: #007bff;
      background-color: #fff;
      border: 1px solid #dee2e6;
    }

    .page-item {
      &.disabled .page-link {
        color: #6c757d;
        pointer-events: none;
        cursor: auto;
        background-color: #fff;
        border-color: #dee2e6;
      }

      &:first-child .page-link {
        margin-left: 0;
        border-top-left-radius: 0.25rem;
        border-bottom-left-radius: 0.25rem;
      }
    }

    .sr-only {
      position: absolute;
      width: 1px;
      height: 1px;
      padding: 0;
      margin: -1px;
      overflow: hidden;
      clip: rect(0, 0, 0, 0);
      white-space: nowrap;
      border: 0;
    }
  }
}

</style>
