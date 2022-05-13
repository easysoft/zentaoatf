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
import {StateType} from "@/src/views/result/store";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import {computed, onMounted} from "vue";
import {useI18n} from "vue-i18n";
const { t } = useI18n();

const store = useStore<{ Result: StateType }>();
const models = computed<any[]>(() => store.state.Result.queryResult.result)
var model = computed<any[]>(() => models.value.length > 0 ? models.value[0] : {})

const list = (page: number) => {
    store.dispatch('Result/list', {
    page: page});
}
list(1);

const showDetail = (item) => {
    console.log(item)
}

onMounted(() => {
    console.log("onMounted")
})

</script>

