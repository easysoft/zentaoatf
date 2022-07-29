<template>
<div v-if="statistic != undefined" class="statistic">
    <span class="statistic-total"><Icon icon="play" />{{statistic.total}}</span>
    <span class="statistic-succ"><Icon icon="checkmark-circle" />{{statistic.succ}}</span>
    <span class="statistic-fail"><Icon icon="close-circle" />{{statistic.fail}}</span>
    <span class="statistic-bug"><Icon icon="bug" />{{statistic.fail}}</span>
</div>
</template>

<script setup lang="ts">
import {useStore} from "vuex";
import {computed} from "vue";
import {useI18n} from "vue-i18n";
import Icon from '@/components/Icon.vue';
import {ScriptData} from "@/views/script/store";

const { t } = useI18n();

const store = useStore<{ Script: ScriptData }>();

const statistic = computed<any>(() => store.state.Script.statistic);

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
    }
    .statistic-bug{
        color: var(--color-red);
    }
    .statistic-total{
        color: var(--color-secondary);
    }
}
</style>
