<template>
  <div class="result-list">
    <List compact divider>
    <div v-for="item, index in models" :key="index" :class="'list-item-container ' + (item.checked==1?'checked':'')" @click="showDetail($event, index)" @mouseenter="changeControlIcon($event, index)" @mouseleave="changeControlIcon($event, index)">
        <ListItem
          icon="checkmark-circle"
          class="inline-left"
          iconClass="text-green"
          v-if="item.fail==0"
          :title="item.total > 1 ? item.workspaceName + '(' + item.total + ')' : item.testScriptName"
          trailingTextClass="muted small"
        >
        </ListItem>
        <ListItem
          icon="close-circle"
          class="inline-left"
          iconClass="text-red"
          v-else
          :title="item.total > 1 ? item.workspaceName + '(' + item.total + ')' : item.testScriptName"
          trailingTextClass="muted small"
        />
        <span v-if="item.checked == 0 || item.checked == undefined">{{momentTime(item.startTime, 'hh:mm')}}</span>
        <div v-else>
            <Icon
                icon="refresh"
                color="#007752"
                class="icon"
                @click="refreshExec($event, item)"
                />
            <Icon
                icon="file"
                color="#007752"
                class="icon"
                @click="showDetail($event, item)"
                />
        </div>
    </div>
    </List>
  </div>
</template>

<script setup lang="ts">
import List from './List.vue';
import ListItem from './ListItem.vue';
import {StateType} from "@/src/views/result/store";
import {PaginationConfig, QueryParams} from "@/types/data";
import Icon from './Icon.vue';
import {useI18n} from "vue-i18n";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import {computed, onMounted, watch, ref} from "vue";
import {momentUnixDefFormat} from "@/utils/datetime";

const { t } = useI18n();
const router = useRouter();

const momentTime = momentUnixDefFormat

const store = useStore<{ Result: StateType }>();
const models = computed<any[]>(() => store.state.Result.queryResult.result)

console.log("test result list:",models)
const pagination = computed<PaginationConfig>(() => store.state.Result.queryResult.pagination);
const queryParams = ref<QueryParams>({
    keywords: '', enabled: '1', page: pagination.value.page, pageSize: pagination.value.pageSize
});

const list = (page: number) => {
    store.dispatch('Result/list', {
    keywords: queryParams.value.keywords,
    enabled: queryParams.value.enabled,
    pageSize: pagination.value.pageSize,
    page: page});
}
list(1);

const refreshExec = (e, item) => {
    console.log(e, item)
    e.stopPropagation()
}

const showDetail = (e, item) => {
    console.log(item)
    e.stopPropagation()
}

const changeControlIcon = (e, index) => {
    for(let i=0; i < models.value.length; i++){
        if(i == index){
            models.value[index].checked = !models.value[index].checked;
        }else{
            models.value[i].checked = false;
        }
    }
}

onMounted(() => {
    console.log("onMounted")
})

</script>

<style scoped>
.list-item-container{
    display: flex;
    align-items: center;
    flex-direction: row;
    justify-content: space-between;
    width: 100%;
}

.icon{
    margin-left: 8px;
    cursor: pointer;
}

.checked{
    background-color: #E2E5E9;
}

.inline-left{
    min-width: 80%;
}
.result-list{
    padding-right: 20px;
}
</style>
