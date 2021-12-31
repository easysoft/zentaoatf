<template>
    <div class="indexlayout-main-conent">
        <a-spin :spinning="loading" size="large">
            <a-card :bordered="false" title="退款申请" style="margin-bottom: 20px">

                <a-descriptions >
                    <a-descriptions-item label="取货单号">
                    {{refundApplication.ladingNo}}
                    </a-descriptions-item>
                    <a-descriptions-item label="状态">
                    {{refundApplication.state}}
                    </a-descriptions-item>
                    <a-descriptions-item label="销售单号">
                    {{refundApplication.saleNo}}
                    </a-descriptions-item>
                    <a-descriptions-item label="子订单">
                    {{refundApplication.childOrders}}
                    </a-descriptions-item>
                </a-descriptions>
            </a-card>

            <a-card :bordered="false" title="用户信息" style="margin-bottom: 20px">

                <a-descriptions>
                    <a-descriptions-item label="用户姓名">
                    {{userInfo.name}}
                    </a-descriptions-item>
                    <a-descriptions-item label="联系电话">
                    {{userInfo.tel}}
                    </a-descriptions-item>
                    <a-descriptions-item label="常用快递">
                    {{userInfo.courier}}
                    </a-descriptions-item>
                    <a-descriptions-item label="取货地址">
                    {{userInfo.address}}
                    </a-descriptions-item>
                    <a-descriptions-item label="备注">
                    {{userInfo.remark}}
                    </a-descriptions-item>
                </a-descriptions>

            </a-card>
            <a-card :bordered="false" title="退货商品" style="margin-bottom: 20px">
                <a-table
                    rowKey="id"
                    :pagination="false"
                    :dataSource="goodsData"
                    :columns="goodsColumns"
                />
            </a-card>

            <a-card :bordered="false" title="退货进度">
                <a-table
                    :pagination="false"
                    :dataSource="returnProgress"
                    :columns="progressColumns"
                >
                    <template #status="{ text }">
                        <a-badge v-if="text === 'success'" status="success" text="成功" />
                        <a-badge v-else status="processing" text="进行中" />;
                    </template>
                </a-table>
            </a-card>

        </a-spin>        
    </div>
</template>
<script lang="ts">
import { computed, defineComponent, onMounted, ref, h, Ref, ComputedRef } from "vue";
import { useStore } from 'vuex';
import { StateType as DetailStateType } from './store';
import { RefundApplicationDataType, ReturnGoodsDataType, ReturnProgressDataType, UserInfoDataType } from './data.d';

interface DetailBasicPageSetupData {
    loading: Ref<boolean>;
    refundApplication: ComputedRef<RefundApplicationDataType>;
    userInfo: ComputedRef<UserInfoDataType>;
    goodsColumns: any;
    goodsData: ComputedRef<ReturnGoodsDataType[]>;
    progressColumns: any;
    returnProgress: ComputedRef<ReturnProgressDataType[]>;
}

export default defineComponent({
    name: 'DetailModulePage',
    setup(): DetailBasicPageSetupData {
        const store = useStore<{ DetailModule: DetailStateType}>();

        // 退款申请 信息
        const refundApplication = computed<RefundApplicationDataType>(() => store.state.DetailModule.detail.refundApplication);

        // 用户信息
        const userInfo = computed<UserInfoDataType>(() => store.state.DetailModule.detail.userInfo);

        // 退货商品
        const returnGoods = computed<ReturnGoodsDataType[]>(() => store.state.DetailModule.detail.returnGoods);
        const goodsData = computed<ReturnGoodsDataType[]>(() => {
            let goodsData: typeof returnGoods.value = [];
            if (returnGoods.value.length) {
                let num = 0;
                let amount = 0;
                returnGoods.value.forEach(item => {
                    num += Number(item.num);
                    amount += Number(item.amount);
                });
                goodsData = returnGoods.value.concat({
                    id: '总计',
                    num,
                    amount,
                });
            }
            return goodsData;
        });
        const renderContent = ({text, index}: {text: any; index: number}) => {
            const obj: {
                children: any;
                props: { colSpan?: number };
            } = {
                children: text,
                props: {},
            };
            if (index === returnGoods.value.length) {
                obj.props.colSpan = 0;
            }
            return obj;
        };
        const goodsColumns = [
            {
                title: '商品编号',
                dataIndex: 'id',
                customRender: ({text, index}: { text: any; index: number}) => {
                    if (index < returnGoods.value.length) {
                        return h('a',{
                            href: 'javascript:;'
                        }, text);
                    }
                    return {
                        children: h('span',{style:"font-weight: 600"},'总计'),
                        props: {
                            colSpan: 4,
                        },
                    };
                },
            },
            {
                title: '商品名称',
                dataIndex: 'name',
                key: 'name',
                customRender: renderContent,
            },
            {
                title: '商品条码',
                dataIndex: 'barcode',
                key: 'barcode',
                customRender: renderContent,
            },
            {
                title: '单价',
                dataIndex: 'price',
                key: 'price',
                align: 'right' as 'left' | 'right' | 'center',
                customRender: renderContent,
            },
            {
                title: '数量（件）',
                dataIndex: 'num',
                key: 'num',
                align: 'right' as 'left' | 'right' | 'center',
                customRender: ({text, index}: { text: any; index: number}) => {
                    if (index < returnGoods.value.length) {
                        return text;
                    }
                    return h('span',{style:"font-weight: 600"},text)
                },
            },
            {
                title: '金额',
                dataIndex: 'amount',
                key: 'amount',
                align: 'right' as 'left' | 'right' | 'center',
                customRender: ({text, index}: { text: any; index: number}) => {
                    if (index < returnGoods.value.length) {
                        return text;
                    }
                    return h('span',{style:"font-weight: 600"},text)
                },
            },
        ];

        // 退货进度
        const returnProgress = computed<ReturnProgressDataType[]>(() => store.state.DetailModule.detail.returnProgress);
        const progressColumns = [
            {
                title: '时间',
                dataIndex: 'time',
            },
            {
                title: '当前进度',
                dataIndex: 'rate',
            },
            {
                title: '状态',
                dataIndex: 'status',
                slots: { customRender: 'status' },
            },

            {
                title: '操作员ID',
                dataIndex: 'operator',
            },
            {
                title: '耗时',
                dataIndex: 'cost',
            },
        ];


        // 读取数据 func
        const loading = ref<boolean>(true);
        const getData = async () => {
            loading.value = true;
            await store.dispatch('DetailModule/queryDetail');
            loading.value = false;
        }

        onMounted(()=> {
           getData();
        })


        return {
            loading,
            refundApplication,
            userInfo,
            goodsColumns,
            goodsData,
            progressColumns,
            returnProgress
        }
    }
})
</script>