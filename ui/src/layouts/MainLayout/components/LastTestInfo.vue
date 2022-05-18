<template>
  <Button class="rounded border-light canvas gap-sm" :hint="t('view_test_result')" @click="showDetail(model)">
    <small class="text-muted">{{t('previous_result')}}</small>
    <small class="text-yellow" :key="model">{{model.duration}}s</small>
    <Icon icon="close-circle" class="text-red space-left" />
    <small class="text-red">{{model.fail}}</small>
    <Icon icon="checkmark-circle" class="text-green space-left" />
    <small class="text-green">{{model.pass}}</small>
  </Button>
</template>

<script setup lang="ts">
import Button from './Button.vue';
import Icon from './Icon.vue';
import {StateType} from "@/views/result/store";
import {useStore} from "vuex";
import {computed, onMounted} from "vue";
import {useI18n} from "vue-i18n";
const { t } = useI18n();

const store = useStore<{ Result: StateType }>();
const model = computed<any[]>(() => store.state.Result.lastResult)

const latest = () => {
    store.dispatch('Result/latest', {});
}
latest();

const showDetail = (item) => {
    console.log(item)
}

onMounted(() => {
    console.log("onMounted")
})

</script>

