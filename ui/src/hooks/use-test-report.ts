import {onMounted, watch, shallowRef} from 'vue';
import {useStore} from 'vuex';
import {StateType} from '@/views/result/store';

export default function useTestReport(seq: string, workspaceId: number): any {
    const store = useStore<{ Result: StateType }>();
    const reportRef = shallowRef();

    watch(() => store.state.Result.detailResult, (report) => {
        if (!report || seq !== report.seq || workspaceId !== report.workspaceId) {
            return;
        }
        reportRef.value = report;
    }, {deep: true});

    onMounted(() => {
        store.dispatch('Result/get', {workspaceId: workspaceId, seq: seq});
    });

    return reportRef;
}
