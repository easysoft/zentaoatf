<template>
    <div class="main-conent-screen">
     
        <div v-if="$slots.header" class="screen-header"><slot name="header"></slot></div>
        <div v-else class="screen-padding" />

        <div class="screen-conent" ref="conentRef">
            <a-table
                bordered
                :rowKey="rowKey"
                :columns="columnsRest"
                :dataSource="dataSource"
                :loading="loading"
                :pagination="false"
                :scroll="{ scrollToFirstRowOnChange: true, y: tableScrollY }"
            >
                <slot></slot>
            </a-table>
        </div>

        <div v-if="pagination" class="screen-footer">
          <a-pagination 
            :total="pagination.total" 
            :current="pagination.current" 
            :page-size="pagination.pageSize"
            :show-size-changer="pagination.showSizeChanger"
            :show-quick-jumper="pagination.showQuickJumper"
            @change="pagination.onChange"
          />
        </div>
        <div v-else class="screen-padding" />
    </div>
</template>
<script lang="ts">
import { computed, defineComponent, onBeforeUnmount, onMounted, PropType, ref } from "vue";
import debounce from "lodash.debounce";
import { PaginationConfig } from "./data.d";

export default defineComponent({
    name: 'ScreenTable',
    props: {
        rowKey: {
            type: String,
        },
        columns: {
            type: Array,
            required: true
        },
        dataSource: {
            type: Array
        },
        loading: {
            type: Boolean
        },
        pagination: {
            type: Object as PropType<PaginationConfig | false | undefined>
        }
    },
    setup(props, { slots }) {

        const columnsRest = computed<any>(() => {
            return props.columns.map((item: any)=> {
                    if(item['slots'] && item['slots']['customRender']) {
                        item['customRender'] = slots[item['slots']['customRender']];
                    }
                    return item;
            });
        })

        const tableScrollY = ref<number>(0);
        const conentRef = ref<HTMLDivElement>();

        const resizeHandler = debounce(() => {
            if (conentRef.value) {           
                tableScrollY.value = conentRef.value.offsetHeight - conentRef.value.getElementsByClassName('ant-table-thead')[0].clientHeight - 2;            
            }
        }, 100);

        onMounted(()=> {
            resizeHandler();

            window.addEventListener('resize', resizeHandler);

        })

        onBeforeUnmount(()=> {
             window.removeEventListener('resize', resizeHandler);
        })

        return {
            conentRef,
            tableScrollY,
            columnsRest
        }

    }
})
</script>
<style lang="less" scoped>
.main-conent-screen {
  display: flex;
  flex-direction: column;
  height: calc(100% - 48px - 50px);
  border-radius: 4px;
  background-color: #fff;
  .screen-header {
    padding: 20px;
    min-height: 33px;
  }
  .screen-footer {
    padding: 20px;
    min-height: 32px;
    text-align: right;
  }
  .screen-conent {
    flex: 1;
    padding: 0 20px;
    overflow: hidden;
  }
  .screen-padding {
    padding-top: 20px;
  }
  ::v-deep(.ant-table.ant-table-bordered) > .ant-table-content {
    border-bottom: 1px solid #f0f0f0;
  }
}
</style>