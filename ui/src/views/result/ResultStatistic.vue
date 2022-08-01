<template>
<div v-if="statistic != undefined" class="statistic">
    <span class="statistic-total"><Icon icon="play" />{{statistic.total}}</span>
    <span class="statistic-succ"><Icon icon="checkmark-circle" />{{statistic.succ}}</span>
    <span class="statistic-fail" @click="showFailureModal=statistic.fail>0"><Icon icon="close-circle" />{{statistic.fail}}</span>
    <span class="statistic-bug" @click="showBugsModal=bugCount>0"><Icon icon="bug" />{{bugCount}}</span>
    <BugsModal
      v-if="showBugsModal"
      @cancel="showBugsModal=!showBugsModal"
      :caseIds="caseIds"
    />
    <FailureModal
      v-if="showFailureModal"
      @cancel="showFailureModal=!showFailureModal"
      :path="props.path"
    />
</div>
</template>

<script setup lang="ts">
import {useStore} from "vuex";
import {computed, ref, withDefaults, defineProps} from "vue";
import {useI18n} from "vue-i18n";
import Icon from '@/components/Icon.vue';
import {StateType} from "@/views/result/store";
import {ZentaoData} from "@/store/zentao";
import { ScriptData } from "@/views/script/store";
import BugsModal from "@/views/result/BugsModal.vue";
import FailureModal from "@/views/result/FailureModal.vue";

export interface ResultStatisticProps {
  path: string;
}
const props = withDefaults(defineProps<ResultStatisticProps>(), {
  path: '',
});

const { t } = useI18n();

const store = useStore<{ Result: StateType, Zentao: ZentaoData, Script: ScriptData }>();

const statistic = computed<any>(() => store.state.Result.statistic);
const bugMap = computed<any>(() => store.state.Zentao.bugMap);
const treeDataMap = computed<any>(() => store.state.Script.treeDataMap);
const showFailureModal = ref(false)
const bugCount = computed(() => {
    if(treeDataMap.value[props.path] == undefined || bugMap.value[treeDataMap.value[props.path].caseId] == undefined) {
        return 0;
    }
    return bugMap.value[treeDataMap.value[props.path].caseId].length
});
const caseIds = computed(() => {
    if(treeDataMap.value[props.path] == undefined || bugMap.value[treeDataMap.value[props.path].caseId] == undefined) {
        return [];
    }
    return [treeDataMap.value[props.path].caseId]
});
const showBugsModal = ref(false)
</script>

<style lang="less" scoped>
.statistic{
    display: flex;
    flex: 1;
    justify-content: flex-start;
    span{
        display: flex;
        align-items: center;
        padding: 0 var(--space-base);
        svg{
            margin-right: 10px;
        }
    }
    .statistic-succ{
        color: var(--color-green);
    }
    .statistic-fail{
        color: var(--color-yellow);
        cursor: pointer;
    }
    .statistic-bug{
        color: var(--color-red);
        cursor: pointer;
    }
    .statistic-total{
        color: var(--color-secondary);
    }
}
</style>
