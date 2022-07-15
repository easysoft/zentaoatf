import { computed, ComputedRef, onMounted, watch } from 'vue';
import { useStore } from 'vuex';
import { StateType } from '@/views/site/store';
import {ZentaoData} from '@/store/zentao';
import {PaginationConfig} from '@/types/data';
import useCurrentProduct from './use-current-product';

interface TestResultInfo {
    buildTool: string;
    duration: number;
    endTime: number;
    execBy: string;
    fail: number;
    no: string;
    pass: number;
    seq: string;
    skip: number;
    startTime: number;
    testEnv: string;
    testTool: string;
    testType: string;
    total: number;
    workspaceId: number;
    workspaceName: string;
    testScriptName: string;
    displayName: string;
}

export default function useResultList(): {results: ComputedRef<TestResultInfo[]>, fetchResults: () => void} {
    const store = useStore<{ Zentao: ZentaoData, Result: StateType }>();
    const results = computed<any[]>(() =>
        store.state.Result.queryResult.result?.map((item) => {
            const displayName = item.testType === "unit" || item.total != 1 ? item.workspaceName + '(' + item.total + ')' : item.testScriptName
                return {
                    ...item,
                    displayName: displayName
                }
            }));
    const currentProduct = useCurrentProduct();

    const pagination = computed<PaginationConfig>(() => store.state.Result.queryResult.pagination);

    function fetchResults(page = 1) {
        const queryParams = {keywords: '', enabled: '1', page: pagination.value.page, pageSize: pagination.value.pageSize};
        store.dispatch('Result/list', {
            keywords: queryParams.keywords,
            enabled: queryParams.enabled,
            pageSize: pagination.value.pageSize,
            page: page
        });
    }

    onMounted(fetchResults);

    watch(currentProduct, () => fetchResults());

    return {results, fetchResults};
}
